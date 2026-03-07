package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

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
