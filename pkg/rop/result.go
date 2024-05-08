package rop

type Result[T any] struct {
	result     T
	err        error
	isSuccess  bool
	isCancel   bool
	isAccepted bool
}

func Success[T any](r T) Result[T] {
	return Result[T]{
		result:    r,
		err:       nil,
		isSuccess: true,
		isCancel:  false,
	}
}

func Fail[T any](err error) Result[T] {
	return Result[T]{
		err:       err,
		isSuccess: false,
		isCancel:  false,
	}
}

func Cancel[T any](err error) Result[T] {
	return Result[T]{
		err:       err,
		isSuccess: false,
		isCancel:  true,
	}
}

func Accept[T any](r Result[T]) Result[T] {
	return Result[T]{
		err:        r.Err(),
		isSuccess:  r.IsSuccess(),
		isCancel:   r.IsCancel(),
		isAccepted: true,
	}
}

func (r Result[T]) Result() T {
	return r.result
}

func (r Result[T]) Err() error {
	return r.err
}

func (r Result[T]) IsSuccess() bool {
	return r.isSuccess
}

func (r Result[T]) IsCancel() bool {
	return r.isCancel
}

func (r Result[T]) IsAccepted() bool {
	return r.isAccepted
}
