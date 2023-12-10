package main

import "sync"

type AsyncResult[T any] struct {
	Value T
	Err   error
}

func (ar AsyncResult[T]) Unwrap() (T, error) {
	return ar.Value, ar.Err
}

func Async[T any](fn func() (T, error), wg *sync.WaitGroup) <-chan AsyncResult[T] {
	done := make(chan AsyncResult[T])

	if wg != nil {
		wg.Add(1)
	}

	go func() {
		value, err := fn()
		done <- AsyncResult[T]{value, err}

		if wg != nil {
			wg.Done()
		}

		close(done)
	}()

	return done
}
