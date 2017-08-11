package work_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/justsocialapps/assert"
	"github.com/justsocialapps/justlib/work"
)

func Example() {

	// create a worker with 3 goroutines that can process jobs in parallel
	worker := work.NewWorker(3, func(p work.Payload) interface{} {
		// this is our worker function that is called for every new job
		fmt.Println(p.Data)
		time.Sleep(200 * time.Millisecond)
		return nil
	}, false)

	// create 100 jobs and dispatch them to the worker.
	for i := 0; i < 100; i++ {
		// this call will block when all 3 goroutines are currently busy.
		err := worker.Dispatch(work.Payload{Data: strconv.Itoa(i)})
		if err != nil {
			return // the worker has already been shut down
		}
	}

	// this call makes sure that the worker stops all goroutines as soon as
	// they have processed all remaining jobs.
	worker.Quit()

}

func TestDispatchShouldReturnErrorWhenWorkerIsShutDown(t *testing.T) {
	worker := work.NewWorker(1, func(p work.Payload) interface{} {
		return fmt.Sprintf("%s.", p.Data)
	}, true)
	worker.Quit()
	err := worker.Dispatch(work.Payload{})
	assert.NewAssert(t).NotNil(err, "Dispatch returned no error")
}

func TestCallingQuitMultipleTimesShouldNotBlock(t *testing.T) {
	worker := work.NewWorker(1, func(p work.Payload) interface{} {
		return fmt.Sprintf("%s.", p.Data)
	}, true)
	worker.Quit()
	worker.Quit()
}

func TestQuitShouldCloseCompletionsChannel(t *testing.T) {
	worker := work.NewWorker(1, func(p work.Payload) interface{} {
		return fmt.Sprintf("%s.", p.Data)
	}, true)
	worker.Quit()
	_, ok := <-worker.Completions()
	assert.NewAssert(t).Equal(ok, false, "The completions channel was not closed")
}

func TestQuitShouldAscertainThatAllJobsHaveCompleted(t *testing.T) {
	resCh := make(chan string, 100)
	worker := work.NewWorker(1, func(p work.Payload) interface{} {
		resCh <- fmt.Sprintf("%s.", p.Data)
		return nil
	}, false)

	for i := 0; i < 100; i++ {
		err := worker.Dispatch(work.Payload{Data: strconv.Itoa(i)})
		if err != nil {
			t.Fail()
			return
		}
	}

	worker.Quit()

	assert.NewAssert(t).Equal(len(resCh), 100, "Not all jobs were completed")
}

func TestWorkerShouldWorkSequentiallyWithOnlyOneGoroutine(t *testing.T) {

	resultCh := make(chan string)
	worker := work.NewWorker(1, func(p work.Payload) interface{} {
		return fmt.Sprintf("%s.", p.Data)
	}, true)

	go func() {
		var result string
		for v := range worker.Completions() {
			result += v.Output.(string)
		}
		resultCh <- result
	}()

	for i := 0; i < 100; i++ {
		err := worker.Dispatch(work.Payload{Data: strconv.Itoa(i)})
		if err != nil {
			t.Fail()
			return
		}
	}

	worker.Quit()

	assert.NewAssert(t).Equal(<-resultCh, "0.1.2.3.4.5.6.7.8.9.10.11.12.13.14.15.16.17.18.19.20.21.22.23.24.25.26.27.28.29.30.31.32.33.34.35.36.37.38.39.40.41.42.43.44.45.46.47.48.49.50.51.52.53.54.55.56.57.58.59.60.61.62.63.64.65.66.67.68.69.70.71.72.73.74.75.76.77.78.79.80.81.82.83.84.85.86.87.88.89.90.91.92.93.94.95.96.97.98.99.", "Jobs were completed in wrong order or incompletely")

}
