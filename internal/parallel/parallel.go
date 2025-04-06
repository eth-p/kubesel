package parallel

import (
	"iter"
	"sync"
)

func Ordered[I any, O any](inputs []I, transformer func(input I) O) iter.Seq2[int, O] {
	type taskResult struct {
		output O
		done   bool
	}

	nTasks := len(inputs)
	tasks := make([]taskResult, nTasks)
	notiCh := make(chan int, nTasks)

	// Start a goroutine for each input.
	for i := range nTasks {
		go func() {
			tasks[i].output = transformer(inputs[i])
			notiCh <- i
		}()
	}

	// Wait until all tasks are finished, yielding whenever possible.
	return func(yield func(int, O) bool) {
		head := 0
		notis := 0

	yieldingLoop:
		for notis < nTasks {
			i := <-notiCh
			notis++
			tasks[i].done = true
			if notis < head {
				continue
			}

			// If the task completed is the next one in the series, yield until
			// we reach an unfinished task.
			for ; head < nTasks; head++ {
				if !tasks[head].done {
					continue yieldingLoop
				}

				if !yield(head, tasks[head].output) {
					goto drainingLoop
				}
			}
		}

	drainingLoop:
		for notis < nTasks {
			<-notiCh
		}

		close(notiCh)
	}
}

func Run[I any](inputs []I, mutator func(input I)) {
	var wg sync.WaitGroup
	nTasks := len(inputs)

	// Start a goroutine for each input.
	wg.Add(nTasks)
	for i := range nTasks {
		go func() {
			mutator(inputs[i])
			wg.Done()
		}()
	}

	// Wait until all tasks are finished.
	wg.Wait()
}

func Mutate[I any](inputs []I, mutator func(input *I)) {
	var wg sync.WaitGroup
	nTasks := len(inputs)

	// Start a goroutine for each input.
	wg.Add(nTasks)
	for i := range nTasks {
		go func() {
			mutator(&inputs[i])
			wg.Done()
		}()
	}

	// Wait until all tasks are finished.
	wg.Wait()
}
