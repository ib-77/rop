package solo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ib-77/rop/pkg/rop"
)

var errCancelled = errors.New("cancelled")

func Validate[T any](input T, validateF func(in T) bool, errMsg string) rop.Result[T] {

	if validateF(input) {
		return rop.Success(input)
	} else {
		return rop.Fail[T](errors.New(errMsg))
	}
}

func ValidateWithCtx[T any](ctx context.Context, input T,
	validateF func(ctx context.Context, in T) bool, errMsg string) rop.Result[T] {

	if validateF(ctx, input) {
		return rop.Success(input)
	} else {
		return rop.Fail[T](errors.New(errMsg))
	}
}

// with cancel
func ValidateWithCtxWithCancel[T any](ctx context.Context, input T,
	validateF func(ctx context.Context, in T) bool, errMsg string) <-chan rop.Result[T] {

	validationCh := make(chan bool)
	resCh := make(chan rop.Result[T])

	go func() {
		validationCh <- validateF(ctx, input)
		close(validationCh)
	}()

	go func() {
		defer close(resCh)

		select {
		case isValid := <-validationCh:
			if isValid {
				resCh <- rop.Success(input)
			} else {
				resCh <- rop.Fail[T](errors.New(errMsg))
			}
		case <-ctx.Done():
			resCh <- rop.Cancel[T](errCancelled)
		}
	}()

	return resCh
}

func ValidateWithErr[T any](input T, validateF func(in T) (bool, error)) rop.Result[T] {

	if ok, err := validateF(input); ok {
		return rop.Success(input)
	} else {
		return rop.Fail[T](err)
	}
}

func ValidateWithErrWithCtx[T any](ctx context.Context, input T,
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

func ValidateCancelWithCtx[T any](ctx context.Context, input T,
	validateF func(ctx context.Context, in T) bool, cancelMsg string) rop.Result[T] {

	if validateF(ctx, input) {
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

func AndValidateWithCtx[T any](ctx context.Context, input rop.Result[T],
	validateF func(ctx context.Context, in T) bool, errMsg string) rop.Result[T] {

	if input.IsSuccess() {

		if validateF(ctx, input.Result()) {
			return rop.Success(input.Result())
		} else {
			return rop.Fail[T](errors.New(errMsg))
		}
	}
	return input
}

// with cancel
func AndValidateWithCtxWithCancel[T any](ctx context.Context, input rop.Result[T],
	validateF func(ctx context.Context, in T) bool, errMsg string) <-chan rop.Result[T] {

	validationCh := make(chan bool)
	resCh := make(chan rop.Result[T])

	go func() {
		if input.IsSuccess() {
			validationCh <- validateF(ctx, input.Result())
		}
		close(validationCh)
	}()

	go func() {
		defer close(resCh)
		if input.IsSuccess() {
			select {
			case isValid := <-validationCh:
				if isValid {
					resCh <- rop.Success(input.Result())
				} else {
					resCh <- rop.Fail[T](errors.New(errMsg))
				}
			case <-ctx.Done():
				resCh <- rop.Cancel[T](errCancelled)
			}
		} else {
			resCh <- input
		}
	}()

	return resCh
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

func AndValidateWithErrWithCtx[T any](ctx context.Context, input rop.Result[T],
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

func AndValidateCancelWithCtx[T any](ctx context.Context, input rop.Result[T],
	validateF func(ctx context.Context, in T) bool, cancelMsg string) rop.Result[T] {

	if input.IsSuccess() {
		if ok := validateF(ctx, input.Result()); ok {
			return rop.Success(input.Result())
		} else {
			return rop.Cancel[T](fmt.Errorf(cancelMsg))
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

// with cancel
func SwitchWithCtxWithCancel[In any, Out any](ctx context.Context,
	input rop.Result[In], switchF func(ctx context.Context, r In) rop.Result[Out]) <-chan rop.Result[Out] {

	switchCh := make(chan rop.Result[Out])
	outCh := make(chan rop.Result[Out])

	go func() {
		if input.IsSuccess() {
			switchCh <- switchF(ctx, input.Result())
		}
		close(switchCh)
	}()

	go func() {
		defer close(outCh)

		if input.IsSuccess() {

			select {
			case res := <-switchCh:
				outCh <- res
			case <-ctx.Done():
				outCh <- rop.Cancel[Out](errCancelled)
			}

		} else {
			if input.IsCancel() {
				outCh <- rop.Cancel[Out](input.Err())
			} else {
				outCh <- rop.Fail[Out](input.Err())
			}
		}
	}()

	return outCh
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

// with cancel
func MapWithCtxWithCancel[In any, Out any](ctx context.Context,
	input rop.Result[In], mapF func(ctx context.Context, r In) Out) <-chan rop.Result[Out] {

	mapCh := make(chan rop.Result[Out])
	outCh := make(chan rop.Result[Out])

	go func() {
		if input.IsSuccess() {
			mapCh <- rop.Success(mapF(ctx, input.Result()))
		}
		close(mapCh)
	}()

	go func() {
		defer close(outCh)

		if input.IsSuccess() {
			select {
			case res := <-mapCh:
				outCh <- res
			case <-ctx.Done():
				outCh <- rop.Cancel[Out](errCancelled)
			}
		} else {
			if input.IsCancel() {
				outCh <- rop.Cancel[Out](input.Err())
			} else {
				outCh <- rop.Fail[Out](input.Err())
			}
		}
	}()

	return outCh
}

func MapWithErrWithCtx[In any, Out any](ctx context.Context,
	input rop.Result[In], mapF func(ctx context.Context, r In) (Out, error)) rop.Result[Out] {

	if input.IsSuccess() {
		r, err := mapF(ctx, input.Result())
		if err != nil {
			return rop.Fail[Out](err)
		}
		return rop.Success(r)
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

// TeeWithError TODO unit test
func TeeWithErr[T any](input rop.Result[T], deadEndF func(r rop.Result[T]) error) rop.Result[T] {

	if input.IsSuccess() {
		err := deadEndF(input)
		if err != nil {
			return rop.Fail[T](err)
		}
	}

	return input
}

func TeeWithErrWithCtx[T any](ctx context.Context, input rop.Result[T],
	deadEndF func(ctx context.Context, r rop.Result[T]) error) rop.Result[T] {

	if input.IsSuccess() {
		err := deadEndF(ctx, input)
		if err != nil {
			return rop.Fail[T](err)
		}
	}

	return input
}

func TeeWithCtx[T any](ctx context.Context, input rop.Result[T],
	deadEndF func(ctx context.Context, r rop.Result[T])) rop.Result[T] {

	if input.IsSuccess() {
		deadEndF(ctx, input)
	}

	return input
}

// with cancel
func TeeWithCtxWithCancel[T any](ctx context.Context, input rop.Result[T],
	deadEndF func(ctx context.Context, r rop.Result[T])) <-chan rop.Result[T] {

	teeCh := make(chan struct{})
	outCh := make(chan rop.Result[T])

	go func() {
		if input.IsSuccess() {
			deadEndF(ctx, input)
		}
		close(teeCh)
	}()

	go func() {
		defer close(outCh)

		if input.IsSuccess() {
			select {
			case <-teeCh:
			case <-ctx.Done():
			}
		}
		outCh <- input
	}()

	return outCh
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

func DoubleTeeWithCtx[T any](ctx context.Context, input rop.Result[T],
	deadEndF func(ctx context.Context, r T),
	deadEndWithErrF func(ctx context.Context, err error)) rop.Result[T] {

	if input.IsSuccess() {
		deadEndF(ctx, input.Result())
	} else {
		deadEndWithErrF(ctx, input.Err())
	}

	return input
}

// with cancel
func DoubleTeeWithCtxWithCancel[T any](ctx context.Context, input rop.Result[T],
	deadEndF func(ctx context.Context, r T),
	deadEndWithErrF func(ctx context.Context, err error)) <-chan rop.Result[T] {

	teeCh := make(chan struct{})
	outCh := make(chan rop.Result[T])

	go func() {
		if input.IsSuccess() {
			deadEndF(ctx, input.Result())
		} else {
			deadEndWithErrF(ctx, input.Err())
		}
		close(teeCh)
	}()

	go func() {
		defer close(outCh)

		if input.IsSuccess() {
			select {
			case <-teeCh:
			case <-ctx.Done():
			}
		}
		outCh <- input
	}()

	return outCh
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

func DoubleMapWithCtx[In any, Out any](ctx context.Context, input rop.Result[In],
	successF func(ctx context.Context, r In) Out, failF func(ctx context.Context, err error) Out,
	cancelF func(ctx context.Context, err error) Out) rop.Result[Out] {

	if input.IsSuccess() {
		return rop.Success(successF(ctx, input.Result()))
	}

	if input.IsCancel() {
		cancelF(ctx, input.Err())
	} else {
		failF(ctx, input.Err())
	}

	if input.IsCancel() {
		return rop.Cancel[Out](input.Err())
	} else {
		return rop.Fail[Out](input.Err())
	}
}

// with cancel
func DoubleMapWithCtxWithCancel[In any, Out any](ctx context.Context, input rop.Result[In],
	successF func(ctx context.Context, r In) Out, failF func(ctx context.Context, err error) Out,
	cancelF func(ctx context.Context, err error) Out) <-chan rop.Result[Out] {

	mapCh := make(chan rop.Result[Out])
	outCh := make(chan rop.Result[Out])

	go func() {

		if input.IsSuccess() {
			mapCh <- rop.Success(successF(ctx, input.Result()))
		} else if input.IsCancel() {
			cancelF(ctx, input.Err())
		} else {
			failF(ctx, input.Err())
		}

		close(mapCh)
	}()

	go func() {
		defer close(outCh)

		if input.IsSuccess() {
			select {
			case res := <-mapCh:
				outCh <- res
			case <-ctx.Done():
				outCh <- rop.Cancel[Out](errCancelled)
			}
		} else {
			select {
			case <-mapCh:
				if input.IsCancel() {
					outCh <- rop.Cancel[Out](input.Err())
				} else {
					outCh <- rop.Fail[Out](input.Err())
				}
			case <-ctx.Done():
				outCh <- rop.Cancel[Out](errCancelled)
			}
		}
	}()

	return outCh
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

func TryWithCtx[In any, Out any](ctx context.Context, input rop.Result[In],
	withErrF func(ctx context.Context, r In) (Out, error)) rop.Result[Out] {

	if input.IsSuccess() {

		out, err := withErrF(ctx, input.Result())
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

func TryWithCtxWithCancel[In any, Out any](ctx context.Context, input rop.Result[In],
	withErrF func(ctx context.Context, r In) (Out, error)) <-chan rop.Result[Out] {

	tryCh := make(chan rop.Result[Out])
	outCh := make(chan rop.Result[Out])

	go func() {
		defer close(tryCh)
		if input.IsSuccess() {
			out, err := withErrF(ctx, input.Result())
			if err != nil {
				outCh <- rop.Fail[Out](err)
				return
			}
			outCh <- rop.Success(out)
		}
	}()

	go func() {
		if input.IsSuccess() {
			select {
			case res := <-tryCh:
				outCh <- res
			case <-ctx.Done():
				outCh <- rop.Cancel[Out](errCancelled)
			}
		} else {
			select {
			case <-tryCh:
				if input.IsCancel() {
					outCh <- rop.Cancel[Out](input.Err())
				} else {
					outCh <- rop.Fail[Out](input.Err())
				}
			case <-ctx.Done():
				outCh <- rop.Cancel[Out](errCancelled)
			}
		}
	}()
	return outCh
}

func RetryWithCtx[In any, Out any](ctx context.Context, input rop.Result[In],
	withErrF func(ctx context.Context, r In) (Out, error)) rop.Result[Out] {

	if input.IsSuccess() {

		rs, ok := rop.GetRetryFromCtx(ctx)
		if !ok {
			return rop.Fail[Out](fmt.Errorf("RetryWithCtx: context  is not set, use rop.WithRetry"))
		}

		var attempt int64 = 0
		var err error
		var out Out
		for {
			out, err = withErrF(ctx, input.Result())
			if err != nil {
				attempt++
				if attempt >= rs.Attempts() {
					break
				}
				time.Sleep(rs.Wait(attempt))
			} else {
				break
			}
		}

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

func CheckWithCtx[In any](ctx context.Context, input rop.Result[In],
	boolF func(ctx context.Context, r In) bool, falseErrMsg string) rop.Result[bool] {

	if input.IsSuccess() {

		if ok := boolF(ctx, input.Result()); ok {
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

func CheckCancelWithCtx[In any](ctx context.Context, input rop.Result[In],
	boolF func(ctx context.Context, r In) bool, falseCancelMsg string) rop.Result[bool] {

	if input.IsSuccess() {

		if ok := boolF(ctx, input.Result()); ok {
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

func FinallyWithCtx[Out, In any](ctx context.Context, input rop.Result[In],
	successF func(ctx context.Context, r In) Out, failOrCancelF func(ctx context.Context, err error) Out) Out {
	if input.IsSuccess() {
		return successF(ctx, input.Result())
	} else {
		return failOrCancelF(ctx, input.Err())
	}
}

func FinallyTeeWithCtx[In any](ctx context.Context, input rop.Result[In],
	successF func(ctx context.Context, r In), failOrCancelF func(ctx context.Context, err error)) {
	if input.IsSuccess() {
		successF(ctx, input.Result())
	} else {
		failOrCancelF(ctx, input.Err())
	}
}

func FinallyTeeWithCtxWithErr[In any](ctx context.Context, input rop.Result[In],
	successF func(ctx context.Context, r In) error, failOrCancelF func(ctx context.Context, err error) error) error {
	if input.IsSuccess() {
		return successF(ctx, input.Result())
	} else {
		return failOrCancelF(ctx, input.Err())
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

func Succeed[In any](input In) rop.Result[In] {
	return rop.Success(input)
}

func SucceedWithCtx[In any, Out any](ctx context.Context, input rop.Result[In],
	successF func(ctx context.Context, r In) Out) rop.Result[Out] {
	return rop.Success(successF(ctx, input.Result()))
}

func FailWith[In any, Out any](input rop.Result[In], failF func(r rop.Result[In]) error) rop.Result[Out] {
	if input.IsSuccess() {
		return rop.Fail[Out](failF(input))
	}

	if input.IsCancel() {
		return rop.Cancel[Out](input.Err()) // strange, is already canceled!
	}

	return rop.Fail[Out](input.Err()) // strange, is already failed
}

func Fail[In any](err error) rop.Result[In] {
	return rop.Fail[In](err)
}

func FailWithCtx[In any, Out any](ctx context.Context, input rop.Result[In],
	failF func(ctx context.Context, r rop.Result[In]) error) rop.Result[Out] {

	if input.IsSuccess() {
		return rop.Fail[Out](failF(ctx, input))
	}

	if input.IsCancel() {
		return rop.Cancel[Out](input.Err()) // strange, is already canceled!
	}

	return rop.Fail[Out](input.Err()) // strange, is already failed
}

func CancelWith[In any, Out any](input rop.Result[In], cancelF func(r rop.Result[In]) error) rop.Result[Out] { // cancelF out
	if input.IsSuccess() {
		return rop.Cancel[Out](cancelF(input))
	}

	if input.IsCancel() {
		return rop.Cancel[Out](input.Err())
	}

	return rop.Fail[Out](input.Err())
}

func Cancel[In any](err error) rop.Result[In] {
	return rop.Cancel[In](err)
}

func CancelWithCtx[In any, Out any](ctx context.Context, input rop.Result[In],
	cancelF func(ctx context.Context, r In) error) rop.Result[Out] {

	if input.IsSuccess() {
		return rop.Cancel[Out](cancelF(ctx, input.Result()))
	}

	if input.IsCancel() {
		return rop.Cancel[Out](input.Err())
	}

	return rop.Fail[Out](input.Err())
}
