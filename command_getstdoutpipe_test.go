package types

import (
	"errors"
	"io"
	"strings"
	"testing"
)

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
