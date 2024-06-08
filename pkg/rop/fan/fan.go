package fan

import (
	"container/ring"
	"context"
	"github.com/ib-77/rop/pkg/rop"
	"sync"
)

func InTee[T any](ctx context.Context, inputChs chan chan rop.Result[T]) <-chan rop.Result[T] {
	outputCh := make(chan rop.Result[T])

	var wg sync.WaitGroup
	wg.Add(len(inputChs))
	for inputCh := range inputChs {
		go func(ch <-chan rop.Result[T]) {
			defer wg.Done()
			for value := range ch {
				select {
				case <-ctx.Done():
					return
				case outputCh <- value:
				}
			}
		}(inputCh)
	}

	go func() {
		wg.Wait()
		close(outputCh)
	}()
	return outputCh
}

func InFinally[T any](ctx context.Context, inputChs chan chan T) <-chan T {
	outputCh := make(chan T)

	var wg sync.WaitGroup
	wg.Add(len(inputChs))
	for inputCh := range inputChs {
		go func(ch <-chan T) {
			defer wg.Done()
			for value := range ch {
				select {
				case <-ctx.Done():
					return
				case outputCh <- value:
				}
			}
		}(inputCh)
	}

	go func() {
		wg.Wait()
		close(outputCh)
	}()
	return outputCh
}

func OutRand[T any](ctx context.Context, inputCh chan rop.Result[T], chCount int) []chan rop.Result[T] {

	outs := makeOutputChs[T](chCount)

	var wg sync.WaitGroup
	wg.Add(chCount)

	for chIndex := 0; chIndex < chCount; chIndex++ {
		go func(ch chan<- rop.Result[T]) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					ch <- <-inputCh
				}
			}
		}(outs[chIndex])
	}

	go func() {
		wg.Wait()
		closeOutputChs[T](outs)
	}()

	return outs
}

func OutNext[T any](ctx context.Context, inputCh chan rop.Result[T], chCount int) []chan rop.Result[T] {
	outs := makeOutputChs[T](chCount)
	nextChIndex := make(chan int, 1)
	r := makeRing(chCount)
	nextChIndex <- r.Value.(int)

	var wg sync.WaitGroup
	wg.Add(chCount)

	for chIndex := 0; chIndex < chCount; chIndex++ {
		outCh := outs[chIndex]
		go func(ch chan<- rop.Result[T], index int) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					select {
					case next, ok := <-nextChIndex:
						if !ok {
							break
						}
						outs[next] <- <-inputCh
						r = r.Next()
						nextChIndex <- r.Value.(int)
						break
					default:
						break
					}
				}
			}
		}(outCh, chIndex)
	}

	go func() {
		wg.Wait()
		closeOutputChs[T](outs)
		close(nextChIndex)
		nextChIndex = nil
	}()

	return outs
}

func ChsToSlice[T any](inputChs chan chan T, count int) []chan T {
	outputChs := make([]chan T, count)
	for i := 0; i < count; i++ {
		outputChs[i] = <-inputChs
	}
	return outputChs
}

// bug
func SliceToChs[T any](inputChs []chan T) chan chan T {
	count := len(inputChs)
	outputChs := make(chan chan T, count)
	defer close(outputChs)
	for i := 0; i < count; i++ {
		outputChs <- inputChs[i]
	}
	return outputChs
}

func makeRing(count int) *ring.Ring {
	r := ring.New(count)
	for i := 0; i < count; i++ {
		r.Value = i
		r = r.Next()
	}
	return r
}

func makeOutputChs[Out any](outputChCount int) []chan rop.Result[Out] {
	outs := make([]chan rop.Result[Out], outputChCount)
	for i := 0; i < outputChCount; i++ {
		outs[i] = make(chan rop.Result[Out])
	}
	return outs
}

func closeOutputChs[Out any](outs []chan rop.Result[Out]) {
	for i := 0; i < len(outs); i++ {
		close(outs[i])
	}
}
