package types

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// Sudo creates a new Command with sudo privileges.
// This is a convenience function equivalent to Cmd(cmd, args...).Sudo().
//
// Example:
//
//	result := types.Sudo("systemctl", "restart", "nginx").Stdout()
func Sudo(cmd string, args ...string) *Command { return Cmd(cmd, args...).Sudo() }

// Command represents a system command that can be chained and piped together.
// Commands are executed lazily - they only run when an output method is called
// (Stdout, Stderr, Error, Run, etc.).
//
// Command supports:
//   - Command chaining via Pipe
//   - Function transformations via PipeFn
//   - Sudo execution
//   - Interactive mode for terminal input/output
//   - Input redirection
//
// Commands are idempotent - calling output methods multiple times executes the
// command only once and returns cached results.
type Command struct {
	// previous is the preceding command in a pipeline
	previous *Command
	// cmd is the command name to execute
	cmd string
	// cmdFn is an optional function to execute instead of a system command
	cmdFn func(stdin string) (stdout, stderr string, err error)
	// args are the command arguments
	args []string
	// interactive indicates if the command should connect to the terminal
	interactive bool
	// input is the stdin source
	input io.Reader
	// useSudo indicates if the command should run with sudo
	useSudo bool
	// executed tracks if the command has been run
	executed bool
	// stdout holds the captured stdout
	stdout string
	// stderr holds the captured stderr
	stderr string
	// err holds any execution error
	err error
	// ctx is the context for cancellation/timeout
	ctx context.Context
	// dir is the working directory for the command
	dir string
	// env holds environment variables to set
	env map[string]string
	// clearEnv indicates if the environment should be cleared
	clearEnv bool
	// exitCode holds the command's exit code
	exitCode int
	// retryCount is the number of retry attempts
	retryCount int
	// retryDelay is the delay between retries
	retryDelay time.Duration
}

// Cmd creates a new Command with the given command name and arguments.
// The command will not execute until an output method is called.
//
// Example:
//
//	cmd := types.Cmd("echo", "hello", "world")
//	output := cmd.Stdout() // "hello world\n"
func Cmd(cmd string, args ...string) *Command {
	return &Command{
		cmd:  cmd,
		args: args,
	}
}

// CmdFn creates a new Command from a function that transforms stdin to stdout/stderr.
// This allows inserting custom Go functions into command pipelines.
//
// The function receives stdin as a string and returns stdout, stderr, and an error.
//
// Example:
//
//	upperCase := types.CmdFn(func(stdin string) (string, string, error) {
//		return strings.ToUpper(stdin), "", nil
//	})
//	result := types.Cmd("echo", "hello").PipeFn(upperCase.cmdFn).Stdout() // "HELLO\n"
func CmdFn(fn func(stdin string) (stdout, stderr string, outErr error)) *Command {
	return &Command{
		cmdFn: fn,
	}
}

// Pipe chains another command to receive this command's stdout as stdin.
// This creates a pipeline similar to shell pipes (|).
//
// Errors from earlier commands in the pipeline prevent later commands from executing.
//
// Example:
//
//	result := types.Cmd("echo", "apple\nbanana\napricot").
//		Pipe("grep", "a").
//		Pipe("wc", "-l").
//		Stdout() // "3\n"
func (c *Command) Pipe(cmd string, args ...string) *Command {
	next := Cmd(cmd, args...)
	next.previous = c

	return next
}

// PipeFn chains a function to receive this command's stdout as stdin.
// This allows inserting custom transformations into command pipelines.
//
// Example:
//
//	result := types.Cmd("echo", "hello").
//		PipeFn(func(stdin string) (string, string, error) {
//			return strings.ToUpper(stdin), "", nil
//		}).
//		Stdout() // "HELLO\n"
func (c *Command) PipeFn(fn func(stdin string) (stdout, stderr string, outErr error)) *Command {
	next := CmdFn(fn)
	next.previous = c

	return next
}

// Interactive sets the command to run in interactive mode.
// In interactive mode, stdin/stdout/stderr are connected directly to the terminal
// instead of being captured. This is useful for commands that require user input
// or display progress.
//
// Example:
//
//	err := types.Cmd("vim", "file.txt").Interactive().Error()
func (c *Command) Interactive() *Command {
	c.interactive = true
	return c
}

// Input sets the stdin for the command from a string.
// This is useful for providing input to commands that read from stdin.
//
// Example:
//
//	result := types.Cmd("grep", "hello").
//		Input("hello world\ngoodbye\nhello again").
//		Stdout() // "hello world\nhello again\n"
func (c *Command) Input(input string) *Command {
	c.input = strings.NewReader(input)
	return c
}

// InputReader sets the stdin for the command from an io.Reader.
// This is useful for streaming input from files or other sources.
//
// Example:
//
//	file, _ := os.Open("input.txt")
//	result := types.Cmd("wc", "-l").InputReader(file).Stdout()
func (c *Command) InputReader(r io.Reader) *Command {
	c.input = r
	return c
}

// Sudo sets the command to run with sudo privileges.
// If sudo authentication is required, the user will be prompted interactively.
//
// Example:
//
//	err := types.Cmd("systemctl", "restart", "nginx").Sudo().Error()
func (c *Command) Sudo() *Command {
	c.useSudo = true
	return c
}

// Run executes the command and returns the Command for chaining.
// This is useful when you want to ensure execution but don't need the output.
//
// Example:
//
//	types.Cmd("mkdir", "-p", "/tmp/test").Run()
func (c *Command) Run() *Command { return c.execute() }

// Stdout executes the command and returns its stdout.
// Multiple calls return the cached result without re-executing.
//
// Example:
//
//	output := types.Cmd("echo", "hello").Stdout() // "hello\n"
func (c *Command) Stdout() string { return c.execute().stdout }

// StdoutErr executes the command and returns both stdout and any error.
// This is useful when you need both the output and error information.
//
// Example:
//
//	output, err := types.Cmd("ls", "/nonexistent").StdoutErr()
func (c *Command) StdoutErr() (string, error) { return c.Stdout(), c.Error() }

// Stderr executes the command and returns its stderr.
// Multiple calls return the cached result without re-executing.
//
// Example:
//
//	errMsg := types.Cmd("ls", "/nonexistent").Stderr()
func (c *Command) Stderr() string { return c.execute().stderr }

// StderrErr executes the command and returns both stderr and any error.
//
// Example:
//
//	errOutput, err := types.Cmd("ls", "/nonexistent").StderrErr()
func (c *Command) StderrErr() (string, error) { return c.Stderr(), c.Error() }

// Error executes the command and returns any error that occurred.
// Returns nil if the command executed successfully.
//
// Example:
//
//	if err := types.Cmd("false").Error(); err != nil {
//		// handle error
//	}
func (c *Command) Error() error { return c.execute().err }

// StdoutStderr executes the command and returns both stdout and stderr concatenated.
// This is useful when you need all output regardless of which stream it came from.
//
// Example:
//
//	allOutput := types.Cmd("ls", "/tmp", "/nonexistent").StdoutStderr()
func (c *Command) StdoutStderr() string { return c.execute().stdout + c.execute().stderr }

// WithContext sets the context for the command.
// The context can be used for cancellation or timeout.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	output := types.Cmd("sleep", "10").WithContext(ctx).Stdout() // cancelled after 5s
func (c *Command) WithContext(ctx context.Context) *Command {
	c.ctx = ctx
	return c
}

// WithTimeout sets a timeout for the command using a context with deadline.
// This is a convenience wrapper around WithContext.
//
// Example:
//
//	output := types.Cmd("sleep", "10").WithTimeout(5*time.Second).Stdout() // cancelled after 5s
func (c *Command) WithTimeout(duration time.Duration) *Command {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	// Store cancel function but don't call it - the command execution will handle it
	_ = cancel
	c.ctx = ctx
	return c
}

// WithDeadline sets an absolute deadline for the command using a context with deadline.
// This is a convenience wrapper around WithContext.
//
// Example:
//
//	deadline := time.Now().Add(5*time.Second)
//	output := types.Cmd("sleep", "10").WithDeadline(deadline).Stdout() // cancelled at deadline
func (c *Command) WithDeadline(t time.Time) *Command {
	ctx, cancel := context.WithDeadline(context.Background(), t)
	_ = cancel
	c.ctx = ctx
	return c
}

// Dir sets the working directory for the command.
// If not set, the command runs in the current working directory.
//
// Example:
//
//	output := types.Cmd("ls").Dir("/tmp").Stdout()
func (c *Command) Dir(path string) *Command {
	c.dir = path
	return c
}

// Env sets a single environment variable for the command.
// Can be called multiple times to set multiple variables.
//
// Example:
//
//	output := types.Cmd("printenv", "MY_VAR").Env("MY_VAR", "hello").Stdout()
func (c *Command) Env(key, value string) *Command {
	if c.env == nil {
		c.env = make(map[string]string)
	}
	c.env[key] = value
	return c
}

// EnvMap sets multiple environment variables from a map.
// Existing environment variables are preserved unless overridden.
//
// Example:
//
//	envVars := map[string]string{"PATH": "/usr/bin", "HOME": "/tmp"}
//	output := types.Cmd("env").EnvMap(envVars).Stdout()
func (c *Command) EnvMap(env map[string]string) *Command {
	if c.env == nil {
		c.env = make(map[string]string)
	}
	for k, v := range env {
		c.env[k] = v
	}
	return c
}

// ClearEnv clears all inherited environment variables.
// Only variables set via Env() or EnvMap() will be available.
//
// Example:
//
//	output := types.Cmd("env").ClearEnv().Env("ONLY_VAR", "value").Stdout()
func (c *Command) ClearEnv() *Command {
	c.clearEnv = true
	return c
}

// ExitCode returns the exit code of the command after execution.
// Returns 0 if the command hasn't been executed yet or succeeded.
// For non-zero exit codes, also check Error() for the error message.
//
// Example:
//
//	cmd := types.Cmd("false").Run()
//	if cmd.ExitCode() != 0 {
//		// handle non-zero exit
//	}
func (c *Command) ExitCode() int {
	c.execute()
	return c.exitCode
}

// Retry sets the number of retry attempts for the command.
// If the command fails, it will be retried up to the specified number of times.
// Use RetryWithBackoff for delays between retries.
//
// Example:
//
//	output := types.Cmd("curl", "http://example.com").Retry(3).Stdout()
func (c *Command) Retry(attempts int) *Command {
	c.retryCount = attempts
	return c
}

// RetryWithBackoff sets retry attempts with a delay between each retry.
// The delay is constant for all retries.
//
// Example:
//
//	output := types.Cmd("curl", "http://example.com").
//		RetryWithBackoff(3, 2*time.Second).
//		Stdout()
func (c *Command) RetryWithBackoff(attempts int, delay time.Duration) *Command {
	c.retryCount = attempts
	c.retryDelay = delay
	return c
}

// String implements fmt.Stringer and returns a string representation of the command.
// This shows the command that will be executed, including arguments.
//
// Example:
//
//	cmd := types.Cmd("echo", "hello", "world")
//	fmt.Println(cmd.String()) // "echo hello world"
func (c *Command) String() string {
	if c.cmd == "" {
		return "<function>"
	}

	parts := []string{c.cmd}
	parts = append(parts, c.args...)

	if c.useSudo {
		parts = append([]string{"sudo"}, parts...)
	}

	return strings.Join(parts, " ")
}

// StdoutTrimmed executes the command and returns stdout with leading/trailing whitespace removed.
// This is useful for commands that output single values with newlines.
//
// Example:
//
//	version := types.Cmd("git", "--version").StdoutTrimmed()
//	// "git version 2.39.0" (without trailing newline)
func (c *Command) StdoutTrimmed() string {
	return strings.TrimSpace(c.Stdout())
}

func (c *Command) execute() *Command {
	if c.executed {
		return c
	}

	c.executed = true

	// Retry logic wrapper
	maxAttempts := c.retryCount + 1
	if maxAttempts < 1 {
		maxAttempts = 1
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 && c.retryDelay > 0 {
			time.Sleep(c.retryDelay)
		}

		c.executeOnce()

		// If successful, break out of retry loop
		if c.err == nil {
			break
		}
	}

	return c
}

func (c *Command) executeOnce() {

	if c.cmdFn != nil {
		// Execute previous command first or read from input
		var stdin string
		if c.previous != nil {
			stdin, c.err = c.previous.StdoutErr()
			if c.err != nil {
				return
			}
		} else if c.input != nil {
			// Read from input reader
			buf := new(strings.Builder)
			_, c.err = io.Copy(buf, c.input)
			if c.err != nil {
				return
			}
			stdin = buf.String()
		}

		c.stdout, c.stderr, c.err = c.cmdFn(stdin)
		return
	}

	// Build command with sudo if needed
	var command *exec.Cmd

	// Use context if provided
	ctx := c.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	if c.useSudo {
		// Check if sudo is already authenticated (non-interactive)
		if err := Cmd("sudo", "-n", "true").Error(); err != nil {
			// Not authenticated, request authentication interactively
			if err := Cmd("sudo", "-v").Interactive().Error(); err != nil {
				c.err = err
				return
			}
		}

		command = exec.CommandContext(ctx, "sudo", append([]string{c.cmd}, c.args...)...)
	} else {
		command = exec.CommandContext(ctx, c.cmd, c.args...)
	}

	// Set working directory
	if c.dir != "" {
		command.Dir = c.dir
	}

	// Set environment variables
	if c.clearEnv {
		command.Env = []string{}
	}
	if c.env != nil {
		if !c.clearEnv {
			command.Env = os.Environ()
		}
		for k, v := range c.env {
			command.Env = append(command.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	// Set stdin
	if c.input != nil {
		command.Stdin = c.input
	} else if c.interactive {
		command.Stdin = os.Stdin
	}

	if c.previous != nil {
		prevOut, err := c.previous.StdoutErr()
		if err != nil {
			c.err = err
			return
		}

		// TODO stream the stdout instead of reading it all at once then making a reader.
		command.Stdin = strings.NewReader(prevOut)
	}

	// Set stdout/stderr based on mode
	if c.interactive {
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		c.err = command.Run()
	} else {
		// Capture stdout and stderr separately
		var stdoutBuf, stderrBuf strings.Builder
		command.Stdout = &stdoutBuf
		command.Stderr = &stderrBuf
		c.err = command.Run()
		c.stdout = stdoutBuf.String()
		c.stderr = stderrBuf.String()
	}

	// Extract exit code from error
	if c.err != nil {
		if exitErr, ok := c.err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				c.exitCode = status.ExitStatus()
			}
		}
	}
}
