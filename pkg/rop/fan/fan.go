package fan

import (
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

func OutRand[T any](ctx context.Context, inputCh chan rop.Result[T], chCount int) chan chan rop.Result[T] {

	outputChs, outs := makeOutputChs[T](chCount)

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
					ch <- <-inputCh // <=> o <- (<-input)
				}
			}
		}(outs[chIndex])
	}

	go func() {
		wg.Wait()
		closeOutputChs[T](outputChs, outs)
	}()

	return outputChs
}

func makeOutputChs[Out any](outputChCount int) (chan chan rop.Result[Out], []chan rop.Result[Out]) {
	outputChs := make(chan chan rop.Result[Out], outputChCount)
	outs := make([]chan rop.Result[Out], outputChCount)
	for i := 0; i < outputChCount; i++ {
		c := make(chan rop.Result[Out])
		outputChs <- c
		outs[i] = c
	}
	return outputChs, outs
}

func closeOutputChs[Out any](outputChs chan chan rop.Result[Out], outs []chan rop.Result[Out]) {
	for i := 0; i < len(outs); i++ {
		close(outs[i])
	}
	close(outputChs)
	outs = nil
}
