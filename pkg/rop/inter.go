package rop

// WithError defines an interface for types that can return a result or an error
type WithError[T any] interface {
	// Result returns the successful result value
	Result() T
	// Err returns the error if operation failed
	Err() error
	// IsSuccess returns true if the operation was successful
	IsSuccess() bool
}

// WithCancel extends WithError with cancellation support
type WithCancel[T any] interface {
	WithError[T]
	// IsCancel returns true if the operation was cancelled
	IsCancel() bool
}
