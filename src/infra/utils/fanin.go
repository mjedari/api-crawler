package utils

import (
	"sync"
)

// receives slice of channels and iterate over them
// for each channel it is created a new go routine and do the process inside it
// then writes the output to the destination channel which is unique

func Funnel(sources ...<-chan []byte) <-chan []byte {
	dest := make(chan []byte) // make sure to close the channel you've made
	var wg sync.WaitGroup

	wg.Add(len(sources))

	for _, ch := range sources {
		go func(c <-chan []byte) {
			defer wg.Done()

			for i := range c {
				dest <- i
			}
		}(ch)
	}

	// Start a goroutine to close dest after all sources close
	go func() { // <---- Note
		wg.Wait()
		defer close(dest)
	}()

	return dest
}
