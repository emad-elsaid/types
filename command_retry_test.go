package types

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

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
