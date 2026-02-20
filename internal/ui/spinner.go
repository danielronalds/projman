package ui

import (
	"fmt"
	"strings"
	"time"
)

var spinnerFrames = []string{"⣾ ", "⣽ ", "⣻ ", "⢿ ", "⡿ ", "⣟ ", "⣯ ", "⣷ "}

func WithSpinner[T any](message string, fn func() (T, error)) (T, error) {
	done := make(chan struct{})

	go func() {
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
	fmt.Printf("\r%s\r", strings.Repeat(" ", len(message)+len(spinnerFrames[0])))

	return result, err
}
