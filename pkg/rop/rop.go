package rop

import (
	"errors"
)

type Rop[T any] interface {
	Result() T
	Err() error
	IsSuccess() bool
	IsCancel() bool
}

type ropResult[T any] struct {
	result    T
	err       error
	isSuccess bool
	isCancel  bool
}

func Success[T any](r T) Rop[T] {
	return ropResult[T]{
		result:    r,
		err:       nil,
		isSuccess: true,
		isCancel:  false,
	}
}

func Fail[T any](err error) Rop[T] {
	return ropResult[T]{
		err:       err,
		isSuccess: false,
		isCancel:  false,
	}
}

func Cancel[T any](err error) Rop[T] {
	return ropResult[T]{
		err:       err,
		isSuccess: false,
		isCancel:  true,
	}
}

func (r ropResult[T]) Result() T {
	return r.result
}

func (r ropResult[T]) Err() error {
	return r.err
}

func (r ropResult[T]) IsSuccess() bool {
	return r.isSuccess
}

func (r ropResult[T]) IsCancel() bool {
	return r.isCancel
}

func Validate[T any](input T, validateF func(in T) bool, errMsg string) Rop[T] {
	if validateF(input) {
		return Success(input)
	} else {
		return Fail[T](errors.New(errMsg))
	}
}

func AndValidate[T any](input Rop[T], validateF func(in T) bool, errMsg string) Rop[T] {
	if input.IsSuccess() {

		if validateF(input.Result()) {
			return Success(input.Result())
		} else {
			return Fail[T](errors.New(errMsg))
		}

	}
	return input
}

func Switch[In any, Out any](input Rop[In], switchF func(r In) Rop[Out]) Rop[Out] {

	if input.IsSuccess() {
		return switchF(input.Result())
	} else {
		return Fail[Out](input.Err())
	}
}

func Map[In any, Out any](input Rop[In], mapF func(r In) Out) Rop[Out] {

	if input.IsSuccess() {
		return Success(mapF(input.Result()))
	} else {
		return Fail[Out](input.Err())
	}
}

func Tee[T any](input Rop[T], deadEndF func(r Rop[T])) Rop[T] {

	if input.IsSuccess() {
		deadEndF(input)
	}

	return input
}

func TeeWithError[T any](input Rop[T], deadEndF func(r Rop[T]) error) Rop[T] {

	if input.IsSuccess() {
		err := deadEndF(input)
		if err != nil {
			return Fail[T](err)
		}
	}

	return input
}

func DoubleMap[In any, Out any](input Rop[In], successF func(r In) Out,
	failF func(err error) Out) Rop[Out] {

	if input.IsSuccess() {
		return Success(successF(input.Result()))
	}

	failF(input.Err())
	return Fail[Out](input.Err())
}

func Try[In any, Out any](input Rop[In], withErrF func(r In) (Out, error)) Rop[Out] {
	if input.IsSuccess() {

		out, err := withErrF(input.Result())
		if err != nil {
			return Fail[Out](err)
		}

		return Success(out)
	}
	return Fail[Out](input.Err())
}

func Check[In any](input Rop[In], boolF func(r In) bool, falseErrMsg string) Rop[bool] {

	if input.IsSuccess() {

		if ok := boolF(input.Result()); ok {
			return Success[bool](true)
		} else {
			return Fail[bool](errors.New(falseErrMsg))
		}
	}

	return Fail[bool](input.Err())
}

func Finally[Out, In any](input Rop[In], successF func(r In) Out, failF func(err error) Out) Out {
	if input.IsSuccess() {
		return successF(input.Result())
	} else {
		return failF(input.Err())
	}
}

func SucceedWith[In any, Out any](input Rop[In], successF func(r In) Out) Rop[Out] {
	return Success(successF(input.Result()))
}

func FailWith[In any, Out any](input Rop[In], failF func(r Rop[In]) error) Rop[Out] {
	return Fail[Out](failF(input))
}

func CancelWith[In any, Out any](input Rop[In], cancelF func(r Rop[In]) error) Rop[Out] { // cancelF out
	return Cancel[Out](cancelF(input))
}
