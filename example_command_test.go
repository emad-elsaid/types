package types

import (
	"fmt"
	"strings"
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
