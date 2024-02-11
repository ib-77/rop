package rop

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
