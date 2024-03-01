package mass

import (
	"github.com/ib-77/rop/pkg/rop"
	"sync"
)

func Aggregate[T any](inputs ...<-chan rop.Result[T]) <-chan rop.Result[T] {
	output := make(chan rop.Result[T])
	var wg sync.WaitGroup

	for _, in := range inputs {
		wg.Add(1)
		go func(int <-chan rop.Result[T]) {
			defer wg.Done()
			// If in is closed, then the
			// loop will ends eventually.
			for x := range in {
				output <- x
			}
		}(in)
	}
	go func() {
		wg.Wait()
		close(output)
	}()
	return output
}

func Divide[T any](input <-chan rop.Result[T], outputs ...chan<- rop.Result[T]) {
	for _, out := range outputs {
		go func(o chan<- rop.Result[T]) {
			for {
				o <- <-input // <=> o <- (<-input)
			}
		}(out)
	}
}
