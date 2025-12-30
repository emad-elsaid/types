package types

// OrderedParallelizeChan splits input channel to multiple worker channels, processes them in parallel,
// then merges results back whilst preserving the original order of items.
func OrderedParallelizeChan[In, Out any](input <-chan In, workers int, process func(<-chan In) <-chan Out) <-chan Out {
	if input == nil {
		return nil
	}

	if workers < 1 {
		workers = 1
	}

	capacity := cap(input)
	workerInputs := make([]chan In, workers)
	for i := range workerInputs {
		workerInputs[i] = make(chan In, capacity)
	}

	// Distribute input items round-robin to worker channels
	go func() {
		index := 0
		for item := range input {
			workerInputs[index%workers] <- item
			index++
		}
		for _, ch := range workerInputs {
			close(ch)
		}
	}()

	// Start workers and collect their output channels
	workerOutputs := make([]<-chan Out, workers)
	for i := range workerInputs {
		workerOutputs[i] = process(workerInputs[i])
	}

	// Merge outputs in round-robin order to preserve sequence
	output := make(chan Out, capacity)
	go func() {
		defer close(output)
		activeCount := workers
		index := 0

		for activeCount > 0 {
			workerIdx := index % workers
			if workerOutputs[workerIdx] != nil {
				if val, ok := <-workerOutputs[workerIdx]; ok {
					output <- val
				} else {
					workerOutputs[workerIdx] = nil
					activeCount--
				}
			}
			index++
		}
	}()

	return output
}
