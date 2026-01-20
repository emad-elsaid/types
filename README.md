# Types

[![Go Report Card](https://goreportcard.com/badge/github.com/emad-elsaid/types)](https://goreportcard.com/report/github.com/emad-elsaid/types)
[![GoDoc](https://godoc.org/github.com/emad-elsaid/types?status.svg)](https://godoc.org/github.com/emad-elsaid/types)
[![codecov](https://codecov.io/gh/emad-elsaid/types/branch/master/graph/badge.svg)](https://codecov.io/gh/emad-elsaid/types)

Go implementation for generic types, imitating Ruby types

# Data structures

## Slice

A slice of `comparable` type that can hold any value

### Example

```go
func ExampleSlice() {
	a := Slice[int]{1, 2, 3, 4, 5, 6}

	// multiply every element by 100
	a = a.Map(func(e int) int {
		return e * 100
	})

	// select any element <= 300
	a = a.KeepIf(func(e int) bool {
		return e <= 300
	})

	fmt.Print(a)

	// Output: [100 200 300]
}
```

### Slice Methods available:

```go
func (a Slice[T]) All(block func(T) bool) bool
func (a Slice[T]) Any(block func(T) bool) bool
func (a Slice[T]) At(index int) *T
func (a Slice[T]) CountBy(block func(T) bool) (count int)
func (a Slice[T]) CountElement(element T) (count int)
func (a Slice[T]) Cycle(count int, block func(T))
func (a Slice[T]) Delete(element T) Slice[T]
func (a Slice[T]) DeleteAt(index int) Slice[T]
func (a Slice[T]) DeleteIf(block func(T) bool) Slice[T]
func (a Slice[T]) Drop(count int) Slice[T]
func (a Slice[T]) Each(block func(T))
func (a Slice[T]) EachIndex(block func(int))
func (a Slice[T]) Fetch(index int, defaultValue T) T
func (a Slice[T]) Fill(element T, start int, length int) Slice[T]
func (a Slice[T]) FillWith(start int, length int, block func(int) T) Slice[T]
func (a Slice[T]) First() *T
func (a Slice[T]) Firsts(count int) Slice[T]
func (a Slice[T]) Include(element T) bool
func (a Slice[T]) Index(element T) int
func (a Slice[T]) IndexBy(block func(T) bool) int
func (a Slice[T]) Insert(index int, elements ...T) Slice[T]
func (a Slice[T]) IsEmpty() bool
func (a Slice[T]) IsEq(other Slice[T]) bool
func (a Slice[T]) KeepIf(block func(T) bool) Slice[T]
func (a Slice[T]) Last() *T
func (a Slice[T]) Lasts(count int) Slice[T]
func (a Slice[T]) Len() int
func (a Slice[T]) Map(block func(T) T) Slice[T]
func (a Slice[T]) Max(block func(T) int) T
func (a Slice[T]) Min(block func(T) int) T
func (a Slice[T]) Partition(predicate func(T) bool) (Slice[T], Slice[T])
func (a Slice[T]) Pop() (Slice[T], T)
func (a Slice[T]) Push(element T) Slice[T]
func (a Slice[T]) Reduce(block func(T) bool) Slice[T]
func (a Slice[T]) Reverse() Slice[T]
func (a Slice[T]) Select(block func(T) bool) Slice[T]
func (a Slice[T]) SelectUntil(block func(T) bool) Slice[T]
func (a Slice[T]) Shift() (T, Slice[T])
func (a Slice[T]) Shuffle() Slice[T]
func (a Slice[T]) Unique() Slice[T]
func (a Slice[T]) Unshift(element T) Slice[T]
func SliceReduce[T comparable, U any](s Slice[T], initial U, fn func(U, T) U) U
```

## Set

An order-preserving set of `comparable` type that stores unique elements

### Example

```go
func ExampleSet() {
	s := NewSet(1, 2, 3, 3, 4, 5)

	// Set automatically deduplicates
	fmt.Println(s.Size()) // 5

	// Add elements
	s.Add(6)
	s.Add(3) // returns false, already exists

	// Set operations
	other := NewSet(4, 5, 6, 7)
	union := s.Union(other)
	intersection := s.Intersection(other)

	// Functional operations
	filtered := s.Filter(func(x int) bool {
		return x > 3
	})

	doubled := SetMap(s, func(x int) int {
		return x * 2
	})

	fmt.Print(filtered)

	// Output: Set{4, 5, 6}
}
```

### Set Methods available:

```go
func NewSet[T comparable](slice ...T) *Set[T]
func (s *Set[T]) Add(item T) bool
func (s *Set[T]) All(predicate func(T) bool) bool
func (s *Set[T]) Any(predicate func(T) bool) bool
func (s *Set[T]) Clear()
func (s *Set[T]) Clone() *Set[T]
func (s *Set[T]) Contains(item T) bool
func (s *Set[T]) Count(predicate func(T) bool) int
func (s *Set[T]) Difference(other *Set[T]) *Set[T]
func (s *Set[T]) Drop(n int) *Set[T]
func (s *Set[T]) Each(fn func(T))
func (s *Set[T]) Equal(other *Set[T]) bool
func (s *Set[T]) Filter(predicate func(T) bool) *Set[T]
func (s *Set[T]) Find(predicate func(T) bool) (T, bool)
func (s *Set[T]) Intersection(other *Set[T]) *Set[T]
func (s *Set[T]) IsDisjoint(other *Set[T]) bool
func (s *Set[T]) IsEmpty() bool
func (s *Set[T]) IsSubset(other *Set[T]) bool
func (s *Set[T]) IsSuperset(other *Set[T]) bool
func (s *Set[T]) None(predicate func(T) bool) bool
func (s *Set[T]) Partition(predicate func(T) bool) (*Set[T], *Set[T])
func (s *Set[T]) Reject(predicate func(T) bool) *Set[T]
func (s *Set[T]) Remove(item T) bool
func (s *Set[T]) Size() int
func (s *Set[T]) String() string
func (s *Set[T]) SymmetricDifference(other *Set[T]) *Set[T]
func (s *Set[T]) Take(n int) *Set[T]
func (s *Set[T]) ToSlice() []T
func (s *Set[T]) Union(other *Set[T]) *Set[T]
func SetMap[T, U comparable](s *Set[T], fn func(T) U) *Set[U]
func SetReduce[T comparable, U any](s *Set[T], initial U, fn func(U, T) U) U
```

## Map

A thread-safe generic wrapper around `sync.Map` with type-safe key-value operations

### Example

```go
func ExampleMap() {
	m := Map[string, int]{}

	// Store values
	m.Store("one", 1)
	m.Store("two", 2)
	m.Store("three", 3)

	// Load values
	if val, ok := m.Load("two"); ok {
		fmt.Println(val) // 2
	}

	// Iterate over all entries
	m.Range(func(k string, v int) bool {
		fmt.Printf("%s: %d\n", k, v)
		return true // continue iteration
	})

	// LoadOrStore
	actual, loaded := m.LoadOrStore("four", 4)
	fmt.Println(loaded) // false (newly stored)

	// Delete
	m.Delete("one")
}
```

### Map Methods available:

```go
func (m *Map[K, V]) Delete(k K)
func (m *Map[K, V]) Load(k K) (v V, ok bool)
func (m *Map[K, V]) LoadAndDelete(k K) (v V, loaded bool)
func (m *Map[K, V]) LoadOrStore(k K, v V) (actual V, loaded bool)
func (m *Map[K, V]) Range(f func(k K, v V) bool)
func (m *Map[K, V]) Store(k K, v V)
func (m *Map[K, V]) Swap(k K, v V) (previous V, loaded bool)
```

## Chan

Generic channel utility functions for functional-style channel operations

### Example

```go
func ExampleChan() {
	// Create input channel
	input := make(chan int)
	go func() {
		defer close(input)
		for i := 1; i <= 10; i++ {
			input <- i
		}
	}()

	// Map: transform each value
	doubled := ChanMap(input, func(x int) int {
		return x * 2
	})

	// Filter: keep only values > 10
	filtered := ChanFilter(doubled, func(x int) bool {
		return x > 10
	})

	// Each: consume and print
	ChanEach(filtered, func(x int) {
		fmt.Println(x)
	})
	// Output: 12 14 16 18 20
}
```

### Chan Functions available:

```go
func OrderedParallelizeChan[In, Out any](input <-chan In, workers int, process func(<-chan In) <-chan Out) <-chan Out
func ChanMap[In, Out any](input <-chan In, processor func(In) Out) <-chan Out
func ChanFilter[T any](input <-chan T, filter func(T) bool) <-chan T
func ChanEach[T any](input <-chan T, fns ...func(T))
```

## Command

A fluent interface for executing system commands with support for piping, chaining, and custom function transformations.

### Features

- **Command Chaining**: Chain commands together with `Pipe`
- **Function Transformations**: Inject Go functions into pipelines with `PipeFn` and `CmdFn`
- **Sudo Support**: Run commands with sudo privileges
- **Interactive Mode**: Connect commands directly to terminal for user input
- **Input Redirection**: Provide stdin from strings or io.Reader
- **Context Support**: Cancel or timeout commands with context
- **Working Directory**: Set the directory where commands execute
- **Environment Variables**: Configure command environment
- **Exit Code Access**: Get command exit codes
- **Retry Logic**: Retry failed commands with optional backoff
- **Lazy Execution**: Commands execute only when output is requested
- **Idempotent**: Multiple calls to output methods return cached results

### Example

```go
func ExampleCommand() {
	// Simple command execution
	output := Cmd("echo", "hello world").Stdout()
	fmt.Println(output) // "hello world\n"

	// Piping commands together
	result := Cmd("echo", "apple\nbanana\napricot").
		Pipe("grep", "a").
		Pipe("wc", "-l").
		Stdout()
	fmt.Println(strings.TrimSpace(result)) // "3"

	// Using custom functions in pipelines
	upperResult := Cmd("echo", "hello").
		PipeFn(func(stdin string) (string, string, error) {
			return strings.ToUpper(stdin), "", nil
		}).
		Stdout()
	fmt.Println(upperResult) // "HELLO\n"

	// Input redirection
	grepResult := Cmd("grep", "hello").
		Input("hello world\ngoodbye\nhello again").
		Stdout()
	fmt.Println(grepResult) // "hello world\nhello again\n"

	// Running with sudo (requires authentication)
	err := Cmd("systemctl", "restart", "nginx").Sudo().Error()
	if err != nil {
		// handle error
	}

	// Interactive mode for commands requiring user input
	err = Cmd("vim", "file.txt").Interactive().Error()
	
	// Timeout for long-running commands
	output = Cmd("sleep", "10").WithTimeout(2*time.Second).Stdout()
	// Returns with context deadline exceeded error
	
	// Set working directory
	output = Cmd("ls").Dir("/tmp").Stdout()
	
	// Set environment variables
	output = Cmd("printenv", "MY_VAR").
		Env("MY_VAR", "hello").
		Env("ANOTHER", "world").
		Stdout()
	
	// Get exit code
	cmd := Cmd("false").Run()
	exitCode := cmd.ExitCode() // 1
	
	// Retry with backoff
	output = Cmd("curl", "http://example.com").
		RetryWithBackoff(3, 2*time.Second).
		Stdout()
	
	// Get trimmed output (no trailing whitespace)
	version := Cmd("git", "--version").StdoutTrimmed()
}
```

### Command Functions and Methods

```go
// Creating commands
func Cmd(cmd string, args ...string) *Command
func CmdFn(fn func(stdin string) (stdout, stderr string, err error)) *Command
func Sudo(cmd string, args ...string) *Command

// Chaining and piping
func (c *Command) Pipe(cmd string, args ...string) *Command
func (c *Command) PipeFn(fn func(stdin string) (stdout, stderr string, err error)) *Command

// Configuration
func (c *Command) Interactive() *Command
func (c *Command) Input(input string) *Command
func (c *Command) InputReader(r io.Reader) *Command
func (c *Command) Sudo() *Command
func (c *Command) Dir(path string) *Command
func (c *Command) Env(key, value string) *Command
func (c *Command) EnvMap(env map[string]string) *Command
func (c *Command) ClearEnv() *Command

// Context and timeout
func (c *Command) WithContext(ctx context.Context) *Command
func (c *Command) WithTimeout(duration time.Duration) *Command
func (c *Command) WithDeadline(t time.Time) *Command

// Retry logic
func (c *Command) Retry(attempts int) *Command
func (c *Command) RetryWithBackoff(attempts int, delay time.Duration) *Command

// Execution and output
func (c *Command) Run() *Command
func (c *Command) Stdout() string
func (c *Command) StdoutTrimmed() string
func (c *Command) Stderr() string
func (c *Command) Error() error
func (c *Command) ExitCode() int
func (c *Command) StdoutErr() (string, error)
func (c *Command) StderrErr() (string, error)
func (c *Command) StdoutStderr() string

// Utility
func (c *Command) String() string
```
