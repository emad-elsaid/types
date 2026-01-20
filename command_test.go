package types

import (
	"context"
	"os"
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
