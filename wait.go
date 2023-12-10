package main

import (
	"fmt"
	"sync"
	"time"
)

func wait(seconds int, wg *sync.WaitGroup) <-chan int {
	done := make(chan int)

	if wg != nil {
		wg.Add(1)
	}

	go func() {
		for i := 0; i < seconds; i++ {
			fmt.Printf(" - %d seconds...\n", seconds-i)
			time.Sleep(time.Second)
		}

		if wg != nil {
			wg.Done()
		}

		close(done)
	}()

	return done
}
