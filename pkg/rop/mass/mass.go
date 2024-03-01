package mass

import (
	"context"
	"github.com/ib-77/rop/pkg/rop"
	"github.com/ib-77/rop/pkg/rop/solo"
)

func Validate[T any](ctx context.Context, inputs <-chan T,
	validateF func(in T) bool, cancelF func(in T) error, errMsg string) <-chan rop.Result[T] {

	out := make(chan rop.Result[T])
	go func(ctx context.Context, inputs <-chan T) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- rop.Cancel[T](cancelF(in)) // cancel current !!!
				ValidateCancelWith(inputs, out, cancelF)
				break
			default:
				out <- solo.Validate(in, validateF, errMsg)
			}
		}
	}(ctx, inputs)

	return out
}

func AndValidate[T any](ctx context.Context, inputs <-chan rop.Result[T],
	validateF func(in T) bool, cancelF func(in rop.Result[T]) error, errMsg string) <-chan rop.Result[T] {

	out := make(chan rop.Result[T])
	go func(ctx context.Context, inputs <-chan rop.Result[T], errMsg string) {
		defer close(out)

		for in := range inputs {
			select {
			case <-ctx.Done():
				out <- rop.Cancel[T](cancelF(in)) // cancel current !!!
				AndValidateCancelWith[T](inputs, out, cancelF)
				break
			default:
				out <- solo.AndValidate(in, validateF, errMsg)
			}
		}

	}(ctx, inputs, errMsg)

	return out
}

func Switch[In any, Out any](ctx context.Context, inputs <-chan rop.Result[In],
	switchF func(r In) rop.Result[Out], cancelF func(r rop.Result[In]) error) <-chan rop.Result[Out] {
	out := make(chan rop.Result[Out])

	go func(ctx context.Context, inputs <-chan rop.Result[In]) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- solo.CancelWith[In, Out](in, cancelF) // cancel current !!!
				CancelWith(inputs, out, cancelF)
				return
			default:
				out <- solo.Switch(in, switchF)
			}

		}
	}(ctx, inputs)
	return out
}

func Map[In any, Out any](ctx context.Context, inputs <-chan rop.Result[In],
	mapF func(r In) Out, cancelF func(r rop.Result[In]) error) <-chan rop.Result[Out] {

	out := make(chan rop.Result[Out])

	go func(ctx context.Context, inputs <-chan rop.Result[In]) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- solo.CancelWith[In, Out](in, cancelF) // cancel current !!!
				CancelWith(inputs, out, cancelF)
				return
			default:
				out <- solo.Map(in, mapF)
			}

		}
	}(ctx, inputs)
	return out
}

func Tee[T any](ctx context.Context, inputs <-chan rop.Result[T],
	deadEndF func(r rop.Result[T]), cancelF func(r rop.Result[T]) error) <-chan rop.Result[T] {

	out := make(chan rop.Result[T])

	go func(ctx context.Context, inputs <-chan rop.Result[T]) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- solo.CancelWith[T, T](in, cancelF) // cancel current !!!
				CancelWith(inputs, out, cancelF)
				return
			default:
				out <- solo.Tee(in, deadEndF)
			}
		}
	}(ctx, inputs)
	return out
}

func DoubleMap[In any, Out any](ctx context.Context, inputs <-chan rop.Result[In],
	successF func(r In) Out, failF func(err error) Out, cancelF func(err error) Out,
	massCancelF func(r rop.Result[In]) error) <-chan rop.Result[Out] {

	out := make(chan rop.Result[Out])

	go func(ctx context.Context, inputs <-chan rop.Result[In]) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- solo.CancelWith[In, Out](in, massCancelF) // cancel current !!!
				CancelWith(inputs, out, massCancelF)
				return
			default:
				out <- solo.DoubleMap(in, successF, failF, cancelF)
			}
		}
	}(ctx, inputs)
	return out
}

func MassSucceedWith[In any, Out any](inputs <-chan rop.Result[In], outs chan rop.Result[Out],
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

func AndValidateCancelWith[In any](inputs <-chan rop.Result[In], outs chan rop.Result[In],
	cancelF func(r rop.Result[In]) error) <-chan rop.Result[In] {

	for c := range inputs {
		outs <- rop.Cancel[In](cancelF(c))
	}
	return outs
}

func Try[In any, Out any](ctx context.Context, inputs <-chan rop.Result[In],
	withErrF func(r In) (Out, error), cancelF func(r rop.Result[In]) error) <-chan rop.Result[Out] {

	out := make(chan rop.Result[Out])

	go func(ctx context.Context, inputs <-chan rop.Result[In]) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- solo.CancelWith[In, Out](in, cancelF) // cancel current !!!
				CancelWith(inputs, out, cancelF)
				return
			default:
				out <- solo.Try(in, withErrF)
			}
		}
	}(ctx, inputs)
	return out
}

func Check[In any](ctx context.Context, inputs <-chan rop.Result[In],
	boolF func(r In) bool, falseErrMsg string, cancelF func(r rop.Result[In]) error) <-chan rop.Result[bool] {

	out := make(chan rop.Result[bool])

	go func(ctx context.Context, inputs <-chan rop.Result[In]) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- rop.Cancel[bool](cancelF(in)) // cancel current !!!
				CheckCancelWith(inputs, out, cancelF)
				return
			default:
				out <- solo.Check(in, boolF, falseErrMsg)
			}
		}
	}(ctx, inputs)
	return out
}

func Finally[Out, In any](ctx context.Context, inputs <-chan rop.Result[In],
	successF func(r In) Out, failF func(err error) Out, cancelF func(r rop.Result[In]) Out) <-chan Out {

	out := make(chan Out)

	go func(ctx context.Context, inputs <-chan rop.Result[In]) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- cancelF(in) // cancel current !!!
				FinallyCancelWith(inputs, out, cancelF)
				return
			default:
				out <- solo.Finally(in, successF, failF)
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

func FinallyCancelWith[Out, In any](inputs <-chan rop.Result[In], outs chan Out,
	cancelF func(r rop.Result[In]) Out) <-chan Out {

	for c := range inputs {
		outs <- cancelF(c)
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
