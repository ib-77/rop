package rop

// interface pollution )
type WithError[T any] interface {
	Result() T
	Err() error
	IsSuccess()
}

type WithCancel[T any] interface {
	WithError[T]
	IsCancel() bool
}

//type Rop[T any] interface {
//	Result() T
//	Err() error
//	IsSuccess() bool
//	IsCancel() bool
//}
