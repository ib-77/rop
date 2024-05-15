package mass

import (
	"context"
	"github.com/ib-77/rop/pkg/rop"
	"github.com/ib-77/rop/pkg/rop/solo"
	"time"
)

func SliceToChan[T any](input []T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)

		for _, in := range input {
			out <- in
		}
	}()
	return out
}

func SliceToChanWithCancelCtx[T any](ctx context.Context, input []T,
	isCancelF func(i int, in T) bool) (context.Context, <-chan T) {

	out := make(chan T)
	newCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		defer close(out)

		for i, in := range input {
			if isCancelF(i, in) {
				cancel()
			}
			out <- in
		}
	}()

	return newCtx, out
}

func SliceToChanWithTimeoutCtx[T any](ctx context.Context, input []T,
	timeout time.Duration) (context.Context, <-chan T) {

	out := make(chan T)
	newCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	go func() {
		defer close(out)

		for _, in := range input {
			out <- in
		}
	}()

	return newCtx, out
}

func Validate[T any](ctx context.Context, inputs <-chan T,
	validateF func(ctx context.Context, in T) bool,
	cancelF func(ctx context.Context, in T) error, errMsg string) <-chan rop.Result[T] {

	out := make(chan rop.Result[T])
	go func(ctx context.Context, inputs <-chan T) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- rop.Cancel[T](cancelF(ctx, in)) // cancel current !!!
				ValidateCancelWithCtx(ctx, inputs, out, cancelF)
				return
			default:
				out <- solo.ValidateWithCtx(ctx, in, validateF, errMsg)
			}
		}
	}(ctx, inputs)

	return out
}

func AndValidate[T any](ctx context.Context, inputs <-chan rop.Result[T],
	validateF func(ctx context.Context, in T) bool,
	cancelF func(ctx context.Context, in rop.Result[T]) error, errMsg string) <-chan rop.Result[T] {

	out := make(chan rop.Result[T])
	go func(ctx context.Context, inputs <-chan rop.Result[T], errMsg string) {
		defer close(out)

		for in := range inputs {
			select {
			case <-ctx.Done():
				out <- rop.Cancel[T](cancelF(ctx, in)) // cancel current !!!
				AndValidateCancelWithCtx[T](ctx, inputs, out, cancelF)
				return
			default:
				out <- solo.AndValidateWithCtx(ctx, in, validateF, errMsg)
			}
		}

	}(ctx, inputs, errMsg)

	return out
}

func Switch[In any, Out any](ctx context.Context, inputs <-chan rop.Result[In],
	switchF func(ctx context.Context, r In) rop.Result[Out],
	cancelF func(ctx context.Context, r rop.Result[In]) error) <-chan rop.Result[Out] {
	out := make(chan rop.Result[Out])

	go func(ctx context.Context, inputs <-chan rop.Result[In]) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- solo.CancelWithCtx[In, Out](ctx, in, cancelF) // cancel current !!!
				CancelWithCtx(ctx, inputs, out, cancelF)
				return
			default:
				out <- solo.SwitchWithCtx(ctx, in, switchF)
			}

		}
	}(ctx, inputs)
	return out
}

func Map[In any, Out any](ctx context.Context, inputs <-chan rop.Result[In],
	mapF func(ctx context.Context, r In) Out,
	cancelF func(ctx context.Context, r rop.Result[In]) error) <-chan rop.Result[Out] {

	out := make(chan rop.Result[Out])

	go func(ctx context.Context, inputs <-chan rop.Result[In]) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- solo.CancelWithCtx[In, Out](ctx, in, cancelF) // cancel current !!!
				CancelWithCtx(ctx, inputs, out, cancelF)
				return
			default:
				out <- solo.MapWithCtx(ctx, in, mapF)
			}

		}
	}(ctx, inputs)
	return out
}

func Tee[T any](ctx context.Context, inputs <-chan rop.Result[T],
	deadEndF func(ctx context.Context, r rop.Result[T]),
	cancelF func(ctx context.Context, r rop.Result[T]) error) <-chan rop.Result[T] {

	out := make(chan rop.Result[T])

	go func(ctx context.Context, inputs <-chan rop.Result[T]) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- solo.CancelWithCtx[T, T](ctx, in, cancelF) // cancel current !!!
				CancelWithCtx(ctx, inputs, out, cancelF)
				return
			default:
				out <- solo.TeeWithCtx(ctx, in, deadEndF)
			}
		}
	}(ctx, inputs)
	return out
}

func DoubleMap[In any, Out any](ctx context.Context, inputs <-chan rop.Result[In],
	successF func(ctx context.Context, r In) Out,
	failF func(ctx context.Context, err error) Out,
	cancelF func(ctx context.Context, err error) Out,
	massCancelF func(ctx context.Context, r rop.Result[In]) error) <-chan rop.Result[Out] {

	out := make(chan rop.Result[Out])

	go func(ctx context.Context, inputs <-chan rop.Result[In]) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- solo.CancelWithCtx[In, Out](ctx, in, massCancelF) // cancel current !!!
				CancelWithCtx(ctx, inputs, out, massCancelF)
				return
			default:
				out <- solo.DoubleMapWithCtx(ctx, in, successF, failF, cancelF)
			}
		}
	}(ctx, inputs)
	return out
}

func SucceedWith[In any, Out any](inputs <-chan rop.Result[In], outs chan rop.Result[Out],
	successF func(r In) Out) <-chan rop.Result[Out] {

	for f := range inputs {
		outs <- solo.SucceedWith(f, successF)
	}
	return outs
}

func FailWith[In any, Out any](inputs <-chan rop.Result[In], outs chan rop.Result[Out],
	failF func(r rop.Result[In]) error) <-chan rop.Result[Out] {

	for f := range inputs {
		outs <- solo.FailWith[In, Out](f, failF)
	}
	return outs
}

func ValidateCancelWith[In any](inputs <-chan In, outs chan rop.Result[In],
	cancelF func(r In) error) <-chan rop.Result[In] {

	for c := range inputs {
		outs <- rop.Cancel[In](cancelF(c))
	}
	return outs
}

func ValidateCancelWithCtx[In any](ctx context.Context, inputs <-chan In,
	outs chan rop.Result[In],
	cancelF func(ctx context.Context, r In) error) <-chan rop.Result[In] {

	for c := range inputs {
		outs <- rop.Cancel[In](cancelF(ctx, c))
	}
	return outs
}

func AndValidateCancelWith[In any](inputs <-chan rop.Result[In], outs chan rop.Result[In],
	cancelF func(r rop.Result[In]) error) <-chan rop.Result[In] {

	for c := range inputs {
		outs <- rop.Cancel[In](cancelF(c))
	}
	return outs
}

func AndValidateCancelWithCtx[In any](ctx context.Context, inputs <-chan rop.Result[In], outs chan rop.Result[In],
	cancelF func(ctx context.Context, r rop.Result[In]) error) <-chan rop.Result[In] {

	for c := range inputs {
		outs <- rop.Cancel[In](cancelF(ctx, c))
	}
	return outs
}

func Try[In any, Out any](ctx context.Context, inputs <-chan rop.Result[In],
	withErrF func(ctx context.Context, r In) (Out, error),
	cancelF func(ctx context.Context, r rop.Result[In]) error) <-chan rop.Result[Out] {

	out := make(chan rop.Result[Out])

	go func(ctx context.Context, inputs <-chan rop.Result[In]) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- solo.CancelWithCtx[In, Out](ctx, in, cancelF) // cancel current !!!
				CancelWithCtx(ctx, inputs, out, cancelF)
				return
			default:
				out <- solo.TryWithCtx(ctx, in, withErrF)
			}
		}
	}(ctx, inputs)
	return out
}

func Check[In any](ctx context.Context, inputs <-chan rop.Result[In],
	boolF func(ctx context.Context, r In) bool, falseErrMsg string,
	cancelF func(ctx context.Context, r rop.Result[In]) error) <-chan rop.Result[bool] {

	out := make(chan rop.Result[bool])

	go func(ctx context.Context, inputs <-chan rop.Result[In], falseErrMsg string) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- rop.Cancel[bool](cancelF(ctx, in)) // cancel current !!!
				CheckCancelWithCtx(ctx, inputs, out, cancelF)
				return
			default:
				out <- solo.CheckWithCtx(ctx, in, boolF, falseErrMsg)
			}
		}
	}(ctx, inputs, falseErrMsg)
	return out
}

func Finally[Out, In any](ctx context.Context, inputs <-chan rop.Result[In],
	successF func(ctx context.Context, r In) Out,
	failF func(ctx context.Context, err error) Out,
	cancelF func(ctx context.Context, r rop.Result[In]) Out) <-chan Out {

	out := make(chan Out)

	go func(ctx context.Context, inputs <-chan rop.Result[In]) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- cancelF(ctx, in) // cancel current !!!
				FinallyCancelWithCtx(ctx, inputs, out, cancelF)
				return
			default:
				out <- solo.FinallyWithCtx(ctx, in, successF, failF)
			}
		}
	}(ctx, inputs)
	return out
}

func CancelWith[In any, Out any](inputs <-chan rop.Result[In], outs chan rop.Result[Out],
	cancelF func(r rop.Result[In]) error) <-chan rop.Result[Out] {

	for c := range inputs {
		outs <- solo.CancelWith[In, Out](c, cancelF)
	}
	return outs
}

func CancelWithCtx[In any, Out any](ctx context.Context, inputs <-chan rop.Result[In],
	outs chan rop.Result[Out],
	cancelF func(ctx context.Context, r rop.Result[In]) error) <-chan rop.Result[Out] {

	for c := range inputs {
		outs <- solo.CancelWithCtx[In, Out](ctx, c, cancelF)
	}
	return outs
}

func FinallyCancelWith[Out, In any](inputs <-chan rop.Result[In], outs chan Out,
	cancelF func(r rop.Result[In]) Out) <-chan Out {

	for c := range inputs {
		outs <- cancelF(c)
	}

	return outs
}

func FinallyCancelWithCtx[Out, In any](ctx context.Context, inputs <-chan rop.Result[In],
	outs chan Out, cancelF func(ctx context.Context, r rop.Result[In]) Out) <-chan Out {

	for c := range inputs {
		outs <- cancelF(ctx, c)
	}

	return outs
}

func CheckCancelWith[In any](inputs <-chan rop.Result[In], outs chan rop.Result[bool],
	cancelF func(r rop.Result[In]) error) <-chan rop.Result[bool] {

	for c := range inputs {
		outs <- rop.Cancel[bool](cancelF(c))
	}
	return outs
}

func CheckCancelWithCtx[In any](ctx context.Context, inputs <-chan rop.Result[In],
	outs chan rop.Result[bool],
	cancelF func(ctx context.Context, r rop.Result[In]) error) <-chan rop.Result[bool] {

	for c := range inputs {
		outs <- rop.Cancel[bool](cancelF(ctx, c))
	}
	return outs
}
