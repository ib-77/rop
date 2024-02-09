package rop

import (
	"context"
	"fmt"
	"sync"
)

func MassValidate[T any](ctx context.Context, inputs <-chan T,
	beginF func(in T) bool, cancelF func(in T) error, errMsg string) <-chan Rop[T] {
	out := make(chan Rop[T])
	go func() {
		for in := range inputs {

			select {
			case <-ctx.Done():
				MassValidateCancelWith(inputs, out, cancelF)
				return
			default:
				out <- Validate(in, beginF, errMsg)
			}
		}
	}()

	return out
}

func MassSwitch[In any, Out any](ctx context.Context, inputs <-chan Rop[In],
	switchF func(r In) Rop[Out], cancelF func(r Rop[In]) error) <-chan Rop[Out] {
	out := make(chan Rop[Out])

	go func() {
		for in := range inputs {

			select {
			case <-ctx.Done():
				MassCancelWith(inputs, out, cancelF)
				return
			default:
				out <- Switch(in, switchF)
			}

		}
	}()
	return out
}

func MassMap[In any, Out any](ctx context.Context, inputs <-chan Rop[In],
	mapF func(r In) Out, cancelF func(r Rop[In]) error) <-chan Rop[Out] {

	out := make(chan Rop[Out])

	go func() {
		for in := range inputs {

			select {
			case <-ctx.Done():
				MassCancelWith(inputs, out, cancelF)
				return
			default:
				out <- Map(in, mapF)
			}

		}
	}()
	return out
}

func MassTee[T any](ctx context.Context, inputs <-chan Rop[T],
	deadEndF func(r Rop[T]), cancelF func(r Rop[T]) error) <-chan Rop[T] {

	out := make(chan Rop[T])

	go func() {
		for in := range inputs {

			select {
			case <-ctx.Done():
				MassCancelWith(inputs, out, cancelF)
				return
			default:
				out <- Tee(in, deadEndF)
			}
		}
	}()
	return out
}

func MassDoubleMap[In any, Out any](ctx context.Context, inputs <-chan Rop[In],
	successF func(r In) Out, failF func(err error) Out, cancelF func(r Rop[In]) error) <-chan Rop[Out] {

	out := make(chan Rop[Out])

	go func() {
		for in := range inputs {

			select {
			case <-ctx.Done():
				MassCancelWith(inputs, out, cancelF)
				return
			default:
				out <- DoubleMap(in, successF, failF)
			}
		}
	}()
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

func MassTry[In any, Out any](ctx context.Context, inputs <-chan Rop[In],
	withErrF func(r In) (Out, error), cancelF func(r Rop[In]) error) <-chan Rop[Out] {

	out := make(chan Rop[Out])

	go func() {
		for in := range inputs {

			select {
			case <-ctx.Done():
				MassCancelWith(inputs, out, cancelF)
				return
			default:
				out <- Try(in, withErrF)
			}
		}
	}()
	return out
}

func MassCheck[In any](ctx context.Context, inputs <-chan Rop[In],
	boolF func(r In) bool, falseErrMsg string, cancelF func(r In) error) <-chan Rop[bool] {

	out := make(chan Rop[bool])

	go func() {
		for in := range inputs {

			select {
			case <-ctx.Done():
				MassCheckCancelWith(inputs, out, cancelF)
				return
			default:
				out <- Check(in, boolF, falseErrMsg)
			}
		}
	}()
	return out
}

func MassFinally[Out, In any](ctx context.Context, inputs <-chan Rop[In],
	successF func(r In) Out, failF func(err error) Out, cancelF func(r In) Out) <-chan Out {

	out := make(chan Out)

	go func() {
		for in := range inputs {

			select {
			case <-ctx.Done():
				MassFinallyCancelWith(inputs, out, cancelF)
				return
			default:
				out <- Finally(in, successF, failF)
			}
		}
	}()
	return out
}

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

func MassCancelWith[In any, Out any](inputs <-chan Rop[In], outs chan Rop[Out],
	cancelF func(r Rop[In]) error) <-chan Rop[Out] {

	for c := range inputs {
		outs <- CancelWith[In, Out](c, cancelF)
	}
	return outs
}

func MassFinallyCancelWith[Out, In any](inputs <-chan Rop[In], outs chan Out,
	cancelF func(r In) Out) <-chan Out {

	for c := range inputs {
		outs <- cancelF(c.Result())
	}

	return outs
}

func MassCheckCancelWith[In any](inputs <-chan Rop[In], outs chan Rop[bool],
	cancelF func(r In) error) <-chan Rop[bool] {

	for c := range inputs {
		outs <- Cancel[bool](cancelF(c.Result()))
	}
	return outs
}

func PrintAll[T any](ctx context.Context, inputs <-chan Rop[T]) {
	MassTee(ctx, inputs, func(r Rop[T]) {
		if r.IsSuccess() {
			fmt.Printf("result: %v", r.Result())
		} else {
			fmt.Printf("error: %v", r.Err())
		}
	}, func(r Rop[T]) error {
		return fmt.Errorf("canceled: %v", r)
	})
}
