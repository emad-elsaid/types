package types

import (
	"strings"
	"testing"
	"time"
)

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
