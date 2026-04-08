package types

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestConvention_TestFilesHaveImplementation enforces the convention that every test file
// (except example_*_test.go and this init_test.go) must have a corresponding implementation file.
//
// Convention: For every *_test.go file, there must be a *.go implementation file with the same base name.
// Example: slice_test.go requires slice.go to exist
//
// This ensures test files don't become orphaned and helps maintain a clear structure.
func TestConvention_TestFilesHaveImplementation(t *testing.T) {
	// Get all test files in current directory
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("Failed to list Go files: %v", err)
	}

	var violations []string

	for _, file := range files {
		// Skip non-test files
		if !strings.HasSuffix(file, "_test.go") {
			continue
		}

		// Skip example tests (they're documentation, not unit tests)
		if strings.HasPrefix(file, "example_") {
			continue
		}

		// Skip this init test file itself
		if file == "init_test.go" {
			continue
		}

		// Determine expected implementation file
		implFile := strings.TrimSuffix(file, "_test.go") + ".go"

		// Check if implementation file exists
		if _, err := os.Stat(implFile); os.IsNotExist(err) {
			violations = append(violations, file)
		}
	}

	if len(violations) > 0 {
		t.Errorf("CONVENTION VIOLATION: The following test files do not have corresponding implementation files:\n\n"+
			"Files missing implementations:\n")
		for _, file := range violations {
			expectedImpl := strings.TrimSuffix(file, "_test.go") + ".go"
			t.Errorf("  - %s (expected: %s)\n", file, expectedImpl)
		}
		t.Errorf("\nCONVENTION: Every *_test.go file must have a corresponding *.go implementation file.\n" +
			"This ensures test files are properly paired with their implementation and prevents orphaned tests.\n")
	}
}
