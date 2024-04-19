package solo

import (
	"context"
	"errors"
	"github.com/ib-77/rop/pkg/rop"
)

func Validate[T any](input T, validateF func(in T) bool, errMsg string) rop.Result[T] {
	if validateF(input) {
		return rop.Success(input)
	} else {
		return rop.Fail[T](errors.New(errMsg))
	}
}

func ValidateWithErr[T any](input T, validateF func(in T) (bool, error)) rop.Result[T] {
	if ok, err := validateF(input); ok {
		return rop.Success(input)
	} else {
		return rop.Fail[T](err)
	}
}

func ValidateWithCtxWithErr[T any](ctx context.Context, input T,
	validateF func(ctx context.Context, in T) (bool, error)) rop.Result[T] {

	if ok, err := validateF(ctx, input); ok {
		return rop.Success(input)
	} else {
		return rop.Fail[T](err)
	}
}

func ValidateCancel[T any](input T, validateF func(in T) bool, cancelMsg string) rop.Result[T] {
	if validateF(input) {
		return rop.Success(input)
	} else {
		return rop.Cancel[T](errors.New(cancelMsg))
	}
}

func ValidateCancelWithErr[T any](input T, validateF func(in T) (bool, error)) rop.Result[T] {
	if ok, cancelErr := validateF(input); ok {
		return rop.Success(input)
	} else {
		return rop.Cancel[T](cancelErr)
	}
}

func AndValidate[T any](input rop.Result[T], validateF func(in T) bool, errMsg string) rop.Result[T] {
	if input.IsSuccess() {

		if validateF(input.Result()) {
			return rop.Success(input.Result())
		} else {
			return rop.Fail[T](errors.New(errMsg))
		}
	}
	return input
}

func AndValidateWithErr[T any](input rop.Result[T], validateF func(in T) (bool, error)) rop.Result[T] {
	if input.IsSuccess() {

		if ok, err := validateF(input.Result()); ok {
			return rop.Success(input.Result())
		} else {
			return rop.Fail[T](err)
		}
	}
	return input
}

func AndValidateWithCtxWithErr[T any](ctx context.Context, input rop.Result[T],
	validateF func(ctx context.Context, in T) (bool, error)) rop.Result[T] {
	if input.IsSuccess() {

		if ok, err := validateF(ctx, input.Result()); ok {
			return rop.Success(input.Result())
		} else {
			return rop.Fail[T](err)
		}
	}
	return input
}

func AndValidateCancelWithErr[T any](input rop.Result[T], validateF func(in T) (bool, error)) rop.Result[T] {
	if input.IsSuccess() {
		if ok, cancelErr := validateF(input.Result()); ok {
			return rop.Success(input.Result())
		} else {
			return rop.Cancel[T](cancelErr)
		}
	}
	return input
}

func Switch[In any, Out any](input rop.Result[In], switchF func(r In) rop.Result[Out]) rop.Result[Out] {

	if input.IsSuccess() {
		return switchF(input.Result())
	} else {
		if input.IsCancel() {
			return rop.Cancel[Out](input.Err())
		} else {
			return rop.Fail[Out](input.Err())
		}
	}
}

func SwitchWithCtx[In any, Out any](ctx context.Context,
	input rop.Result[In], switchF func(ctx context.Context, r In) rop.Result[Out]) rop.Result[Out] {

	if input.IsSuccess() {
		return switchF(ctx, input.Result())
	} else {
		if input.IsCancel() {
			return rop.Cancel[Out](input.Err())
		} else {
			return rop.Fail[Out](input.Err())
		}
	}
}

func Map[In any, Out any](input rop.Result[In], mapF func(r In) Out) rop.Result[Out] {

	if input.IsSuccess() {
		return rop.Success(mapF(input.Result()))
	} else {
		if input.IsCancel() {
			return rop.Cancel[Out](input.Err())
		} else {
			return rop.Fail[Out](input.Err())
		}
	}
}

func MapWithCtx[In any, Out any](ctx context.Context,
	input rop.Result[In], mapF func(ctx context.Context, r In) Out) rop.Result[Out] {

	if input.IsSuccess() {
		return rop.Success(mapF(ctx, input.Result()))
	} else {
		if input.IsCancel() {
			return rop.Cancel[Out](input.Err())
		} else {
			return rop.Fail[Out](input.Err())
		}
	}
}

func Tee[T any](input rop.Result[T], deadEndF func(r rop.Result[T])) rop.Result[T] {

	if input.IsSuccess() {
		deadEndF(input)
	}

	return input
}

// DoubleTee TODO unit test
func DoubleTee[T any](input rop.Result[T], deadEndF func(r T),
	deadEndWithErrF func(err error)) rop.Result[T] {

	if input.IsSuccess() {
		deadEndF(input.Result())
	} else {
		deadEndWithErrF(input.Err())
	}

	return input
}

// TeeWithError TODO unit test
func TeeWithError[T any](input rop.Result[T], deadEndF func(r rop.Result[T]) error) rop.Result[T] {

	if input.IsSuccess() {
		err := deadEndF(input)
		if err != nil {
			return rop.Fail[T](err)
		}
	}

	return input
}

func DoubleMap[In any, Out any](input rop.Result[In], successF func(r In) Out,
	failF func(err error) Out, cancelF func(err error) Out) rop.Result[Out] {

	if input.IsSuccess() {
		return rop.Success(successF(input.Result()))
	}

	if input.IsCancel() {
		cancelF(input.Err())
	} else {
		failF(input.Err())
	}

	if input.IsCancel() {
		return rop.Cancel[Out](input.Err())
	} else {
		return rop.Fail[Out](input.Err())
	}
}

func Try[In any, Out any](input rop.Result[In], withErrF func(r In) (Out, error)) rop.Result[Out] {
	if input.IsSuccess() {

		out, err := withErrF(input.Result())
		if err != nil {
			return rop.Fail[Out](err)
		}

		return rop.Success(out)
	}

	if input.IsCancel() {
		return rop.Cancel[Out](input.Err())
	} else {
		return rop.Fail[Out](input.Err())
	}
}

// Check TODO unit test
func Check[In any](input rop.Result[In], boolF func(r In) bool, falseErrMsg string) rop.Result[bool] {

	if input.IsSuccess() {

		if ok := boolF(input.Result()); ok {
			return rop.Success[bool](true)
		} else {
			return rop.Fail[bool](errors.New(falseErrMsg))
		}
	}

	if input.IsCancel() {
		return rop.Cancel[bool](input.Err())
	} else {
		return rop.Fail[bool](input.Err())
	}
}

// CheckCancel TODO unit test
func CheckCancel[In any](input rop.Result[In], boolF func(r In) bool, falseCancelMsg string) rop.Result[bool] {

	if input.IsSuccess() {

		if ok := boolF(input.Result()); ok {
			return rop.Success[bool](true)
		} else {
			return rop.Cancel[bool](errors.New(falseCancelMsg))
		}
	}

	if input.IsCancel() {
		return rop.Cancel[bool](input.Err())
	} else {
		return rop.Fail[bool](input.Err())
	}
}

func Finally[Out, In any](input rop.Result[In], successF func(r In) Out,
	failOrCancelF func(err error) Out) Out {
	if input.IsSuccess() {
		return successF(input.Result())
	} else {
		return failOrCancelF(input.Err())
	}
}

func FinallyWithErr[Out, In any](input rop.Result[In], successF func(r In) (Out, error),
	failOrCancelF func(err error) (Out, error)) (Out, error) {
	if input.IsSuccess() {
		return successF(input.Result())
	} else {
		return failOrCancelF(input.Err())
	}
}

func FinallyWithCtxWithErr[Out, In any](ctx context.Context,
	input rop.Result[In], successF func(ctx context.Context, r In) (Out, error),
	failOrCancelF func(ctx context.Context, err error) (Out, error)) (Out, error) {
	if input.IsSuccess() {
		return successF(ctx, input.Result())
	} else {
		return failOrCancelF(ctx, input.Err())
	}
}

func SucceedWith[In any, Out any](input rop.Result[In], successF func(r In) Out) rop.Result[Out] {
	return rop.Success(successF(input.Result()))
}

func FailWith[In any, Out any](input rop.Result[In], failF func(r rop.Result[In]) error) rop.Result[Out] {
	return rop.Fail[Out](failF(input))
}

func CancelWith[In any, Out any](input rop.Result[In], cancelF func(r rop.Result[In]) error) rop.Result[Out] { // cancelF out
	return rop.Cancel[Out](cancelF(input))
}
