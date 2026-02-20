package ui

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var spinnerFrames = []string{"⣾ ", "⣽ ", "⣻ ", "⢿ ", "⡿ ", "⣟ ", "⣯ ", "⣷ "}

func WithSpinner[T any](message string, fn func() (T, error)) (T, error) {
	done := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 0
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r%s%s", spinnerFrames[i%len(spinnerFrames)], message)
				time.Sleep(100 * time.Millisecond)
				i++
			}
		}
	}()

	result, err := fn()

	close(done)
	wg.Wait()
	fmt.Printf("\r%s\r", strings.Repeat(" ", len(message)+len(spinnerFrames[0])))

	return result, err
}
