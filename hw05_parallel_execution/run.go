package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var workersCount int
	tasksLen := len(tasks)

	wg := sync.WaitGroup{}

	ch := make(chan Task, tasksLen)
	chErr := make(chan bool, tasksLen)

	defer close(chErr)

	if tasksLen < n {
		workersCount = tasksLen
	} else {
		workersCount = n
	}

	for _, task := range tasks {
		ch <- task
	}

	close(ch)

	for i := 0; i < workersCount; i++ {
		wg.Add(1)

		go func() {
			for task := range ch {
				if m > 0 && len(chErr) >= m {
					break
				}

				err := task()

				if m > 0 && err != nil {
					chErr <- true
				}
			}

			wg.Done()
		}()
	}

	wg.Wait()

	if m > 0 && len(chErr) >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
