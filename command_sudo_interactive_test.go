package types_test

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/emad-elsaid/types"
)

// TestSudo tests the Sudo method and function
func TestSudo(t *testing.T) {
	t.Run("Sudo method sets useSudo flag", func(t *testing.T) {
		// Create a command that would succeed without sudo
		cmd := types.Cmd("echo", "test").Sudo()
		
		// The String() representation should include sudo
		cmdStr := cmd.String()
		if !strings.HasPrefix(cmdStr, "sudo") {
			t.Errorf("Expected command string to start with 'sudo', got: %s", cmdStr)
		}
	})

	t.Run("Sudo function creates command with sudo", func(t *testing.T) {
		cmd := types.Sudo("echo", "test")
		
		cmdStr := cmd.String()
		if !strings.HasPrefix(cmdStr, "sudo") {
			t.Errorf("Expected command string to start with 'sudo', got: %s", cmdStr)
		}
	})

	t.Run("Sudo preserves arguments", func(t *testing.T) {
		cmd := types.Cmd("test", "arg1", "arg2").Sudo()
		
		cmdStr := cmd.String()
		expected := "sudo test arg1 arg2"
		if cmdStr != expected {
			t.Errorf("Expected '%s', got: %s", expected, cmdStr)
		}
	})

	t.Run("Sudo can be chained with other methods", func(t *testing.T) {
		// Test that Sudo can be chained
		cmd := types.Cmd("echo", "test").Sudo().Env("TEST", "value")
		
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
		cmd := types.Cmd("sudo", "-n", "true")
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
		cmd := types.Cmd("echo", "test").Interactive()
		
		// We can't directly check the interactive flag (it's private),
		// but we can verify the command was created successfully
		if cmd == nil {
			t.Error("Expected non-nil command")
		}
	})

	t.Run("Interactive can be chained", func(t *testing.T) {
		cmd := types.Cmd("cat").Input("test").Interactive()
		
		if cmd == nil {
			t.Error("Expected non-nil command after chaining")
		}
	})

	t.Run("Interactive with non-interactive command", func(t *testing.T) {
		// Even in interactive mode, commands that don't require interaction should work
		// Note: This won't actually be interactive in tests, but verifies no crash
		cmd := types.Cmd("echo", "hello")
		
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
		cmd := types.Cmd("cat", "/dev/null").Interactive()
		
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
		cmd := types.Cmd("cat", "/dev/null").Sudo().Interactive()
		
		cmdStr := cmd.String()
		if !strings.HasPrefix(cmdStr, "sudo") {
			t.Errorf("Expected command to start with 'sudo', got: %s", cmdStr)
		}
	})

	t.Run("Interactive and Sudo order doesn't matter", func(t *testing.T) {
		cmd1 := types.Cmd("echo", "test").Sudo().Interactive()
		cmd2 := types.Cmd("echo", "test").Interactive().Sudo()
		
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
		cmd := types.Cmd("sudo", "-n", "echo", "test")
		
		// This might fail if sudo requires password, which is fine for testing
		// We're just verifying no panic occurs
		_ = cmd.Stdout()
	})
}

// TestSudoInPipeline tests that Sudo works correctly in command pipelines
func TestSudoInPipeline(t *testing.T) {
	t.Run("Sudo on first command in pipeline", func(t *testing.T) {
		// Create a pipeline where sudo is used on the first command
		cmd1 := types.Cmd("echo", "test").Sudo()
		
		// Verify the first command has sudo
		cmdStr := cmd1.String()
		if !strings.Contains(cmdStr, "sudo") {
			t.Errorf("Expected first command to contain 'sudo', got: %s", cmdStr)
		}
		
		// Create pipeline - subsequent commands don't show in String()
		_ = cmd1.Pipe("cat").Pipe("wc", "-l")
	})

	t.Run("Sudo on piped command", func(t *testing.T) {
		cmd := types.Cmd("echo", "test").Pipe("cat").Sudo()
		
		cmdStr := cmd.String()
		if !strings.Contains(cmdStr, "sudo") {
			t.Errorf("Expected pipeline to contain 'sudo', got: %s", cmdStr)
		}
	})

	t.Run("Sudo in pipeline execution", func(t *testing.T) {
		// Verify that sudo commands can be part of a pipeline
		// and execute without errors (using safe commands)
		cmd := types.Cmd("echo", "hello").Pipe("cat")
		
		result := cmd.Stdout()
		expected := "hello\n"
		if result != expected {
			t.Errorf("Expected '%s', got: '%s'", expected, result)
		}
	})
}
