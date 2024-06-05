package bridge

import (
	"context"
	"github.com/ib-77/rop/pkg/rop"
	"github.com/ib-77/rop/pkg/rop/mass"
	"github.com/ib-77/rop/pkg/rop/solo"
)

func Validate[T any](ctx context.Context, inputChs <-chan <-chan T,
	validateF func(ctx context.Context, in T) bool,
	cancelF func(ctx context.Context, in T) error, errMsg string) <-chan rop.Result[T] {

	outCh := make(chan rop.Result[T])
	go func() {
		defer close(outCh)

		for {
			var input <-chan T
			select {

			case maybeInput, ok := <-inputChs:
				if !ok {
					return
				}
				input = maybeInput

			case <-ctx.Done():
				return
			}

			for out := range mass.Validate(ctx, input, validateF, cancelF, errMsg) {
				select {
				case outCh <- out:
					//case <-ctx.Done():    << don't skip cancelled
				}
			}
		}
	}()
	return outCh
}

func AndValidate[T any](ctx context.Context, inputChs <-chan <-chan rop.Result[T],
	validateF func(ctx context.Context, in T) bool,
	cancelF func(ctx context.Context, in rop.Result[T]) error, errMsg string) <-chan rop.Result[T] {

	outCh := make(chan rop.Result[T])
	go func() {
		defer close(outCh)

		for {
			var input <-chan rop.Result[T]
			select {

			case maybeInput, ok := <-inputChs:
				if !ok {
					return
				}
				input = maybeInput

			case <-ctx.Done():
				return
			}

			for out := range mass.AndValidate(ctx, input, validateF, cancelF, errMsg) {
				select {
				case outCh <- out:
					//case <-ctx.Done():    << don't skip cancelled
				}
			}
		}
	}()
	return outCh
}

func Map[In, Out any](ctx context.Context, inputChs <-chan <-chan rop.Result[In],
	mapF func(ctx context.Context, r In) Out,
	cancelF func(ctx context.Context, r rop.Result[In]) error) <-chan rop.Result[Out] {

	outCh := make(chan rop.Result[Out])
	go func() {
		defer close(outCh)

		for {
			var input <-chan rop.Result[In]
			select {

			case maybeInput, ok := <-inputChs:
				if !ok {
					return
				}
				input = maybeInput

			case <-ctx.Done():
				return
			}

			for out := range mass.Map(ctx, input, mapF, cancelF) {
				select {
				case outCh <- out:
					//case <-ctx.Done():    << don't skip cancelled
				}
			}
		}
	}()
	return outCh
}

func Tee[T any](ctx context.Context, inputChs chan chan rop.Result[T],
	deadEndF func(ctx context.Context, r rop.Result[T]),
	cancelF func(ctx context.Context, r rop.Result[T]) error) <-chan rop.Result[T] {

	outCh := make(chan rop.Result[T])
	go func() {
		defer close(outCh)

		for {
			var input <-chan rop.Result[T]
			select {

			case maybeInput, ok := <-inputChs:
				if !ok {
					return
				}
				input = maybeInput

			case <-ctx.Done():
				return
			}

			for out := range mass.Tee(ctx, input, deadEndF, cancelF) {
				select {
				case outCh <- out:
					//case <-ctx.Done():    << don't skip cancelled
				}
			}
		}
	}()
	return outCh
}

func Switch[In, Out any](ctx context.Context, inputChs <-chan <-chan rop.Result[In],
	switchF func(ctx context.Context, r In) rop.Result[Out],
	cancelF func(ctx context.Context, r rop.Result[In]) error) <-chan rop.Result[Out] {

	outCh := make(chan rop.Result[Out])
	go func() {
		defer close(outCh)

		for {
			var input <-chan rop.Result[In]
			select {

			case maybeInput, ok := <-inputChs:
				if !ok {
					return
				}
				input = maybeInput

			case <-ctx.Done():
				return
			}

			for out := range mass.Switch(ctx, input, switchF, cancelF) {
				select {
				case outCh <- out:
					//case <-ctx.Done():    << don't skip cancelled
				}
			}
		}
	}()
	return outCh
}

func DoubleMap[In any, Out any](ctx context.Context, inputChs <-chan <-chan rop.Result[In],
	successF func(ctx context.Context, r In) Out,
	failF func(ctx context.Context, err error) Out,
	cancelF func(ctx context.Context, err error) Out,
	massCancelF func(ctx context.Context, r rop.Result[In]) error) <-chan rop.Result[Out] {

	outCh := make(chan rop.Result[Out])
	go func() {
		defer close(outCh)

		for {
			var input <-chan rop.Result[In]
			select {

			case maybeInput, ok := <-inputChs:
				if !ok {
					return
				}
				input = maybeInput

			case <-ctx.Done():
				return
			}

			for out := range mass.DoubleMap(ctx, input, successF, failF, cancelF, massCancelF) {
				select {
				case outCh <- out:
					//case <-ctx.Done():    << don't skip cancelled
				}
			}
		}
	}()
	return outCh
}

func Try[In any, Out any](ctx context.Context, inputChs <-chan <-chan rop.Result[In],
	withErrF func(ctx context.Context, r In) (Out, error),
	cancelF func(ctx context.Context, r rop.Result[In]) error) <-chan rop.Result[Out] {

	outCh := make(chan rop.Result[Out])
	go func() {
		defer close(outCh)

		for {
			var input <-chan rop.Result[In]
			select {

			case maybeInput, ok := <-inputChs:
				if !ok {
					return
				}
				input = maybeInput

			case <-ctx.Done():
				return
			}

			for out := range mass.Try(ctx, input, withErrF, cancelF) {
				select {
				case outCh <- out:
					//case <-ctx.Done():    << don't skip cancelled
				}
			}
		}
	}()
	return outCh
}

func Finally[Out, In any](ctx context.Context, inputChs <-chan <-chan rop.Result[In],
	successF func(ctx context.Context, r In) Out,
	failF func(ctx context.Context, err error) Out,
	cancelF func(ctx context.Context, r rop.Result[In]) Out) <-chan Out {

	outCh := make(chan Out)
	go func() {
		defer close(outCh)

		for {
			var input <-chan rop.Result[In]
			select {

			case maybeInput, ok := <-inputChs:
				if !ok {
					return
				}
				input = maybeInput

			case <-ctx.Done():
				return
			}

			for out := range mass.Finally(ctx, input, successF, failF, cancelF) {
				select {
				case outCh <- out:
					//case <-ctx.Done():    << don't skip cancelled
				}
			}
		}
	}()
	return outCh
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
