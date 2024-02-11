package rop

import "sync"

func MassAggregate[T any](inputs ...<-chan Rop[T]) <-chan Rop[T] {
	output := make(chan Rop[T])
	var wg sync.WaitGroup

	for _, in := range inputs {
		wg.Add(1)
		go func(int <-chan Rop[T]) {
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

func MassDivide[T any](input <-chan Rop[T], outputs ...chan<- Rop[T]) {
	for _, out := range outputs {
		go func(o chan<- Rop[T]) {
			for {
				o <- <-input // <=> o <- (<-input)
			}
		}(out)
	}
}

//func MassPrint[T any](ctx context.Context, inputs <-chan Rop[T]) {
//	MassTee(ctx, inputs, func(r Rop[T]) {
//		if r.IsSuccess() {
//			fmt.Printf("result: %v", r.Result())
//		} else {
//			fmt.Printf("error: %v", r.Err())
//		}
//	}, func(r Rop[T]) error {
//		return fmt.Errorf("canceled: %v", r)
//	})
//}
