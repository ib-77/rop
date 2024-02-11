package rop

import (
	"context"
)

func MassValidate[T any](ctx context.Context, inputs <-chan T,
	validateF func(in T) bool, cancelF func(in T) error, errMsg string) <-chan Rop[T] {

	out := make(chan Rop[T])
	go func(ctx context.Context, inputs <-chan T) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- Cancel[T](cancelF(in)) // cancel current !!!
				MassValidateCancelWith(inputs, out, cancelF)
				break
			default:
				out <- Validate(in, validateF, errMsg)
			}
		}
	}(ctx, inputs)

	return out
}

func MassAndValidate[T any](ctx context.Context, inputs <-chan Rop[T],
	validateF func(in T) bool, cancelF func(in Rop[T]) error, errMsg string) <-chan Rop[T] {

	out := make(chan Rop[T])
	go func(ctx context.Context, inputs <-chan Rop[T], errMsg string) {
		defer close(out)

		for in := range inputs {
			select {
			case <-ctx.Done():
				out <- Cancel[T](cancelF(in)) // cancel current !!!
				MassAndValidateCancelWith[T](inputs, out, cancelF)
				break
			default:
				out <- AndValidate(in, validateF, errMsg)
			}
		}

	}(ctx, inputs, errMsg)

	return out
}

func MassSwitch[In any, Out any](ctx context.Context, inputs <-chan Rop[In],
	switchF func(r In) Rop[Out], cancelF func(r Rop[In]) error) <-chan Rop[Out] {
	out := make(chan Rop[Out])

	go func(ctx context.Context, inputs <-chan Rop[In]) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- CancelWith[In, Out](in, cancelF) // cancel current !!!
				MassCancelWith(inputs, out, cancelF)
				return
			default:
				out <- Switch(in, switchF)
			}

		}
	}(ctx, inputs)
	return out
}

func MassMap[In any, Out any](ctx context.Context, inputs <-chan Rop[In],
	mapF func(r In) Out, cancelF func(r Rop[In]) error) <-chan Rop[Out] {

	out := make(chan Rop[Out])

	go func(ctx context.Context, inputs <-chan Rop[In]) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- CancelWith[In, Out](in, cancelF) // cancel current !!!
				MassCancelWith(inputs, out, cancelF)
				return
			default:
				out <- Map(in, mapF)
			}

		}
	}(ctx, inputs)
	return out
}

func MassTee[T any](ctx context.Context, inputs <-chan Rop[T],
	deadEndF func(r Rop[T]), cancelF func(r Rop[T]) error) <-chan Rop[T] {

	out := make(chan Rop[T])

	go func(ctx context.Context, inputs <-chan Rop[T]) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- CancelWith[T, T](in, cancelF) // cancel current !!!
				MassCancelWith(inputs, out, cancelF)
				return
			default:
				out <- Tee(in, deadEndF)
			}
		}
	}(ctx, inputs)
	return out
}

func MassDoubleMap[In any, Out any](ctx context.Context, inputs <-chan Rop[In],
	successF func(r In) Out, failF func(err error) Out, cancelF func(r Rop[In]) error) <-chan Rop[Out] {

	out := make(chan Rop[Out])

	go func(ctx context.Context, inputs <-chan Rop[In]) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- CancelWith[In, Out](in, cancelF) // cancel current !!!
				MassCancelWith(inputs, out, cancelF)
				return
			default:
				out <- DoubleMap(in, successF, failF)
			}
		}
	}(ctx, inputs)
	return out
}

func MassSucceedWith[In any, Out any](inputs <-chan Rop[In], outs chan Rop[Out],
	successF func(r In) Out) <-chan Rop[Out] {

	for f := range inputs {
		outs <- SucceedWith(f, successF)
	}
	return outs
}

func MassFailWith[In any, Out any](inputs <-chan Rop[In], outs chan Rop[Out],
	failF func(r Rop[In]) error) <-chan Rop[Out] {

	for f := range inputs {
		outs <- FailWith[In, Out](f, failF)
	}
	return outs
}

func MassValidateCancelWith[In any](inputs <-chan In, outs chan Rop[In],
	cancelF func(r In) error) <-chan Rop[In] {

	for c := range inputs {
		outs <- Cancel[In](cancelF(c))
	}
	return outs
}

func MassAndValidateCancelWith[In any](inputs <-chan Rop[In], outs chan Rop[In],
	cancelF func(r Rop[In]) error) <-chan Rop[In] {

	for c := range inputs {
		outs <- Cancel[In](cancelF(c))
	}
	return outs
}

func MassTry[In any, Out any](ctx context.Context, inputs <-chan Rop[In],
	withErrF func(r In) (Out, error), cancelF func(r Rop[In]) error) <-chan Rop[Out] {

	out := make(chan Rop[Out])

	go func(ctx context.Context, inputs <-chan Rop[In]) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- CancelWith[In, Out](in, cancelF) // cancel current !!!
				MassCancelWith(inputs, out, cancelF)
				return
			default:
				out <- Try(in, withErrF)
			}
		}
	}(ctx, inputs)
	return out
}

func MassCheck[In any](ctx context.Context, inputs <-chan Rop[In],
	boolF func(r In) bool, falseErrMsg string, cancelF func(r Rop[In]) error) <-chan Rop[bool] {

	out := make(chan Rop[bool])

	go func(ctx context.Context, inputs <-chan Rop[In]) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- Cancel[bool](cancelF(in)) // cancel current !!!
				MassCheckCancelWith(inputs, out, cancelF)
				return
			default:
				out <- Check(in, boolF, falseErrMsg)
			}
		}
	}(ctx, inputs)
	return out
}

func MassFinally[Out, In any](ctx context.Context, inputs <-chan Rop[In],
	successF func(r In) Out, failF func(err error) Out, cancelF func(r Rop[In]) Out) <-chan Out {

	out := make(chan Out)

	go func(ctx context.Context, inputs <-chan Rop[In]) {
		defer close(out)

		for in := range inputs {

			select {
			case <-ctx.Done():
				out <- cancelF(in) // cancel current !!!
				MassFinallyCancelWith(inputs, out, cancelF)
				return
			default:
				out <- Finally(in, successF, failF)
			}
		}
	}(ctx, inputs)
	return out
}

func MassCancelWith[In any, Out any](inputs <-chan Rop[In], outs chan Rop[Out],
	cancelF func(r Rop[In]) error) <-chan Rop[Out] {

	for c := range inputs {
		outs <- CancelWith[In, Out](c, cancelF)
	}
	return outs
}

func MassFinallyCancelWith[Out, In any](inputs <-chan Rop[In], outs chan Out,
	cancelF func(r Rop[In]) Out) <-chan Out {

	for c := range inputs {
		outs <- cancelF(c)
	}

	return outs
}

func MassCheckCancelWith[In any](inputs <-chan Rop[In], outs chan Rop[bool],
	cancelF func(r Rop[In]) error) <-chan Rop[bool] {

	for c := range inputs {
		outs <- Cancel[bool](cancelF(c))
	}
	return outs
}
