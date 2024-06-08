package bridge

import (
	"context"
	"github.com/ib-77/rop/pkg/rop"
	"github.com/ib-77/rop/pkg/rop/mass"
	"github.com/ib-77/rop/pkg/rop/solo"
	"sync"
)

func Validate[T any](ctx context.Context, inputChs chan chan T,
	validateF func(ctx context.Context, in T) bool,
	cancelF func(ctx context.Context, in T) error, errMsg string) chan chan rop.Result[T] {

	outChs, outs := makeOutputChs[T](len(inputChs))

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		chIndex := 0
		for inputCh := range inputChs {

			outCh := outs[chIndex]

			select {
			case <-ctx.Done():
				return
			default:
			}

			wg.Add(1)
			go func(inCh <-chan T, ouCh chan rop.Result[T]) {
				defer wg.Done()

				for out := range mass.Validate(ctx, inCh, validateF, cancelF, errMsg) {
					select {
					case ouCh <- out:
						//case <-ctx.Done():    << don't skip cancelled
					}
				}
			}(inputCh, outCh)

			chIndex++
		}
	}()

	go func() {
		wg.Wait()
		closeOutputChs(outChs, outs)
	}()

	return outChs
}

func AndValidate[T any](ctx context.Context, inputChs chan chan rop.Result[T],
	validateF func(ctx context.Context, in T) bool,
	cancelF func(ctx context.Context, in rop.Result[T]) error, errMsg string) chan chan rop.Result[T] {

	outChs, outs := makeOutputChs[T](len(inputChs))

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		chIndex := 0
		for inputCh := range inputChs {

			outCh := outs[chIndex]

			select {
			case <-ctx.Done():
				return
			default:
			}

			wg.Add(1)
			go func(inCh <-chan rop.Result[T], ouCh chan rop.Result[T]) {
				defer wg.Done()

				for out := range mass.AndValidate(ctx, inCh, validateF, cancelF, errMsg) {
					select {
					case ouCh <- out:
						//case <-ctx.Done():    << don't skip cancelled
					}
				}
			}(inputCh, outCh)

			chIndex++
		}
	}()

	go func() {
		wg.Wait()
		closeOutputChs(outChs, outs)
	}()

	return outChs
}

func Map[In, Out any](ctx context.Context, inputChs chan chan rop.Result[In],
	mapF func(ctx context.Context, r In) Out,
	cancelF func(ctx context.Context, r rop.Result[In]) error) chan chan rop.Result[Out] {

	outChs, outs := makeOutputChs[Out](len(inputChs))

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		chIndex := 0
		for inputCh := range inputChs {

			outCh := outs[chIndex]

			select {
			case <-ctx.Done():
				return
			default:
			}

			wg.Add(1)
			go func(inCh <-chan rop.Result[In], ouCh chan rop.Result[Out]) {
				defer wg.Done()

				for out := range mass.Map(ctx, inCh, mapF, cancelF) {
					select {
					case ouCh <- out:
						//case <-ctx.Done():    << don't skip cancelled
					}
				}
			}(inputCh, outCh)

			chIndex++
		}
	}()

	go func() {
		wg.Wait()
		closeOutputChs(outChs, outs)
	}()

	return outChs

}

func Tee[T any](ctx context.Context, inputChs chan chan rop.Result[T],
	deadEndF func(ctx context.Context, r rop.Result[T]),
	cancelF func(ctx context.Context, r rop.Result[T]) error) chan chan rop.Result[T] {

	outChs, outs := makeOutputChs[T](len(inputChs))

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		chIndex := 0
		for inputCh := range inputChs {

			outCh := outs[chIndex]

			select {
			case <-ctx.Done():
				return
			default:
			}

			wg.Add(1)
			go func(inCh <-chan rop.Result[T], ouCh chan rop.Result[T]) {
				defer wg.Done()

				for out := range mass.Tee(ctx, inCh, deadEndF, cancelF) {
					select {
					case ouCh <- out:
						//case <-ctx.Done():    << don't skip cancelled
					}
				}
			}(inputCh, outCh)

			chIndex++
		}
	}()

	go func() {
		wg.Wait()
		closeOutputChs(outChs, outs)
	}()

	return outChs

}

func Switch[In, Out any](ctx context.Context, inputChs chan chan rop.Result[In],
	switchF func(ctx context.Context, r In) rop.Result[Out],
	cancelF func(ctx context.Context, r rop.Result[In]) error) chan chan rop.Result[Out] {

	outChs, outs := makeOutputChs[Out](len(inputChs))

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		chIndex := 0
		for inputCh := range inputChs {

			outCh := outs[chIndex]

			select {
			case <-ctx.Done():
				return
			default:
			}

			wg.Add(1)
			go func(inCh <-chan rop.Result[In], ouCh chan rop.Result[Out]) {
				defer wg.Done()

				for out := range mass.Switch(ctx, inCh, switchF, cancelF) {
					select {
					case ouCh <- out:
						//case <-ctx.Done():    << don't skip cancelled
					}
				}
			}(inputCh, outCh)

			chIndex++
		}
	}()

	go func() {
		wg.Wait()
		closeOutputChs(outChs, outs)
	}()

	return outChs

}

func DoubleMap[In any, Out any](ctx context.Context, inputChs chan chan rop.Result[In],
	successF func(ctx context.Context, r In) Out,
	failF func(ctx context.Context, err error) Out,
	cancelF func(ctx context.Context, err error) Out,
	massCancelF func(ctx context.Context, r rop.Result[In]) error) chan chan rop.Result[Out] {

	outChs, outs := makeOutputChs[Out](len(inputChs))

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		chIndex := 0
		for inputCh := range inputChs {

			outCh := outs[chIndex]

			select {
			case <-ctx.Done():
				return
			default:
			}

			wg.Add(1)
			go func(inCh <-chan rop.Result[In], ouCh chan rop.Result[Out]) {
				defer wg.Done()

				for out := range mass.DoubleMap(ctx, inCh, successF, failF, cancelF, massCancelF) {
					select {
					case ouCh <- out:
						//case <-ctx.Done():    << don't skip cancelled
					}
				}
			}(inputCh, outCh)

			chIndex++
		}
	}()

	go func() {
		wg.Wait()
		closeOutputChs(outChs, outs)
	}()

	return outChs

}

func Try[In, Out any](ctx context.Context, inputChs chan chan rop.Result[In],
	withErrF func(ctx context.Context, r In) (Out, error),
	cancelF func(ctx context.Context, r rop.Result[In]) error) chan chan rop.Result[Out] {

	outChs, outs := makeOutputChs[Out](len(inputChs))

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		chIndex := 0
		for inputCh := range inputChs {

			outCh := outs[chIndex]

			select {
			case <-ctx.Done():
				return
			default:
			}

			wg.Add(1)
			go func(inCh <-chan rop.Result[In], ouCh chan rop.Result[Out]) {
				defer wg.Done()

				for out := range mass.Try(ctx, inCh, withErrF, cancelF) {
					select {
					case ouCh <- out:
						//case <-ctx.Done():    << don't skip cancelled
					}
				}
			}(inputCh, outCh)

			chIndex++
		}
	}()

	go func() {
		wg.Wait()
		closeOutputChs(outChs, outs)
	}()

	return outChs

}

func Finally[In, Out any](ctx context.Context, inputChs chan chan rop.Result[In],
	successF func(ctx context.Context, r In) Out,
	failF func(ctx context.Context, err error) Out,
	cancelF func(ctx context.Context, r rop.Result[In]) Out) chan chan Out {

	outChs, outs := makeFinallyOutputChs[Out](len(inputChs))

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		chIndex := 0
		for inputCh := range inputChs {

			outCh := outs[chIndex]

			select {
			case <-ctx.Done():
				return
			default:
			}

			wg.Add(1)
			go func(inCh <-chan rop.Result[In], ouCh chan Out) {
				defer wg.Done()

				for out := range mass.Finally(ctx, inCh, successF, failF, cancelF) {
					select {
					case ouCh <- out:
						//case <-ctx.Done():    << don't skip cancelled
					}
				}
			}(inputCh, outCh)

			chIndex++
		}
	}()

	go func() {
		wg.Wait()
		closeFinallyOutputChs(outChs, outs)
	}()

	return outChs

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

func makeFinallyOutputChs[Out any](outputChCount int) (chan chan Out, []chan Out) {
	outputChs := make(chan chan Out, outputChCount)
	outs := make([]chan Out, outputChCount)
	for i := 0; i < outputChCount; i++ {
		c := make(chan Out)
		outputChs <- c
		outs[i] = c
	}
	return outputChs, outs
}

func closeFinallyOutputChs[Out any](outputChs chan chan Out, outs []chan Out) {
	for i := 0; i < len(outs); i++ {
		close(outs[i])
	}
	close(outputChs)
	outs = nil
}

func alternativeToRange[T any](ctx context.Context, input <-chan rop.Result[T],
	deadEndF func(ctx context.Context, r rop.Result[T]),
	cancelF func(ctx context.Context, r rop.Result[T]) error) <-chan rop.Result[T] {
	outCh := make(chan rop.Result[T])
	go func() {
		defer close(outCh)
		for {
			select {
			case <-ctx.Done():
				mass.CancelWithCtx(ctx, input, outCh, cancelF)
				return
			case in, ok := <-input:
				if !ok {
					return
				}
				select {
				case outCh <- solo.TeeWithCtx(ctx, in, deadEndF):
				case <-ctx.Done():
					outCh <- solo.CancelWithCtx[T, T](ctx, in, cancelF)
				}
			}
		}
	}()
	return outCh
}
