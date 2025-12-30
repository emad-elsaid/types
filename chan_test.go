package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Helper to create a simple process function
func makeProcessor[In, Out any](fn func(In) Out) func(<-chan In) <-chan Out {
	return func(in <-chan In) <-chan Out {
		out := make(chan Out)
		go func() {
			defer close(out)
			for v := range in {
				out <- fn(v)
			}
		}()
		return out
	}
}

func TestOrderedParallelizeChan_BasicOrder(t *testing.T) {
	tests := []struct {
		name    string
		input   []int
		workers int
		process func(<-chan int) <-chan int
		want    []int
	}{
		{
			name:    "simple sequence",
			input:   []int{1, 2, 3, 4, 5},
			workers: 3,
			process: makeProcessor(func(x int) int { return x * 2 }),
			want:    []int{2, 4, 6, 8, 10},
		},
		{
			name:    "single worker",
			input:   []int{1, 2, 3},
			workers: 1,
			process: makeProcessor(func(x int) int { return x * 2 }),
			want:    []int{2, 4, 6},
		},
		{
			name:    "many workers",
			input:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			workers: 10,
			process: makeProcessor(func(x int) int { return x * 3 }),
			want:    []int{3, 6, 9, 12, 15, 18, 21, 24, 27, 30},
		},
		{
			name:    "zero workers defaults to one",
			input:   []int{1, 2, 3},
			workers: 0,
			process: makeProcessor(func(x int) int { return x + 1 }),
			want:    []int{2, 3, 4},
		},
		{
			name:    "negative workers defaults to one",
			input:   []int{1, 2, 3},
			workers: -5,
			process: makeProcessor(func(x int) int { return x + 1 }),
			want:    []int{2, 3, 4},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := make(chan int)
			go func() {
				defer close(input)
				for _, v := range tc.input {
					input <- v
				}
			}()

			output := OrderedParallelizeChan(input, tc.workers, tc.process)

			result := []int{}
			for v := range output {
				result = append(result, v)
			}

			require.Equal(t, tc.want, result)
		})
	}
}

func TestOrderedParallelizeChan_NilInput(t *testing.T) {
	var input chan int = nil

	output := OrderedParallelizeChan(input, 3, makeProcessor(func(x int) int { return x * 2 }))

	require.Nil(t, output)
}

func TestOrderedParallelizeChan_EmptyInput(t *testing.T) {
	input := make(chan int)
	close(input)

	output := OrderedParallelizeChan(input, 3, makeProcessor(func(x int) int { return x * 2 }))

	result := []int{}
	for v := range output {
		result = append(result, v)
	}

	require.Empty(t, result)
}

func TestOrderedParallelizeChan_VaryingProcessingTime(t *testing.T) {
	input := make(chan int)
	go func() {
		defer close(input)
		for i := 1; i <= 10; i++ {
			input <- i
		}
	}()

	// Process function where later items finish faster to test order preservation
	process := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for x := range in {
				time.Sleep(time.Millisecond * time.Duration(10-x))
				out <- x * 2
			}
		}()
		return out
	}

	output := OrderedParallelizeChan(input, 5, process)

	result := []int{}
	for v := range output {
		result = append(result, v)
	}

	expected := []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}
	require.Equal(t, expected, result)
}

func TestOrderedParallelizeChan_StringType(t *testing.T) {
	input := make(chan string)
	go func() {
		defer close(input)
		words := []string{"hello", "world", "test"}
		for _, w := range words {
			input <- w
		}
	}()

	output := OrderedParallelizeChan(input, 2, makeProcessor(func(s string) string {
		return s + "!"
	}))

	result := []string{}
	for v := range output {
		result = append(result, v)
	}

	expected := []string{"hello!", "world!", "test!"}
	require.Equal(t, expected, result)
}

func TestOrderedParallelizeChan_TypeConversion(t *testing.T) {
	input := make(chan int)
	go func() {
		defer close(input)
		for i := 1; i <= 5; i++ {
			input <- i
		}
	}()

	output := OrderedParallelizeChan(input, 3, makeProcessor(func(x int) string {
		return string(rune('A' + x - 1))
	}))

	result := []string{}
	for v := range output {
		result = append(result, v)
	}

	expected := []string{"A", "B", "C", "D", "E"}
	require.Equal(t, expected, result)
}

func TestOrderedParallelizeChan_StructType(t *testing.T) {
	type Item struct {
		ID    int
		Value string
	}

	input := make(chan Item)
	go func() {
		defer close(input)
		for i := 1; i <= 3; i++ {
			input <- Item{ID: i, Value: "test"}
		}
	}()

	output := OrderedParallelizeChan(input, 2, makeProcessor(func(item Item) Item {
		return Item{ID: item.ID * 10, Value: item.Value + "!"}
	}))

	result := []Item{}
	for v := range output {
		result = append(result, v)
	}

	expected := []Item{
		{ID: 10, Value: "test!"},
		{ID: 20, Value: "test!"},
		{ID: 30, Value: "test!"},
	}
	require.Equal(t, expected, result)
}

func TestOrderedParallelizeChan_LargeDataset(t *testing.T) {
	count := 1000
	input := make(chan int)
	go func() {
		defer close(input)
		for i := 0; i < count; i++ {
			input <- i
		}
	}()

	output := OrderedParallelizeChan(input, 20, makeProcessor(func(x int) int {
		return x * 2
	}))

	result := []int{}
	for v := range output {
		result = append(result, v)
	}

	require.Len(t, result, count)
	for i := 0; i < count; i++ {
		require.Equal(t, i*2, result[i])
	}
}

func TestOrderedParallelizeChan_SingleItem(t *testing.T) {
	input := make(chan int)
	go func() {
		defer close(input)
		input <- 42
	}()

	output := OrderedParallelizeChan(input, 5, makeProcessor(func(x int) int {
		return x * 2
	}))

	result := []int{}
	for v := range output {
		result = append(result, v)
	}

	require.Equal(t, []int{84}, result)
}

func TestOrderedParallelizeChan_WorkerReceivesMultipleItems(t *testing.T) {
	input := make(chan int)
	go func() {
		defer close(input)
		for i := 1; i <= 9; i++ {
			input <- i
		}
	}()

	// With 3 workers, each should receive 3 items
	// Worker 0: 1, 4, 7
	// Worker 1: 2, 5, 8
	// Worker 2: 3, 6, 9
	output := OrderedParallelizeChan(input, 3, makeProcessor(func(x int) int {
		return x * 10
	}))

	result := []int{}
	for v := range output {
		result = append(result, v)
	}

	expected := []int{10, 20, 30, 40, 50, 60, 70, 80, 90}
	require.Equal(t, expected, result)
}

func TestOrderedParallelizeChan_CustomChannelProcessor(t *testing.T) {
	input := make(chan int)
	go func() {
		defer close(input)
		for i := 1; i <= 6; i++ {
			input <- i
		}
	}()

	// Custom processor that batches items
	batchProcessor := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			sum := 0
			count := 0
			for v := range in {
				sum += v
				count++
			}
			if count > 0 {
				out <- sum
			}
		}()
		return out
	}

	// With 3 workers processing 6 items:
	// Worker 0 gets: 1, 4 -> sum=5
	// Worker 1 gets: 2, 5 -> sum=7
	// Worker 2 gets: 3, 6 -> sum=9
	output := OrderedParallelizeChan(input, 3, batchProcessor)

	result := []int{}
	for v := range output {
		result = append(result, v)
	}

	expected := []int{5, 7, 9}
	require.Equal(t, expected, result)
}
