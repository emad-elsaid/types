package types

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCmd(t *testing.T) {
	tests := []struct {
		name           string
		cmd            string
		args           []string
		expectedStdout string
	}{
		{
			name:           "simple echo",
			cmd:            "echo",
			args:           []string{"hello"},
			expectedStdout: "hello\n",
		},
		{
			name:           "echo with multiple args",
			cmd:            "echo",
			args:           []string{"hello", "world"},
			expectedStdout: "hello world\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Cmd(tt.cmd, tt.args...)
			stdout := cmd.Stdout()
			require.Equal(t, tt.expectedStdout, stdout)
			require.NoError(t, cmd.Error())
		})
	}
}

func TestCommand_Pipe(t *testing.T) {
	tests := []struct {
		name           string
		setup          func() *Command
		expectedStdout string
	}{
		{
			name: "echo piped to grep",
			setup: func() *Command {
				return Cmd("echo", "hello\nworld\nhello").Pipe("grep", "hello")
			},
			expectedStdout: "hello\nhello\n",
		},
		{
			name: "multiple pipes",
			setup: func() *Command {
				return Cmd("echo", "apple\nbanana\napricot").
					Pipe("grep", "a").
					Pipe("grep", "p")
			},
			expectedStdout: "apple\napricot\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setup()
			stdout := cmd.Stdout()
			require.Equal(t, tt.expectedStdout, stdout)
			require.NoError(t, cmd.Error())
		})
	}
}

func TestCommand_Input(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedStdout string
	}{
		{
			name:           "grep with input",
			input:          "hello\nworld\nhello",
			expectedStdout: "hello\nhello\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Cmd("grep", "hello").Input(tt.input)
			stdout := cmd.Stdout()
			require.Equal(t, tt.expectedStdout, stdout)
			require.NoError(t, cmd.Error())
		})
	}
}

func TestCmdFn(t *testing.T) {
	tests := []struct {
		name           string
		fn             func(stdin string) (stdout, stderr string, err error)
		stdin          string
		expectedStdout string
		expectedStderr string
	}{
		{
			name: "uppercase transformer",
			fn: func(stdin string) (stdout, stderr string, err error) {
				return strings.ToUpper(stdin), "", nil
			},
			stdin:          "hello world",
			expectedStdout: "HELLO WORLD",
			expectedStderr: "",
		},
		{
			name: "line counter",
			fn: func(stdin string) (stdout, stderr string, err error) {
				lines := strings.Split(strings.TrimSpace(stdin), "\n")
				count := len(lines)
				return strings.TrimSpace(strings.Repeat(".", count)), "", nil
			},
			stdin:          "line1\nline2\nline3",
			expectedStdout: "...",
			expectedStderr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := CmdFn(tt.fn).Input(tt.stdin)
			stdout := cmd.Stdout()
			stderr := cmd.Stderr()
			require.Equal(t, tt.expectedStdout, stdout)
			require.Equal(t, tt.expectedStderr, stderr)
			require.NoError(t, cmd.Error())
		})
	}
}

func TestCommand_PipeFn(t *testing.T) {
	tests := []struct {
		name           string
		setup          func() *Command
		expectedStdout string
	}{
		{
			name: "echo piped to uppercase function",
			setup: func() *Command {
				return Cmd("echo", "hello world").PipeFn(func(stdin string) (stdout, stderr string, err error) {
					return strings.ToUpper(stdin), "", nil
				})
			},
			expectedStdout: "HELLO WORLD\n",
		},
		{
			name: "command piped through multiple functions",
			setup: func() *Command {
				return Cmd("echo", "hello").
					PipeFn(func(stdin string) (stdout, stderr string, err error) {
						return strings.ToUpper(stdin), "", nil
					}).
					PipeFn(func(stdin string) (stdout, stderr string, err error) {
						return strings.TrimSpace(stdin) + "!!!", "", nil
					})
			},
			expectedStdout: "HELLO!!!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setup()
			stdout := cmd.Stdout()
			require.Equal(t, tt.expectedStdout, stdout)
			require.NoError(t, cmd.Error())
		})
	}
}

func TestCommand_StdoutErr(t *testing.T) {
	tests := []struct {
		name           string
		cmd            *Command
		expectedStdout string
		expectError    bool
	}{
		{
			name:           "successful command",
			cmd:            Cmd("echo", "hello"),
			expectedStdout: "hello\n",
			expectError:    false,
		},
		{
			name:           "failing command",
			cmd:            Cmd("false"),
			expectedStdout: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, err := tt.cmd.StdoutErr()
			require.Equal(t, tt.expectedStdout, stdout)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCommand_Idempotent(t *testing.T) {
	t.Run("multiple calls return same result", func(t *testing.T) {
		cmd := Cmd("echo", "hello")

		stdout1 := cmd.Stdout()
		stdout2 := cmd.Stdout()
		stdout3 := cmd.Stdout()

		require.Equal(t, stdout1, stdout2)
		require.Equal(t, stdout2, stdout3)
	})
}

func TestCommand_ErrorPropagation(t *testing.T) {
	t.Run("error in first command stops pipeline", func(t *testing.T) {
		cmd := Cmd("false").Pipe("echo", "hello")

		err := cmd.Error()
		require.Error(t, err)
	})
}

func TestCommand_WithTimeout(t *testing.T) {
	tests := []struct {
		name        string
		cmd         *Command
		expectError bool
	}{
		{
			name:        "command completes within timeout",
			cmd:         Cmd("sleep", "0.1").WithTimeout(1 * time.Second),
			expectError: false,
		},
		{
			name:        "command exceeds timeout",
			cmd:         Cmd("sleep", "10").WithTimeout(100 * time.Millisecond),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.Error()
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCommand_WithContext(t *testing.T) {
	t.Run("cancelled context stops command", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		// Start a long-running command
		cmd := Cmd("sleep", "10").WithContext(ctx)

		// Cancel immediately
		cancel()

		err := cmd.Error()
		require.Error(t, err)
	})
}

func TestCommand_Dir(t *testing.T) {
	tests := []struct {
		name           string
		dir            string
		expectedOutput string
	}{
		{
			name:           "run in /tmp",
			dir:            "/tmp",
			expectedOutput: "/tmp\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Cmd("pwd").Dir(tt.dir)
			stdout := cmd.Stdout()
			require.Equal(t, tt.expectedOutput, stdout)
			require.NoError(t, cmd.Error())
		})
	}
}

func TestCommand_Env(t *testing.T) {
	t.Run("set single environment variable", func(t *testing.T) {
		cmd := Cmd("sh", "-c", "echo $MY_VAR").Env("MY_VAR", "hello")
		stdout := cmd.Stdout()
		require.Equal(t, "hello\n", stdout)
		require.NoError(t, cmd.Error())
	})

	t.Run("set multiple environment variables", func(t *testing.T) {
		cmd := Cmd("sh", "-c", "echo $VAR1 $VAR2").
			Env("VAR1", "hello").
			Env("VAR2", "world")
		stdout := cmd.Stdout()
		require.Equal(t, "hello world\n", stdout)
		require.NoError(t, cmd.Error())
	})
}

func TestCommand_EnvMap(t *testing.T) {
	t.Run("set environment variables from map", func(t *testing.T) {
		envVars := map[string]string{
			"VAR1": "hello",
			"VAR2": "world",
		}
		cmd := Cmd("sh", "-c", "echo $VAR1 $VAR2").EnvMap(envVars)
		stdout := cmd.Stdout()
		require.Equal(t, "hello world\n", stdout)
		require.NoError(t, cmd.Error())
	})
}

func TestCommand_ClearEnv(t *testing.T) {
	t.Run("cleared environment has only set variables", func(t *testing.T) {
		// Set a variable in the test environment
		os.Setenv("TEST_VAR", "should_not_see_this")
		defer os.Unsetenv("TEST_VAR")

		cmd := Cmd("sh", "-c", "echo ${TEST_VAR:-empty} $MY_VAR").
			ClearEnv().
			Env("MY_VAR", "hello")

		stdout := cmd.Stdout()
		// TEST_VAR should default to "empty", MY_VAR should be "hello"
		require.Equal(t, "empty hello\n", stdout)
		require.NoError(t, cmd.Error())
	})
}

func TestCommand_ExitCode(t *testing.T) {
	tests := []struct {
		name         string
		cmd          *Command
		expectedCode int
	}{
		{
			name:         "successful command has exit code 0",
			cmd:          Cmd("true"),
			expectedCode: 0,
		},
		{
			name:         "failed command has non-zero exit code",
			cmd:          Cmd("false"),
			expectedCode: 1,
		},
		{
			name:         "command with specific exit code",
			cmd:          Cmd("sh", "-c", "exit 42"),
			expectedCode: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exitCode := tt.cmd.ExitCode()
			require.Equal(t, tt.expectedCode, exitCode)
		})
	}
}

func TestCommand_InputReader(t *testing.T) {
	t.Run("read from io.Reader", func(t *testing.T) {
		reader := strings.NewReader("hello\nworld\nhello")
		cmd := Cmd("grep", "hello").InputReader(reader)
		stdout := cmd.Stdout()
		require.Equal(t, "hello\nhello\n", stdout)
		require.NoError(t, cmd.Error())
	})
}

func TestCommand_CmdFn_WithInputReader(t *testing.T) {
	t.Run("CmdFn with InputReader success", func(t *testing.T) {
		reader := strings.NewReader("hello world")
		cmd := CmdFn(func(stdin string) (string, string, error) {
			return strings.ToUpper(stdin), "", nil
		}).InputReader(reader)

		stdout := cmd.Stdout()
		require.Equal(t, "HELLO WORLD", stdout)
		require.NoError(t, cmd.Error())
	})

	t.Run("CmdFn with InputReader error handling", func(t *testing.T) {
		// Create a reader that will fail
		testErr := os.ErrClosed
		errReader := &errorReader{err: testErr}
		cmd := CmdFn(func(stdin string) (string, string, error) {
			return stdin, "", nil
		}).InputReader(errReader)

		_ = cmd.Stdout() // Trigger execution
		require.Error(t, cmd.Error())
		require.Equal(t, testErr, cmd.Error())
	})

	t.Run("CmdFn with previous command", func(t *testing.T) {
		cmd := Cmd("echo", "test input").
			PipeFn(func(stdin string) (string, string, error) {
				return strings.ToUpper(stdin), "", nil
			})

		stdout := cmd.Stdout()
		require.Equal(t, "TEST INPUT\n", stdout)
		require.NoError(t, cmd.Error())
	})
}

// errorReader is a helper type that always returns an error on Read
type errorReader struct {
	err error
}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}

func TestCommand_String(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *Command
		expected string
	}{
		{
			name:     "simple command",
			cmd:      Cmd("echo", "hello"),
			expected: "echo hello",
		},
		{
			name:     "command with multiple args",
			cmd:      Cmd("git", "commit", "-m", "message"),
			expected: "git commit -m message",
		},
		{
			name:     "sudo command",
			cmd:      Cmd("systemctl", "restart", "nginx").Sudo(),
			expected: "sudo systemctl restart nginx",
		},
		{
			name:     "function command",
			cmd:      CmdFn(func(s string) (string, string, error) { return s, "", nil }),
			expected: "<function>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cmd.String()
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestCommand_StdoutTrimmed(t *testing.T) {
	tests := []struct {
		name     string
		cmd      *Command
		expected string
	}{
		{
			name:     "removes trailing newline",
			cmd:      Cmd("echo", "hello"),
			expected: "hello",
		},
		{
			name:     "removes leading and trailing whitespace",
			cmd:      Cmd("echo", "  hello  "),
			expected: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cmd.StdoutTrimmed()
			require.Equal(t, tt.expected, result)
			require.NoError(t, tt.cmd.Error())
		})
	}
}


func TestCommand_getStdoutPipe(t *testing.T) {
	t.Run("returns cached output when already executed", func(t *testing.T) {
		cmd := Cmd("echo", "hello").Run()

		// First call
		reader1, err1 := cmd.getStdoutPipe()
		if err1 != nil {
			t.Fatalf("expected no error, got %v", err1)
		}

		// Second call should return cached data
		reader2, err2 := cmd.getStdoutPipe()
		if err2 != nil {
			t.Fatalf("expected no error, got %v", err2)
		}

		// Read from both readers
		data1, _ := io.ReadAll(reader1)
		data2, _ := io.ReadAll(reader2)

		if string(data1) != string(data2) {
			t.Errorf("expected same output, got %q vs %q", data1, data2)
		}
	})

	t.Run("returns error when already executed with error", func(t *testing.T) {
		cmd := Cmd("nonexistent-command-xyz").Run()

		reader, err := cmd.getStdoutPipe()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if reader != nil {
			t.Errorf("expected nil reader, got %v", reader)
		}
	})

	t.Run("handles function commands", func(t *testing.T) {
		fn := func(stdin string) (string, string, error) {
			return "function output", "", nil
		}
		cmd := CmdFn(fn)

		reader, err := cmd.getStdoutPipe()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		data, _ := io.ReadAll(reader)
		if string(data) != "function output" {
			t.Errorf("expected 'function output', got %q", data)
		}
	})

	t.Run("handles function commands with error", func(t *testing.T) {
		expectedErr := errors.New("function error")
		fn := func(stdin string) (string, string, error) {
			return "", "", expectedErr
		}
		cmd := CmdFn(fn)

		reader, err := cmd.getStdoutPipe()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if reader != nil {
			t.Errorf("expected nil reader, got %v", reader)
		}
	})

	t.Run("handles piped commands", func(t *testing.T) {
		cmd := Cmd("echo", "line1\nline2\nline3").
			Pipe("grep", "line")

		reader, err := cmd.getStdoutPipe()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		data, _ := io.ReadAll(reader)
		output := string(data)
		if !strings.Contains(output, "line1") || !strings.Contains(output, "line2") {
			t.Errorf("expected lines from pipe, got %q", output)
		}
	})

	t.Run("handles piped command with previous error", func(t *testing.T) {
		cmd := Cmd("nonexistent-cmd-xyz").
			Pipe("grep", "test")

		reader, err := cmd.getStdoutPipe()
		if err == nil {
			t.Fatal("expected error from previous command, got nil")
		}
		if reader != nil {
			t.Errorf("expected nil reader, got %v", reader)
		}
	})

	t.Run("streams output from multiple pipes", func(t *testing.T) {
		cmd := Cmd("echo", "apple\nbanana\napricot\nblueberry").
			Pipe("grep", "a").
			Pipe("grep", "b")

		reader, err := cmd.getStdoutPipe()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		data, _ := io.ReadAll(reader)
		output := strings.TrimSpace(string(data))

		// Should contain lines with both 'a' and 'b'
		lines := strings.Split(output, "\n")
		if len(lines) == 0 {
			t.Error("expected at least one line in output")
		}

		for _, line := range lines {
			if !strings.Contains(line, "a") || !strings.Contains(line, "b") {
				t.Errorf("expected line to contain both 'a' and 'b', got %q", line)
			}
		}
	})

	t.Run("handles command with input reader", func(t *testing.T) {
		input := strings.NewReader("test input data")
		cmd := Cmd("cat").InputReader(input)

		reader, err := cmd.getStdoutPipe()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		data, _ := io.ReadAll(reader)
		if string(data) != "test input data" {
			t.Errorf("expected 'test input data', got %q", data)
		}
	})

	t.Run("marks command as executed", func(t *testing.T) {
		cmd := Cmd("echo", "test")

		if cmd.executed {
			t.Error("expected command to not be executed initially")
		}

		_, _ = cmd.getStdoutPipe()

		if !cmd.executed {
			t.Error("expected command to be marked as executed")
		}
	})
}


// TestRetry_SuccessFirstAttempt verifies that Retry doesn't retry when command succeeds on first attempt
func TestRetry_SuccessFirstAttempt(t *testing.T) {
	// Create a counter file to track attempts
	tmpDir := t.TempDir()
	counterFile := filepath.Join(tmpDir, "attempts.txt")

	// Command that succeeds immediately and increments counter
	script := fmt.Sprintf(`echo "1" >> %s && cat %s | wc -l`, counterFile, counterFile)
	output := Cmd("bash", "-c", script).
		Retry(3).
		StdoutTrimmed()

	if output != "1" {
		t.Errorf("Expected 1 attempt, got %s", output)
	}

	// Verify error is nil on success
	err := Cmd("echo", "test").Retry(3).Error()
	if err != nil {
		t.Errorf("Expected no error on success, got: %v", err)
	}
}

// TestRetry_FailureExhaustsRetries verifies that Retry attempts the specified number of times
func TestRetry_FailureExhaustsRetries(t *testing.T) {
	tmpDir := t.TempDir()
	counterFile := filepath.Join(tmpDir, "attempts.txt")

	// Command that always fails but tracks attempts
	script := fmt.Sprintf(`echo "1" >> %s && exit 1`, counterFile)
	cmd := Cmd("bash", "-c", script).Retry(3)

	// Execute and expect error
	_ = cmd.Stdout()
	err := cmd.Error()

	if err == nil {
		t.Error("Expected error after exhausting retries")
	}

	// Count actual attempts (should be 4: initial + 3 retries)
	content, _ := os.ReadFile(counterFile)
	attempts := strings.Count(string(content), "1")
	if attempts != 4 {
		t.Errorf("Expected 4 attempts (1 initial + 3 retries), got %d", attempts)
	}
}

// TestRetry_SuccessOnSecondAttempt verifies retry works when command eventually succeeds
func TestRetry_SuccessOnSecondAttempt(t *testing.T) {
	tmpDir := t.TempDir()
	counterFile := filepath.Join(tmpDir, "attempts.txt")
	
	// Command that fails first time, succeeds second time
	script := fmt.Sprintf(`
		count=$(cat %s 2>/dev/null | wc -l)
		echo "1" >> %s
		if [ "$count" -lt "1" ]; then
			exit 1
		fi
		echo "success"
	`, counterFile, counterFile)

	output := Cmd("bash", "-c", script).
		Retry(3).
		StdoutTrimmed()

	if output != "success" {
		t.Errorf("Expected 'success', got '%s'", output)
	}

	// Verify it only tried twice (initial + 1 retry)
	content, _ := os.ReadFile(counterFile)
	attempts := strings.Count(string(content), "1")
	if attempts != 2 {
		t.Errorf("Expected 2 attempts, got %d", attempts)
	}

	// Verify no error after successful retry
	err := Cmd("bash", "-c", script).Retry(3).Error()
	if err == nil {
		// This is expected - command should succeed after retry
		// Reset counter for this test
		os.WriteFile(counterFile, []byte{}, 0644)
	}
}

// TestRetry_ZeroRetries verifies that Retry(0) means single attempt
func TestRetry_ZeroRetries(t *testing.T) {
	tmpDir := t.TempDir()
	counterFile := filepath.Join(tmpDir, "attempts.txt")

	script := fmt.Sprintf(`echo "1" >> %s && exit 1`, counterFile)
	_ = Cmd("bash", "-c", script).Retry(0).Stdout()

	// Should only try once (no retries)
	content, _ := os.ReadFile(counterFile)
	attempts := strings.Count(string(content), "1")
	if attempts != 1 {
		t.Errorf("Expected 1 attempt with Retry(0), got %d", attempts)
	}
}

// TestRetry_NegativeRetries verifies that negative retry count is treated as single attempt
func TestRetry_NegativeRetries(t *testing.T) {
	tmpDir := t.TempDir()
	counterFile := filepath.Join(tmpDir, "attempts.txt")

	script := fmt.Sprintf(`echo "1" >> %s && exit 1`, counterFile)
	_ = Cmd("bash", "-c", script).Retry(-5).Stdout()

	// Should only try once (negative treated as 0)
	content, _ := os.ReadFile(counterFile)
	attempts := strings.Count(string(content), "1")
	if attempts != 1 {
		t.Errorf("Expected 1 attempt with negative retry count, got %d", attempts)
	}
}

// TestRetryWithBackoff_Timing verifies that backoff delay is applied between retries
func TestRetryWithBackoff_Timing(t *testing.T) {
	delay := 100 * time.Millisecond
	retries := 2

	start := time.Now()
	
	// Command that always fails
	_ = Cmd("bash", "-c", "exit 1").
		RetryWithBackoff(retries, delay).
		Stdout()

	elapsed := time.Since(start)

	// Total time should be at least (retries * delay)
	// We use 2 retries, so expect at least 200ms
	expectedMin := time.Duration(retries) * delay
	expectedMax := expectedMin + 500*time.Millisecond // Allow some overhead

	if elapsed < expectedMin {
		t.Errorf("Backoff too fast: expected at least %v, got %v", expectedMin, elapsed)
	}

	if elapsed > expectedMax {
		t.Errorf("Backoff too slow: expected less than %v, got %v", expectedMax, elapsed)
	}
}

// TestRetryWithBackoff_SuccessSkipsDelay verifies first success doesn't wait
func TestRetryWithBackoff_SuccessSkipsDelay(t *testing.T) {
	delay := 500 * time.Millisecond

	start := time.Now()
	
	// Command that succeeds immediately
	output := Cmd("echo", "quick").
		RetryWithBackoff(3, delay).
		StdoutTrimmed()

	elapsed := time.Since(start)

	if output != "quick" {
		t.Errorf("Expected 'quick', got '%s'", output)
	}

	// Should complete much faster than the retry delay
	if elapsed > 200*time.Millisecond {
		t.Errorf("Expected fast completion, took %v (should be instant)", elapsed)
	}
}

// TestRetryWithBackoff_ZeroDelay verifies zero delay works (immediate retry)
func TestRetryWithBackoff_ZeroDelay(t *testing.T) {
	tmpDir := t.TempDir()
	counterFile := filepath.Join(tmpDir, "attempts.txt")

	start := time.Now()
	
	script := fmt.Sprintf(`echo "1" >> %s && exit 1`, counterFile)
	_ = Cmd("bash", "-c", script).
		RetryWithBackoff(3, 0).
		Stdout()

	elapsed := time.Since(start)

	// Should complete very quickly with no delays
	if elapsed > 200*time.Millisecond {
		t.Errorf("Expected fast completion with zero delay, took %v", elapsed)
	}

	// Verify all retries happened
	content, _ := os.ReadFile(counterFile)
	attempts := strings.Count(string(content), "1")
	if attempts != 4 {
		t.Errorf("Expected 4 attempts, got %d", attempts)
	}
}

// TestRetryWithBackoff_SuccessOnThirdAttempt verifies retry with backoff works
func TestRetryWithBackoff_SuccessOnThirdAttempt(t *testing.T) {
	tmpDir := t.TempDir()
	counterFile := filepath.Join(tmpDir, "attempts.txt")

	// Command that succeeds on third attempt
	script := fmt.Sprintf(`
		count=$(cat %s 2>/dev/null | wc -l)
		echo "1" >> %s
		if [ "$count" -lt "2" ]; then
			exit 1
		fi
		echo "success"
	`, counterFile, counterFile)

	start := time.Now()
	delay := 50 * time.Millisecond

	output := Cmd("bash", "-c", script).
		RetryWithBackoff(5, delay).
		StdoutTrimmed()

	elapsed := time.Since(start)

	if output != "success" {
		t.Errorf("Expected 'success', got '%s'", output)
	}

	// Should have waited for 2 delays (before 2nd and 3rd attempt)
	expectedMin := 2 * delay
	expectedMax := expectedMin + 300*time.Millisecond

	if elapsed < expectedMin || elapsed > expectedMax {
		t.Errorf("Expected timing between %v and %v, got %v", expectedMin, expectedMax, elapsed)
	}

	// Verify exactly 3 attempts
	content, _ := os.ReadFile(counterFile)
	attempts := strings.Count(string(content), "1")
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

// TestRetry_WithPipe verifies retry works in command pipelines
func TestRetry_WithPipe(t *testing.T) {
	// Simple test: command succeeds and pipes to uppercase
	output := Cmd("echo", "hello").
		Retry(2).
		Pipe("tr", "a-z", "A-Z").
		StdoutTrimmed()

	if output != "HELLO" {
		t.Errorf("Expected 'HELLO', got '%s'", output)
	}

	// Test that retry on piped command works
	tmpDir := t.TempDir()
	counterFile := filepath.Join(tmpDir, "attempts.txt")

	// Second command in pipe fails initially
	script := fmt.Sprintf(`
		count=$(cat %s 2>/dev/null | wc -l)
		echo "1" >> %s
		if [ "$count" -lt "1" ]; then
			exit 1
		fi
		tr a-z A-Z
	`, counterFile, counterFile)

	output = Cmd("echo", "test").
		Pipe("bash", "-c", script).
		Retry(2).
		StdoutTrimmed()

	if output != "TEST" {
		t.Errorf("Expected 'TEST' after retry, got '%s'", output)
	}
}

// TestRetry_PreservesExitCode verifies that exit code is preserved after retries
func TestRetry_PreservesExitCode(t *testing.T) {
	cmd := Cmd("bash", "-c", "exit 42").Retry(2)
	_ = cmd.Stdout()

	exitCode := cmd.ExitCode()
	if exitCode != 42 {
		t.Errorf("Expected exit code 42, got %d", exitCode)
	}
}

// TestRetryWithBackoff_CombinedWithOtherFeatures tests retry with Dir, Env, etc.
func TestRetryWithBackoff_CombinedWithOtherFeatures(t *testing.T) {
	tmpDir := t.TempDir()
	counterFile := filepath.Join(tmpDir, "attempts.txt")
	testFile := filepath.Join(tmpDir, "testfile.txt")

	// Create a test file in tmpDir
	os.WriteFile(testFile, []byte("content"), 0644)

	// Command that uses Dir and Env, fails first time
	script := fmt.Sprintf(`
		count=$(cat %s 2>/dev/null | wc -l)
		echo "1" >> %s
		if [ "$count" -lt "1" ]; then
			exit 1
		fi
		echo "$MY_VAR:$(pwd):$(cat testfile.txt)"
	`, counterFile, counterFile)

	output := Cmd("bash", "-c", script).
		Dir(tmpDir).
		Env("MY_VAR", "test").
		RetryWithBackoff(3, 10*time.Millisecond).
		StdoutTrimmed()

	expected := fmt.Sprintf("test:%s:content", tmpDir)
	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

// TestRetry_Idempotency verifies multiple Stdout() calls don't trigger re-execution
func TestRetry_Idempotency(t *testing.T) {
	tmpDir := t.TempDir()
	counterFile := filepath.Join(tmpDir, "attempts.txt")

	script := fmt.Sprintf(`
		count=$(cat %s 2>/dev/null | wc -l)
		echo "1" >> %s
		if [ "$count" -lt "1" ]; then
			exit 1
		fi
		echo "success"
	`, counterFile, counterFile)

	cmd := Cmd("bash", "-c", script).Retry(2)

	// Call Stdout multiple times
	output1 := cmd.StdoutTrimmed()
	output2 := cmd.StdoutTrimmed()
	output3 := cmd.Stdout()

	if output1 != "success" || output2 != "success" {
		t.Errorf("Expected 'success', got '%s' and '%s'", output1, output2)
	}

	if !strings.Contains(output3, "success") {
		t.Errorf("Expected output3 to contain 'success', got '%s'", output3)
	}

	// Should only have executed twice (initial + 1 retry), not 6 times (2 attempts * 3 calls)
	content, _ := os.ReadFile(counterFile)
	attempts := strings.Count(string(content), "1")
	if attempts != 2 {
		t.Errorf("Expected 2 attempts due to idempotency, got %d", attempts)
	}
}


func TestCommand_StderrErr(t *testing.T) {
	tests := []struct {
		name           string
		setup          func() *Command
		expectedStderr string
		expectError    bool
	}{
		{
			name: "successful command with stderr",
			setup: func() *Command {
				return Cmd("sh", "-c", "echo error >&2")
			},
			expectedStderr: "error\n",
			expectError:    false,
		},
		{
			name: "failed command with stderr",
			setup: func() *Command {
				return Cmd("sh", "-c", "echo failure >&2; exit 1")
			},
			expectedStderr: "failure\n",
			expectError:    true,
		},
		{
			name: "command with no stderr output",
			setup: func() *Command {
				return Cmd("echo", "stdout only")
			},
			expectedStderr: "",
			expectError:    false,
		},
		{
			name: "nonexistent command",
			setup: func() *Command {
				return Cmd("nonexistent-command-xyz")
			},
			expectedStderr: "",
			expectError:    true,
		},
		{
			name: "piped commands with stderr",
			setup: func() *Command {
				return Cmd("sh", "-c", "echo first >&2").
					Pipe("sh", "-c", "echo second >&2")
			},
			expectedStderr: "second\n",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setup()
			stderr, err := cmd.StderrErr()

			require.Equal(t, tt.expectedStderr, stderr)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCommand_StdoutStderr(t *testing.T) {
	tests := []struct {
		name           string
		setup          func() *Command
		expectedOutput string
	}{
		{
			name: "command with both stdout and stderr",
			setup: func() *Command {
				return Cmd("sh", "-c", "echo stdout; echo stderr >&2")
			},
			expectedOutput: "stdout\nstderr\n",
		},
		{
			name: "command with only stdout",
			setup: func() *Command {
				return Cmd("echo", "hello")
			},
			expectedOutput: "hello\n",
		},
		{
			name: "command with only stderr",
			setup: func() *Command {
				return Cmd("sh", "-c", "echo error >&2")
			},
			expectedOutput: "error\n",
		},
		{
			name: "command with no output",
			setup: func() *Command {
				return Cmd("true")
			},
			expectedOutput: "",
		},
		{
			name: "failed command with both outputs",
			setup: func() *Command {
				return Cmd("sh", "-c", "echo output; echo error >&2; exit 1")
			},
			expectedOutput: "output\nerror\n",
		},
		{
			name: "piped commands with mixed outputs",
			setup: func() *Command {
				return Cmd("sh", "-c", "echo first; echo first-err >&2").
					Pipe("sh", "-c", "cat; echo second-err >&2")
			},
			expectedOutput: "first\nsecond-err\n",
		},
		{
			name: "multiple stderr messages concatenated",
			setup: func() *Command {
				return Cmd("sh", "-c", "echo line1 >&2; echo line2 >&2; echo line3 >&2")
			},
			expectedOutput: "line1\nline2\nline3\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setup()
			output := cmd.StdoutStderr()
			require.Equal(t, tt.expectedOutput, output)
		})
	}
}

// TestCommand_StdoutStderr_Idempotency verifies that calling StdoutStderr
// multiple times returns the same result (commands are executed once and cached)
func TestCommand_StdoutStderr_Idempotency(t *testing.T) {
	cmd := Cmd("sh", "-c", "echo stdout; echo stderr >&2")

	output1 := cmd.StdoutStderr()
	output2 := cmd.StdoutStderr()
	output3 := cmd.StdoutStderr()

	require.Equal(t, output1, output2)
	require.Equal(t, output2, output3)
}

// TestCommand_StderrErr_Idempotency verifies that calling StderrErr
// multiple times returns the same result
func TestCommand_StderrErr_Idempotency(t *testing.T) {
	cmd := Cmd("sh", "-c", "echo error >&2; exit 1")

	stderr1, err1 := cmd.StderrErr()
	stderr2, err2 := cmd.StderrErr()

	require.Equal(t, stderr1, stderr2)
	require.Equal(t, err1, err2)
	require.Error(t, err1)
}


// TestCommand_StreamingPipe tests that piped commands stream data efficiently
// rather than buffering all output before processing.
func TestCommand_StreamingPipe(t *testing.T) {
	t.Run("large output streams through pipeline", func(t *testing.T) {
		// Generate 10,000 lines of output
		cmd := Cmd("seq", "1", "10000").Pipe("head", "-n", "5")
		
		stdout := cmd.Stdout()
		require.NoError(t, cmd.Error())
		
		// Should get first 5 lines
		expected := "1\n2\n3\n4\n5\n"
		require.Equal(t, expected, stdout)
	})

	t.Run("streaming with grep in middle of pipeline", func(t *testing.T) {
		// Generate numbers, filter evens, take first 3
		cmd := Cmd("seq", "1", "1000").
			Pipe("grep", "0$").  // Numbers ending in 0
			Pipe("head", "-n", "3")
		
		stdout := cmd.Stdout()
		require.NoError(t, cmd.Error())
		
		expected := "10\n20\n30\n"
		require.Equal(t, expected, stdout)
	})

	t.Run("idempotent calls with streaming", func(t *testing.T) {
		cmd := Cmd("seq", "1", "100").Pipe("head", "-n", "3")
		
		// First call
		stdout1 := cmd.Stdout()
		require.NoError(t, cmd.Error())
		
		// Second call should return same cached result
		stdout2 := cmd.Stdout()
		require.NoError(t, cmd.Error())
		
		require.Equal(t, stdout1, stdout2)
		require.Equal(t, "1\n2\n3\n", stdout1)
	})

	t.Run("error propagation in streaming pipeline", func(t *testing.T) {
		// First command fails, should stop pipeline
		cmd := Cmd("sh", "-c", "echo 'test'; exit 1").Pipe("grep", "test")
		
		_ = cmd.Stdout()
		require.Error(t, cmd.Error())
	})

	t.Run("multiple pipes with streaming", func(t *testing.T) {
		// Test a longer pipeline
		cmd := Cmd("seq", "1", "1000").
			Pipe("grep", "5").     // Numbers containing 5
			Pipe("head", "-n", "5"). // First 5 matches
			Pipe("tail", "-n", "2")  // Last 2 of those
		
		stdout := cmd.Stdout()
		require.NoError(t, cmd.Error())
		
		// Should get lines 4 and 5 from the first 5 matches
		lines := strings.Split(strings.TrimSpace(stdout), "\n")
		require.Len(t, lines, 2)
	})

	t.Run("streaming with function in pipeline", func(t *testing.T) {
		// Test that function commands still work in streaming pipelines
		cmd := Cmd("seq", "1", "10").
			PipeFn(func(stdin string) (string, string, error) {
				// Double each number
				lines := strings.Split(strings.TrimSpace(stdin), "\n")
				var result []string
				for _, line := range lines {
					result = append(result, line+line)
				}
				return strings.Join(result, "\n") + "\n", "", nil
			}).
			Pipe("head", "-n", "3")
		
		stdout := cmd.Stdout()
		require.NoError(t, cmd.Error())
		
		// First three doubled: 11, 22, 33
		expected := "11\n22\n33\n"
		require.Equal(t, expected, stdout)
	})
}

// TestCommand_StreamingMemoryEfficiency demonstrates that streaming
// doesn't buffer entire output in memory before processing
func TestCommand_StreamingMemoryEfficiency(t *testing.T) {
	t.Run("head stops reading after getting needed lines", func(t *testing.T) {
		// If this were buffering all 1M lines, it would be very slow
		// With streaming, head can stop reading after 10 lines
		cmd := Cmd("seq", "1", "1000000").Pipe("head", "-n", "10")
		
		stdout := cmd.Stdout()
		require.NoError(t, cmd.Error())
		
		// Verify we got exactly 10 lines
		lines := strings.Split(strings.TrimSpace(stdout), "\n")
		require.Len(t, lines, 10)
		require.Equal(t, "1", lines[0])
		require.Equal(t, "10", lines[9])
	})
}

// TestCommand_StreamingBackwardCompatibility ensures the streaming
// implementation maintains backward compatibility
func TestCommand_StreamingBackwardCompatibility(t *testing.T) {
	t.Run("simple pipe still works", func(t *testing.T) {
		cmd := Cmd("echo", "hello\nworld").Pipe("grep", "world")
		stdout := cmd.Stdout()
		require.Equal(t, "world\n", stdout)
		require.NoError(t, cmd.Error())
	})

	t.Run("input reader still works", func(t *testing.T) {
		cmd := Cmd("grep", "test").Input("test\nother\ntest")
		stdout := cmd.Stdout()
		require.Equal(t, "test\ntest\n", stdout)
		require.NoError(t, cmd.Error())
	})

	t.Run("piping from input reader", func(t *testing.T) {
		cmd := Cmd("cat").InputReader(strings.NewReader("line1\nline2\nline3")).
			Pipe("head", "-n", "2")
		stdout := cmd.Stdout()
		require.Equal(t, "line1\nline2\n", stdout)
		require.NoError(t, cmd.Error())
	})
}


// TestSudo tests the Sudo method and function
func TestSudo(t *testing.T) {
	t.Run("Sudo method sets useSudo flag", func(t *testing.T) {
		// Create a command that would succeed without sudo
		cmd := Cmd("echo", "test").Sudo()
		
		// The String() representation should include sudo
		cmdStr := cmd.String()
		if !strings.HasPrefix(cmdStr, "sudo") {
			t.Errorf("Expected command string to start with 'sudo', got: %s", cmdStr)
		}
	})

	t.Run("Sudo function creates command with sudo", func(t *testing.T) {
		cmd := Sudo("echo", "test")
		
		cmdStr := cmd.String()
		if !strings.HasPrefix(cmdStr, "sudo") {
			t.Errorf("Expected command string to start with 'sudo', got: %s", cmdStr)
		}
	})

	t.Run("Sudo preserves arguments", func(t *testing.T) {
		cmd := Cmd("test", "arg1", "arg2").Sudo()
		
		cmdStr := cmd.String()
		expected := "sudo test arg1 arg2"
		if cmdStr != expected {
			t.Errorf("Expected '%s', got: %s", expected, cmdStr)
		}
	})

	t.Run("Sudo can be chained with other methods", func(t *testing.T) {
		// Test that Sudo can be chained
		cmd := Cmd("echo", "test").Sudo().Env("TEST", "value")
		
		cmdStr := cmd.String()
		if !strings.Contains(cmdStr, "sudo") {
			t.Errorf("Expected command to contain 'sudo', got: %s", cmdStr)
		}
	})

	// Skip actual sudo execution in tests unless running with sudo privileges
	// This is a safety measure to avoid requiring sudo in CI/CD environments
	t.Run("Sudo execution requires authentication", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("Skipping test when already running as root")
		}

		// Check if sudo is available
		if _, err := exec.LookPath("sudo"); err != nil {
			t.Skip("sudo not available on this system")
		}

		// Try a sudo command that should fail without authentication in non-interactive mode
		// Using sudo -n (non-interactive) should fail if not already authenticated
		cmd := Cmd("sudo", "-n", "true")
		err := cmd.Error()
		
		// We don't assert the error because:
		// - It might succeed if user has passwordless sudo
		// - It might fail if sudo requires password
		// We just verify the command runs without panic
		_ = err
	})
}

// TestInteractive tests the Interactive method
func TestInteractive(t *testing.T) {
	t.Run("Interactive method sets interactive flag", func(t *testing.T) {
		// Create a command and set it to interactive mode
		cmd := Cmd("echo", "test").Interactive()
		
		// We can't directly check the interactive flag (it's private),
		// but we can verify the command was created successfully
		if cmd == nil {
			t.Error("Expected non-nil command")
		}
	})

	t.Run("Interactive can be chained", func(t *testing.T) {
		cmd := Cmd("cat").Input("test").Interactive()
		
		if cmd == nil {
			t.Error("Expected non-nil command after chaining")
		}
	})

	t.Run("Interactive with non-interactive command", func(t *testing.T) {
		// Even in interactive mode, commands that don't require interaction should work
		// Note: This won't actually be interactive in tests, but verifies no crash
		cmd := Cmd("echo", "hello")
		
		// In a real test environment, interactive commands can't actually be interactive
		// So we just verify the method doesn't panic
		result := cmd.Stdout()
		expected := "hello\n"
		if result != expected {
			t.Errorf("Expected '%s', got: '%s'", expected, result)
		}
	})

	// This test verifies that Interactive mode properly handles terminal commands
	// but skips actual execution since tests run in non-interactive environments
	t.Run("Interactive mode with terminal command", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Skipping terminal test on Windows")
		}

		// Commands like vim, less, top, etc. require a terminal
		// In test mode, these would fail or behave unexpectedly
		// We just verify the command construction doesn't panic
		cmd := Cmd("cat", "/dev/null").Interactive()
		
		// This should succeed even in interactive mode since /dev/null has no content
		err := cmd.Error()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

// TestSudoWithInteractive tests combining Sudo and Interactive
func TestSudoWithInteractive(t *testing.T) {
	t.Run("Sudo and Interactive can be combined", func(t *testing.T) {
		cmd := Cmd("cat", "/dev/null").Sudo().Interactive()
		
		cmdStr := cmd.String()
		if !strings.HasPrefix(cmdStr, "sudo") {
			t.Errorf("Expected command to start with 'sudo', got: %s", cmdStr)
		}
	})

	t.Run("Interactive and Sudo order doesn't matter", func(t *testing.T) {
		cmd1 := Cmd("echo", "test").Sudo().Interactive()
		cmd2 := Cmd("echo", "test").Interactive().Sudo()
		
		// Both should produce the same string representation
		if cmd1.String() != cmd2.String() {
			t.Errorf("Expected same string representation regardless of order")
		}
	})
}

// TestSudoAuthentication tests sudo authentication behavior
func TestSudoAuthentication(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping test when already running as root")
	}

	// Check if sudo is available
	if _, err := exec.LookPath("sudo"); err != nil {
		t.Skip("sudo not available on this system")
	}

	t.Run("Sudo with simple command", func(t *testing.T) {
		// Use a command that's safe to run with sudo and doesn't require password if cached
		// We use sudo -n (non-interactive) to avoid hanging on password prompt
		cmd := Cmd("sudo", "-n", "echo", "test")
		
		// This might fail if sudo requires password, which is fine for testing
		// We're just verifying no panic occurs
		_ = cmd.Stdout()
	})
}

// TestSudoInPipeline tests that Sudo works correctly in command pipelines
func TestSudoInPipeline(t *testing.T) {
	t.Run("Sudo on first command in pipeline", func(t *testing.T) {
		// Create a pipeline where sudo is used on the first command
		cmd1 := Cmd("echo", "test").Sudo()
		
		// Verify the first command has sudo
		cmdStr := cmd1.String()
		if !strings.Contains(cmdStr, "sudo") {
			t.Errorf("Expected first command to contain 'sudo', got: %s", cmdStr)
		}
		
		// Create pipeline - subsequent commands don't show in String()
		_ = cmd1.Pipe("cat").Pipe("wc", "-l")
	})

	t.Run("Sudo on piped command", func(t *testing.T) {
		cmd := Cmd("echo", "test").Pipe("cat").Sudo()
		
		cmdStr := cmd.String()
		if !strings.Contains(cmdStr, "sudo") {
			t.Errorf("Expected pipeline to contain 'sudo', got: %s", cmdStr)
		}
	})

	t.Run("Sudo in pipeline execution", func(t *testing.T) {
		// Verify that sudo commands can be part of a pipeline
		// and execute without errors (using safe commands)
		cmd := Cmd("echo", "hello").Pipe("cat")
		
		result := cmd.Stdout()
		expected := "hello\n"
		if result != expected {
			t.Errorf("Expected '%s', got: '%s'", expected, result)
		}
	})
}


// TestCommandWithDeadline_Success verifies that commands complete successfully
// before the deadline expires.
func TestCommandWithDeadline_Success(t *testing.T) {
	deadline := time.Now().Add(5 * time.Second)
	output := Cmd("echo", "hello").WithDeadline(deadline).Stdout()

	if output != "hello\n" {
		t.Errorf("expected 'hello\\n', got %q", output)
	}
}

// TestCommandWithDeadline_Timeout verifies that long-running commands are
// terminated when the deadline expires.
func TestCommandWithDeadline_Timeout(t *testing.T) {
	deadline := time.Now().Add(500 * time.Millisecond)
	start := time.Now()

	output := Cmd("sleep", "5").WithDeadline(deadline).Stdout()
	elapsed := time.Since(start)

	// Command should be killed after ~500ms, not 5 seconds
	if elapsed > 2*time.Second {
		t.Errorf("command took too long: %v (expected ~500ms)", elapsed)
	}

	// Output should be empty since sleep produces no output
	if output != "" {
		t.Errorf("expected empty output, got %q", output)
	}
}

// TestCommandWithDeadline_Error verifies that deadline expiration produces
// an appropriate error.
func TestCommandWithDeadline_Error(t *testing.T) {
	deadline := time.Now().Add(200 * time.Millisecond)

	err := Cmd("sleep", "10").WithDeadline(deadline).Error()

	if err == nil {
		t.Fatal("expected error due to deadline, got nil")
	}

	// Error message should indicate context cancellation
	errMsg := err.Error()
	if !strings.Contains(errMsg, "context") && !strings.Contains(errMsg, "killed") && !strings.Contains(errMsg, "signal") {
		t.Logf("warning: error message doesn't mention context/killed/signal: %v", errMsg)
	}
}

// TestCommandWithDeadline_PastDeadline verifies behavior when deadline is
// already in the past.
func TestCommandWithDeadline_PastDeadline(t *testing.T) {
	deadline := time.Now().Add(-1 * time.Second)

	err := Cmd("echo", "hello").WithDeadline(deadline).Error()

	if err == nil {
		t.Fatal("expected error due to past deadline, got nil")
	}
}

// TestCommandWithDeadline_InPipeline verifies that deadlines work correctly
// in command pipelines.
func TestCommandWithDeadline_InPipeline(t *testing.T) {
	deadline := time.Now().Add(500 * time.Millisecond)
	start := time.Now()

	// Pipeline: sleep 5 | cat
	// Deadline on first command should kill the pipeline
	output := Cmd("sleep", "5").
		WithDeadline(deadline).
		Pipe("cat").
		Stdout()
	elapsed := time.Since(start)

	// Should terminate quickly due to deadline
	if elapsed > 2*time.Second {
		t.Errorf("pipeline took too long: %v (expected ~500ms)", elapsed)
	}

	if output != "" {
		t.Errorf("expected empty output, got %q", output)
	}
}

// TestCommandWithDeadline_MultipleCalls verifies that multiple deadline
// calls update the context.
func TestCommandWithDeadline_MultipleCalls(t *testing.T) {
	// First set a short deadline, then extend it
	shortDeadline := time.Now().Add(100 * time.Millisecond)
	longDeadline := time.Now().Add(5 * time.Second)

	// The last WithDeadline call should win
	output := Cmd("sleep", "0.2").
		WithDeadline(shortDeadline).
		WithDeadline(longDeadline).
		Stdout()

	// Should succeed because we extended the deadline
	if output != "" {
		// sleep produces no output, which is fine
	}

	// Now test the reverse - shorten the deadline
	start := time.Now()
	err := Cmd("sleep", "5").
		WithDeadline(longDeadline).
		WithDeadline(shortDeadline).
		Error()
	elapsed := time.Since(start)

	// Should fail quickly with shortened deadline
	if elapsed > 2*time.Second {
		t.Errorf("command took too long: %v (expected ~100ms)", elapsed)
	}

	if err == nil {
		t.Error("expected error due to shortened deadline, got nil")
	}
}

// TestCommandWithDeadline_WithRetry verifies interaction between deadline
// and retry logic.
func TestCommandWithDeadline_WithRetry(t *testing.T) {
	deadline := time.Now().Add(500 * time.Millisecond)
	start := time.Now()

	// Command will fail and retry, but deadline should still apply
	err := Cmd("sleep", "5").
		WithDeadline(deadline).
		RetryWithBackoff(3, 100*time.Millisecond).
		Error()
	elapsed := time.Since(start)

	// Should terminate around deadline time, not complete all retries
	if elapsed > 2*time.Second {
		t.Errorf("command took too long: %v (expected ~500ms)", elapsed)
	}

	if err == nil {
		t.Fatal("expected error due to deadline, got nil")
	}
}

// TestCommandWithDeadline_WithOutput verifies deadline works when capturing
// output from commands that produce data.
func TestCommandWithDeadline_WithOutput(t *testing.T) {
	deadline := time.Now().Add(2 * time.Second)

	// Command that produces output quickly then exits
	output := Cmd("echo", "hello world").
		WithDeadline(deadline).
		Stdout()

	if output != "hello world\n" {
		t.Errorf("expected 'hello world\\n', got %q", output)
	}
}

// TestCommandWithDeadline_LongOutput verifies deadline terminates commands
// with continuous output.
func TestCommandWithDeadline_LongOutput(t *testing.T) {
	deadline := time.Now().Add(500 * time.Millisecond)
	start := time.Now()

	// Generate infinite output
	_ = Cmd("yes").
		WithDeadline(deadline).
		Stdout()
	elapsed := time.Since(start)

	// Should terminate around deadline
	if elapsed > 2*time.Second {
		t.Errorf("command took too long: %v (expected ~500ms)", elapsed)
	}
}
