package types_test

import (
	"fmt"

	"github.com/emad-elsaid/types"
)

func ExampleMap_Merge() {
	// Create first map with default settings
	defaults := types.Map[string, string]{}
	defaults.Store("host", "localhost")
	defaults.Store("port", "8080")
	defaults.Store("timeout", "30")

	// Create second map with user overrides
	overrides := types.Map[string, string]{}
	overrides.Store("port", "3000")
	overrides.Store("debug", "true")

	// Merge: overrides take precedence over defaults
	config := defaults.Merge(&overrides)

	// Print configuration values
	if v, ok := config.Load("host"); ok {
		fmt.Println("host:", v)
	}
	if v, ok := config.Load("port"); ok {
		fmt.Println("port:", v)
	}
	if v, ok := config.Load("timeout"); ok {
		fmt.Println("timeout:", v)
	}
	if v, ok := config.Load("debug"); ok {
		fmt.Println("debug:", v)
	}

	// Output:
	// host: localhost
	// port: 3000
	// timeout: 30
	// debug: true
}
