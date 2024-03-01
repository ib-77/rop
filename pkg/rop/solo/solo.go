package solo

import (
	"errors"
	"github.com/ib-77/rop/pkg/rop"
)

func Validate[T any](input T, validateF func(in T) bool, errMsg string) rop.Rop[T] {
	if validateF(input) {
		return rop.Success(input)
	} else {
		return rop.Fail[T](errors.New(errMsg))
	}
}

func ValidateCancel[T any](input T, validateF func(in T) bool, cancelMsg string) rop.Rop[T] {
	if validateF(input) {
		return rop.Success(input)
	} else {
		return rop.Cancel[T](errors.New(cancelMsg))
	}
}

func AndValidate[T any](input rop.Rop[T], validateF func(in T) bool, errMsg string) rop.Rop[T] {
	if input.IsSuccess() {

		if validateF(input.Result()) {
			return rop.Success(input.Result())
		} else {
			return rop.Fail[T](errors.New(errMsg))
		}
	}
	return input
}

func AndValidateCancel[T any](input rop.Rop[T], validateF func(in T) bool, cancelMsg string) rop.Rop[T] {
	if input.IsSuccess() {

		if validateF(input.Result()) {
			return rop.Success(input.Result())
		} else {
			return rop.Cancel[T](errors.New(cancelMsg))
		}
	}
	return input
}

func Switch[In any, Out any](input rop.Rop[In], switchF func(r In) rop.Rop[Out]) rop.Rop[Out] {

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

func Map[In any, Out any](input rop.Rop[In], mapF func(r In) Out) rop.Rop[Out] {

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

func Tee[T any](input rop.Rop[T], deadEndF func(r rop.Rop[T])) rop.Rop[T] {

	if input.IsSuccess() {
		deadEndF(input)
	}

	return input
}

// TeeWithError TODO unit test
func TeeWithError[T any](input rop.Rop[T], deadEndF func(r rop.Rop[T]) error) rop.Rop[T] {

	if input.IsSuccess() {
		err := deadEndF(input)
		if err != nil {
			return rop.Fail[T](err)
		}
	}

	return input
}

func DoubleMap[In any, Out any](input rop.Rop[In], successF func(r In) Out,
	failF func(err error) Out, cancelF func(err error) Out) rop.Rop[Out] {

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

func Try[In any, Out any](input rop.Rop[In], withErrF func(r In) (Out, error)) rop.Rop[Out] {
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
func Check[In any](input rop.Rop[In], boolF func(r In) bool, falseErrMsg string) rop.Rop[bool] {

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
func CheckCancel[In any](input rop.Rop[In], boolF func(r In) bool, falseCancelMsg string) rop.Rop[bool] {

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

func Finally[Out, In any](input rop.Rop[In], successF func(r In) Out,
	failOrCancelF func(err error) Out) Out {
	if input.IsSuccess() {
		return successF(input.Result())
	} else {
		return failOrCancelF(input.Err())
	}
}

func SucceedWith[In any, Out any](input rop.Rop[In], successF func(r In) Out) rop.Rop[Out] {
	return rop.Success(successF(input.Result()))
}

func FailWith[In any, Out any](input rop.Rop[In], failF func(r rop.Rop[In]) error) rop.Rop[Out] {
	return rop.Fail[Out](failF(input))
}

func CancelWith[In any, Out any](input rop.Rop[In], cancelF func(r rop.Rop[In]) error) rop.Rop[Out] { // cancelF out
	return rop.Cancel[Out](cancelF(input))
}
