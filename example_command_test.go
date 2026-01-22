package types

import (
	"fmt"
	"strings"
	"time"
)

func ExampleCmd() {
	output := Cmd("echo", "hello world").Stdout()
	fmt.Print(output)
	// Output: hello world
}

func ExampleCmd_pipe() {
	result := Cmd("echo", "apple\nbanana\napricot").
		Pipe("grep", "a").
		Stdout()
	fmt.Print(result)
	// Output: apple
	// banana
	// apricot
}

func ExampleCmd_input() {
	result := Cmd("grep", "hello").
		Input("hello world\ngoodbye\nhello again").
		Stdout()
	fmt.Print(result)
	// Output: hello world
	// hello again
}

func ExampleCmdFn() {
	upperCase := CmdFn(func(stdin string) (string, string, error) {
		return strings.ToUpper(stdin), "", nil
	})

	result := upperCase.Input("hello world").Stdout()
	fmt.Print(result)
	// Output: HELLO WORLD
}

func ExampleCommand_Pipe() {
	result := Cmd("echo", "one\ntwo\nthree").
		Pipe("wc", "-l").
		Stdout()
	fmt.Print(strings.TrimSpace(result))
	// Output: 3
}

func ExampleCommand_PipeFn() {
	result := Cmd("echo", "hello").
		PipeFn(func(stdin string) (string, string, error) {
			return strings.ToUpper(stdin), "", nil
		}).
		Stdout()
	fmt.Print(result)
	// Output: HELLO
}

func ExampleCommand_StdoutErr() {
	output, err := Cmd("echo", "success").StdoutErr()
	fmt.Printf("Output: %s, Error: %v", strings.TrimSpace(output), err)
	// Output: Output: success, Error: <nil>
}

func ExampleCommand_WithTimeout() {
	// Command that completes within timeout
	output := Cmd("echo", "hello").WithTimeout(5 * time.Second).Stdout()
	fmt.Print(output)
	// Output: hello
}

func ExampleCommand_Dir() {
	// Run command in a specific directory
	output := Cmd("pwd").Dir("/tmp").Stdout()
	fmt.Print(strings.TrimSpace(output))
	// Output: /tmp
}

func ExampleCommand_Env() {
	// Set environment variables for the command
	output := Cmd("sh", "-c", "echo $MY_VAR").Env("MY_VAR", "hello").Stdout()
	fmt.Print(output)
	// Output: hello
}

func ExampleCommand_ExitCode() {
	// Get the exit code of a failed command
	cmd := Cmd("sh", "-c", "exit 42").Run()
	fmt.Printf("Exit code: %d", cmd.ExitCode())
	// Output: Exit code: 42
}

func ExampleCommand_StdoutTrimmed() {
	// Get stdout without trailing newline
	version := Cmd("echo", "v1.2.3").StdoutTrimmed()
	fmt.Printf("Version: %s", version)
	// Output: Version: v1.2.3
}

func ExampleCommand_String() {
	// Get string representation of a command
	cmd := Cmd("git", "commit", "-m", "message")
	fmt.Println(cmd.String())
	// Output: git commit -m message
}
