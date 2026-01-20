package types

import (
	"strings"
	"testing"

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
