// Package work contains routines to asynchronously dispatch jobs to
// a fixed set of Goroutines. Use this package to e.g. create a worker with 10
// Goroutines and distribute work among them.
package work

import "errors"

// Payload wraps the data that represents a job to be processed.
type Payload struct {
	Data interface{}
}

// Result wraps the result of the worker function for each input.
type Result struct {
	Input  interface{}
	Output interface{}
}

// Worker represents a worker unit that distributes jobs between its working
// Goroutines.
type Worker struct {
	workerCount int
	workerFunc  func(Payload) interface{}
	dispatcher  chan Payload  // channel for dispatching jobs to the worker
	completions chan Result   // channel with one entry per completed job
	quit        chan struct{} // channel used by Quit to signal the workers to exit
	inProgress  int
	shutdown    bool
}

// NewWorker creates a new worker that maintains exactly workerCount
// Goroutines. Each Goroutine calls workerFunc for processing the given data.
//
// If strictCompletions is true, then the result of every job is sent to the
// completions channel blockingly; else it is sent there non-blockingly which
// might result in in the loss of the result. You must read from the channel
// returned by Completions when you set strictCompletions to true, otherwise
// you'll have a deadlock.
func NewWorker(workerCount int, workerFunc func(Payload) interface{}, strictCompletions bool) *Worker {
	dispatcher := make(chan Payload)
	completions := make(chan Result)
	quit := make(chan struct{})

	// spawn worker goroutines
	for i := 0; i < workerCount; i++ {
		go func() {
			for {
				select {
				case job := <-dispatcher:
					output := workerFunc(job)
					// write the result blockingly
					if strictCompletions {
						completions <- Result{
							Input:  job.Data,
							Output: output,
						}
					} else {
						// write the result non-blockingly,
						// probably losing it when the
						// channel is full.
						select {
						case completions <- Result{Input: job.Data, Output: output}:
						default:
						}
					}
				case <-quit:
					// Quit was called, so we stop
					// processing jobs and exit.
					return
				}
			}
		}()
	}

	worker := &Worker{
		workerCount,
		workerFunc,
		dispatcher,
		completions,
		quit,
		0,
		false,
	}

	return worker

}

// Completions returns a channel that is sent a value to every time when a job
// is completed. This lets you keep track of completed jobs. Be aware that you
// need to start reading from this channel before dispatching jobs to the
// Worker, otherwise you would not receive all completions.
func (w *Worker) Completions() <-chan Result {
	return w.completions
}

// Dispatch feeds a new job to the Worker. If no idle Goroutine is available,
// this function blocks until the job can be processed. If the worker has
// already been shut down with Quit, an error is returned.
func (w *Worker) Dispatch(job Payload) error {
	if w.shutdown {
		return errors.New("WORKER IS SHUT DOWN")
	}
	w.dispatcher <- job
	return nil
}

// Quit waits for all workers to complete their jobs and returns afterwards.
// After this function returns, all worker routines are stopped and you cannot
// use this worker, anymore. Be aware that the job queue filled with Dispatch
// is not drained when you call Quit so that jobs might get lost. If you want
// to make sure that all jobs that are dispatched are also completed, read from
// the completions channel returned by Completions and only call Quit after you
// have received all results.
func (w *Worker) Quit() {
	if w.shutdown {
		return
	}
	for i := 0; i < w.workerCount; i++ {
		w.quit <- struct{}{}
	}
	w.shutdown = true
	close(w.completions) // signal to consumers that we're done so they don't block
}
