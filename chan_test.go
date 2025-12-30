package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOrderedParallelizeChan(t *testing.T) {
	type Item struct {
		ID    int
		Value string
	}

	makeProcessor := func(fn func(int) int) func(<-chan int) <-chan int {
		return func(in <-chan int) <-chan int {
			out := make(chan int)
			go func() {
				defer close(out)
				for v := range in {
					out <- fn(v)
				}
			}()
			return out
		}
	}

	tests := []struct {
		name       string
		setupInput func() any
		workers    int
		processor  any
		want       any
	}{
		{
			name: "simple sequence",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					for _, v := range []int{1, 2, 3, 4, 5} {
						ch <- v
					}
				}()
				return ch
			},
			workers:   3,
			processor: makeProcessor(func(x int) int { return x * 2 }),
			want:      []int{2, 4, 6, 8, 10},
		},
		{
			name: "single worker",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					for _, v := range []int{1, 2, 3} {
						ch <- v
					}
				}()
				return ch
			},
			workers:   1,
			processor: makeProcessor(func(x int) int { return x * 2 }),
			want:      []int{2, 4, 6},
		},
		{
			name: "many workers",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					for _, v := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
						ch <- v
					}
				}()
				return ch
			},
			workers:   10,
			processor: makeProcessor(func(x int) int { return x * 3 }),
			want:      []int{3, 6, 9, 12, 15, 18, 21, 24, 27, 30},
		},
		{
			name: "zero workers defaults to one",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					for _, v := range []int{1, 2, 3} {
						ch <- v
					}
				}()
				return ch
			},
			workers:   0,
			processor: makeProcessor(func(x int) int { return x + 1 }),
			want:      []int{2, 3, 4},
		},
		{
			name: "negative workers defaults to one",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					for _, v := range []int{1, 2, 3} {
						ch <- v
					}
				}()
				return ch
			},
			workers:   -5,
			processor: makeProcessor(func(x int) int { return x + 1 }),
			want:      []int{2, 3, 4},
		},
		{
			name: "nil input",
			setupInput: func() any {
				var ch chan int = nil
				return ch
			},
			workers:   3,
			processor: makeProcessor(func(x int) int { return x * 2 }),
			want:      nil,
		},
		{
			name: "empty input",
			setupInput: func() any {
				ch := make(chan int)
				close(ch)
				return ch
			},
			workers:   3,
			processor: makeProcessor(func(x int) int { return x * 2 }),
			want:      []int{},
		},
		{
			name: "varying processing time preserves order",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					for i := 1; i <= 10; i++ {
						ch <- i
					}
				}()
				return ch
			},
			workers: 5,
			processor: func(in <-chan int) <-chan int {
				out := make(chan int)
				go func() {
					defer close(out)
					for x := range in {
						time.Sleep(time.Millisecond * time.Duration(10-x))
						out <- x * 2
					}
				}()
				return out
			},
			want: []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20},
		},
		{
			name: "string type",
			setupInput: func() any {
				ch := make(chan string)
				go func() {
					defer close(ch)
					for _, w := range []string{"hello", "world", "test"} {
						ch <- w
					}
				}()
				return ch
			},
			workers: 2,
			processor: func(in <-chan string) <-chan string {
				out := make(chan string)
				go func() {
					defer close(out)
					for v := range in {
						out <- v + "!"
					}
				}()
				return out
			},
			want: []string{"hello!", "world!", "test!"},
		},
		{
			name: "type conversion int to string",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					for i := 1; i <= 5; i++ {
						ch <- i
					}
				}()
				return ch
			},
			workers: 3,
			processor: func(in <-chan int) <-chan string {
				out := make(chan string)
				go func() {
					defer close(out)
					for v := range in {
						out <- string(rune('A' + v - 1))
					}
				}()
				return out
			},
			want: []string{"A", "B", "C", "D", "E"},
		},
		{
			name: "struct type",
			setupInput: func() any {
				ch := make(chan Item)
				go func() {
					defer close(ch)
					for i := 1; i <= 3; i++ {
						ch <- Item{ID: i, Value: "test"}
					}
				}()
				return ch
			},
			workers: 2,
			processor: func(in <-chan Item) <-chan Item {
				out := make(chan Item)
				go func() {
					defer close(out)
					for v := range in {
						out <- Item{ID: v.ID * 10, Value: v.Value + "!"}
					}
				}()
				return out
			},
			want: []Item{
				{ID: 10, Value: "test!"},
				{ID: 20, Value: "test!"},
				{ID: 30, Value: "test!"},
			},
		},
		{
			name: "large dataset",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					for i := 0; i < 100; i++ {
						ch <- i
					}
				}()
				return ch
			},
			workers:   20,
			processor: makeProcessor(func(x int) int { return x * 2 }),
			want: func() []int {
				result := make([]int, 100)
				for i := 0; i < 100; i++ {
					result[i] = i * 2
				}
				return result
			}(),
		},
		{
			name: "single item",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					ch <- 42
				}()
				return ch
			},
			workers:   5,
			processor: makeProcessor(func(x int) int { return x * 2 }),
			want:      []int{84},
		},
		{
			name: "worker receives multiple items round robin",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					for i := 1; i <= 9; i++ {
						ch <- i
					}
				}()
				return ch
			},
			workers:   3,
			processor: makeProcessor(func(x int) int { return x * 10 }),
			want:      []int{10, 20, 30, 40, 50, 60, 70, 80, 90},
		},
		{
			name: "custom batch processor",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					for i := 1; i <= 6; i++ {
						ch <- i
					}
				}()
				return ch
			},
			workers: 3,
			processor: func(in <-chan int) <-chan int {
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
			},
			want: []int{5, 7, 9},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := tc.setupInput()

			switch in := input.(type) {
			case chan int:
				switch proc := tc.processor.(type) {
				case func(<-chan int) <-chan int:
					output := OrderedParallelizeChan(in, tc.workers, proc)

					if tc.want == nil {
						require.Nil(t, output)
						return
					}

					var result []int
					for v := range output {
						result = append(result, v)
					}

					want := tc.want.([]int)
					if len(want) == 0 && len(result) == 0 {
						return
					}
					require.Equal(t, want, result)

				case func(<-chan int) <-chan string:
					output := OrderedParallelizeChan(in, tc.workers, proc)

					var result []string
					for v := range output {
						result = append(result, v)
					}
					require.Equal(t, tc.want, result)
				}

			case chan string:
				proc := tc.processor.(func(<-chan string) <-chan string)
				output := OrderedParallelizeChan(in, tc.workers, proc)

				var result []string
				for v := range output {
					result = append(result, v)
				}
				require.Equal(t, tc.want, result)

			case chan Item:
				proc := tc.processor.(func(<-chan Item) <-chan Item)
				output := OrderedParallelizeChan(in, tc.workers, proc)

				var result []Item
				for v := range output {
					result = append(result, v)
				}
				require.Equal(t, tc.want, result)
			}
		})
	}
}

func TestChanProcessor(t *testing.T) {
	type Item struct {
		ID    int
		Value string
	}

	tests := []struct {
		name       string
		setupInput func() any
		processor  any
		want       any
		wantCap    int
	}{
		{
			name: "double numbers",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					for _, v := range []int{1, 2, 3, 4, 5} {
						ch <- v
					}
				}()
				return ch
			},
			processor: func(x int) int { return x * 2 },
			want:      []int{2, 4, 6, 8, 10},
		},
		{
			name: "add one",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					for _, v := range []int{10, 20, 30} {
						ch <- v
					}
				}()
				return ch
			},
			processor: func(x int) int { return x + 1 },
			want:      []int{11, 21, 31},
		},
		{
			name: "square numbers",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					for _, v := range []int{2, 3, 4} {
						ch <- v
					}
				}()
				return ch
			},
			processor: func(x int) int { return x * x },
			want:      []int{4, 9, 16},
		},
		{
			name: "nil input",
			setupInput: func() any {
				var ch chan int = nil
				return ch
			},
			processor: func(x int) int { return x * 2 },
			want:      nil,
		},
		{
			name: "empty input",
			setupInput: func() any {
				ch := make(chan int)
				close(ch)
				return ch
			},
			processor: func(x int) int { return x * 2 },
			want:      []int{},
		},
		{
			name: "single item",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					ch <- 42
				}()
				return ch
			},
			processor: func(x int) int { return x * 2 },
			want:      []int{84},
		},
		{
			name: "type conversion int to string",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					for i := 1; i <= 5; i++ {
						ch <- i
					}
				}()
				return ch
			},
			processor: func(x int) string { return string(rune('A' + x - 1)) },
			want:      []string{"A", "B", "C", "D", "E"},
		},
		{
			name: "string processing",
			setupInput: func() any {
				ch := make(chan string)
				go func() {
					defer close(ch)
					for _, w := range []string{"hello", "world", "test"} {
						ch <- w
					}
				}()
				return ch
			},
			processor: func(s string) string { return s + "!" },
			want:      []string{"hello!", "world!", "test!"},
		},
		{
			name: "struct type",
			setupInput: func() any {
				ch := make(chan Item)
				go func() {
					defer close(ch)
					for i := 1; i <= 3; i++ {
						ch <- Item{ID: i, Value: "test"}
					}
				}()
				return ch
			},
			processor: func(item Item) Item {
				return Item{ID: item.ID * 10, Value: item.Value + "!"}
			},
			want: []Item{
				{ID: 10, Value: "test!"},
				{ID: 20, Value: "test!"},
				{ID: 30, Value: "test!"},
			},
		},
		{
			name: "preserves capacity",
			setupInput: func() any {
				ch := make(chan int, 10)
				go func() {
					defer close(ch)
					for i := 1; i <= 5; i++ {
						ch <- i
					}
				}()
				return ch
			},
			processor: func(x int) int { return x * 2 },
			want:      []int{2, 4, 6, 8, 10},
			wantCap:   10,
		},
		{
			name: "large dataset",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					for i := 0; i < 100; i++ {
						ch <- i
					}
				}()
				return ch
			},
			processor: func(x int) int { return x * 2 },
			want: func() []int {
				result := make([]int, 100)
				for i := 0; i < 100; i++ {
					result[i] = i * 2
				}
				return result
			}(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := tc.setupInput()

			switch in := input.(type) {
			case chan int:
				switch proc := tc.processor.(type) {
				case func(int) int:
					output := ChanProcessor(in, proc)

					if tc.want == nil {
						require.Nil(t, output)
						return
					}

					if tc.wantCap > 0 {
						require.Equal(t, tc.wantCap, cap(output))
					}

					var result []int
					for v := range output {
						result = append(result, v)
					}

					want := tc.want.([]int)
					if len(want) == 0 && len(result) == 0 {
						return
					}
					require.Equal(t, want, result)

				case func(int) string:
					output := ChanProcessor(in, proc)

					var result []string
					for v := range output {
						result = append(result, v)
					}
					require.Equal(t, tc.want, result)
				}

			case chan string:
				proc := tc.processor.(func(string) string)
				output := ChanProcessor(in, proc)

				var result []string
				for v := range output {
					result = append(result, v)
				}
				require.Equal(t, tc.want, result)

			case chan Item:
				proc := tc.processor.(func(Item) Item)
				output := ChanProcessor(in, proc)

				var result []Item
				for v := range output {
					result = append(result, v)
				}
				require.Equal(t, tc.want, result)
			}
		})
	}
}

func TestChanFilter(t *testing.T) {
	type Item struct {
		ID    int
		Value string
	}

	tests := []struct {
		name       string
		setupInput func() any
		filter     any
		want       any
		wantCap    int
	}{
		{
			name: "filter even numbers",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					for _, v := range []int{1, 2, 3, 4, 5, 6} {
						ch <- v
					}
				}()
				return ch
			},
			filter: func(x int) bool { return x%2 == 0 },
			want:   []int{2, 4, 6},
		},
		{
			name: "filter odd numbers",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					for _, v := range []int{1, 2, 3, 4, 5, 6} {
						ch <- v
					}
				}()
				return ch
			},
			filter: func(x int) bool { return x%2 != 0 },
			want:   []int{1, 3, 5},
		},
		{
			name: "filter greater than threshold",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					for _, v := range []int{1, 5, 10, 15, 3, 20} {
						ch <- v
					}
				}()
				return ch
			},
			filter: func(x int) bool { return x > 10 },
			want:   []int{15, 20},
		},
		{
			name: "nil input",
			setupInput: func() any {
				var ch chan int = nil
				return ch
			},
			filter: func(x int) bool { return x > 0 },
			want:   nil,
		},
		{
			name: "empty input",
			setupInput: func() any {
				ch := make(chan int)
				close(ch)
				return ch
			},
			filter: func(x int) bool { return x > 0 },
			want:   []int{},
		},
		{
			name: "all items pass filter",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					for _, v := range []int{2, 4, 6, 8} {
						ch <- v
					}
				}()
				return ch
			},
			filter: func(x int) bool { return x%2 == 0 },
			want:   []int{2, 4, 6, 8},
		},
		{
			name: "no items pass filter",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					for _, v := range []int{1, 3, 5, 7} {
						ch <- v
					}
				}()
				return ch
			},
			filter: func(x int) bool { return x%2 == 0 },
			want:   []int{},
		},
		{
			name: "single item passes",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					ch <- 42
				}()
				return ch
			},
			filter: func(x int) bool { return x > 0 },
			want:   []int{42},
		},
		{
			name: "single item does not pass",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					ch <- 5
				}()
				return ch
			},
			filter: func(x int) bool { return x > 10 },
			want:   []int{},
		},
		{
			name: "string type filter",
			setupInput: func() any {
				ch := make(chan string)
				go func() {
					defer close(ch)
					for _, w := range []string{"hello", "world", "hi", "test"} {
						ch <- w
					}
				}()
				return ch
			},
			filter: func(s string) bool { return len(s) > 3 },
			want:   []string{"hello", "world", "test"},
		},
		{
			name: "struct type filter",
			setupInput: func() any {
				ch := make(chan Item)
				go func() {
					defer close(ch)
					for _, item := range []Item{
						{ID: 1, Value: "test"},
						{ID: 5, Value: "hello"},
						{ID: 3, Value: "world"},
						{ID: 10, Value: "foo"},
					} {
						ch <- item
					}
				}()
				return ch
			},
			filter: func(item Item) bool { return item.ID > 3 },
			want: []Item{
				{ID: 5, Value: "hello"},
				{ID: 10, Value: "foo"},
			},
		},
		{
			name: "preserves capacity",
			setupInput: func() any {
				ch := make(chan int, 10)
				go func() {
					defer close(ch)
					for i := 1; i <= 10; i++ {
						ch <- i
					}
				}()
				return ch
			},
			filter:  func(x int) bool { return x%2 == 0 },
			want:    []int{2, 4, 6, 8, 10},
			wantCap: 10,
		},
		{
			name: "large dataset",
			setupInput: func() any {
				ch := make(chan int)
				go func() {
					defer close(ch)
					for i := 0; i < 100; i++ {
						ch <- i
					}
				}()
				return ch
			},
			filter: func(x int) bool { return x%2 == 0 },
			want: func() []int {
				result := make([]int, 50)
				for i := 0; i < 50; i++ {
					result[i] = i * 2
				}
				return result
			}(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := tc.setupInput()

			switch in := input.(type) {
			case chan int:
				filter := tc.filter.(func(int) bool)
				output := ChanFilter(in, filter)

				if tc.want == nil {
					require.Nil(t, output)
					return
				}

				if tc.wantCap > 0 {
					require.Equal(t, tc.wantCap, cap(output))
				}

				var result []int
				for v := range output {
					result = append(result, v)
				}

				want := tc.want.([]int)
				if len(want) == 0 && len(result) == 0 {
					return
				}
				require.Equal(t, want, result)

			case chan string:
				filter := tc.filter.(func(string) bool)
				output := ChanFilter(in, filter)

				var result []string
				for v := range output {
					result = append(result, v)
				}
				require.Equal(t, tc.want, result)

			case chan Item:
				filter := tc.filter.(func(Item) bool)
				output := ChanFilter(in, filter)

				var result []Item
				for v := range output {
					result = append(result, v)
				}
				require.Equal(t, tc.want, result)
			}
		})
	}
}
