package group

import (
	"context"
	"errors"
	"fmt"
	"github.com/ib-77/rop/pkg/rop"
)

const (
	EmptyResult = "empty result"
)

// AndWithCtx TODO ctx done
func AndWithCtx[In any](ctx context.Context,
	validateF func(ctx context.Context, resId int, in rop.Result[In]) rop.Result[In],
	accumF func(ctx context.Context, resId int, in rop.Result[In]) rop.Result[In],
	fs ...func(ctx context.Context) rop.Result[In]) rop.Result[In] {

	if len(fs) == 0 {
		return rop.Fail[In](fmt.Errorf(EmptyResult))
	}

	var accum rop.Result[In]
	for id, f := range fs {

		r := f(ctx)

		if !r.IsSuccess() {
			if validateF != nil {

				r = validateF(ctx, id, r)
				if !r.IsSuccess() {
					return rop.Fail[In](r.Err())
				}

			} else {
				return rop.Fail[In](r.Err())
			}
		}

		if accumF != nil {
			accum = accumF(ctx, id, r)
		} else {
			accum = r // last result
		}
	}
	return accum
}

// OrWithCtx TODO ctx done
func OrWithCtx[In any](ctx context.Context,
	validateF func(ctx context.Context, resId int, in rop.Result[In]) rop.Result[In],
	fs ...func(ctx context.Context) rop.Result[In]) rop.Result[In] {

	if len(fs) == 0 {
		return rop.Fail[In](fmt.Errorf(EmptyResult))
	}

	var err error = nil
	for id, f := range fs {

		r := f(ctx)

		if r.IsSuccess() {

			if validateF != nil {

				r = validateF(ctx, id, r)
				if r.IsSuccess() {
					return r
				}

			} else {
				return r
			}
		}

		err = rop.Iif(err == nil, r.Err(), errors.Join(err, r.Err()))
	}

	return rop.Fail[In](err)
}
