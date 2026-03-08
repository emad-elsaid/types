package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
