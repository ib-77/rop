package group

import (
	"context"
	"github.com/ib-77/rop/pkg/rop"
)

func AndTeeWithCtx[In any](ctx context.Context, input rop.Result[In],
	fs ...func(ctx context.Context, in rop.Result[In]) rop.Result[In]) rop.Result[In] {

	if !input.IsSuccess() {
		return input
	}

	res := input
	for _, f := range fs {

		res = f(ctx, res)

		if !res.IsSuccess() {
			return rop.Fail[In](res.Err())
		}
	}

	return res
}

func AndSwitchWithCtx[In, Out any](ctx context.Context, input rop.Result[In],
	switchF func(ctx context.Context, r rop.Result[In]) rop.Result[Out],
	fs ...func(ctx context.Context, in rop.Result[In]) rop.Result[In]) rop.Result[Out] {

	if !input.IsSuccess() {
		return switchF(ctx, input)
	}

	res := input
	for _, f := range fs {

		res = f(ctx, res)

		if !res.IsSuccess() {
			return rop.Fail[Out](res.Err())
		}
	}

	return switchF(ctx, res)
}

func OrTeeWithCtx[In any](ctx context.Context, input rop.Result[In],
	fs ...func(ctx context.Context, in rop.Result[In]) (accepted bool, res rop.Result[In])) rop.Result[In] {

	if !input.IsSuccess() {
		return input
	}

	if len(fs) == 0 {
		return input
	}

	for _, f := range fs {

		accepted, r := f(ctx, input)

		if accepted {
			if r.IsSuccess() {
				return r
			}
			return rop.Fail[In](r.Err())
		}

	}

	return input
}

func OrSwitchWithCtx[In, Out any](ctx context.Context, input rop.Result[In],
	switchF func(ctx context.Context, r rop.Result[In]) rop.Result[Out],
	fs ...func(ctx context.Context, in rop.Result[In]) (accepted bool, res rop.Result[In])) rop.Result[Out] {

	if !input.IsSuccess() {
		return switchF(ctx, input)
	}

	if len(fs) == 0 {
		return switchF(ctx, input)
	}

	for _, f := range fs {

		accepted, r := f(ctx, input)

		if accepted {

			if r.IsSuccess() {
				return switchF(ctx, r)
			}

			return rop.Fail[Out](r.Err())
		}
	}

	return switchF(ctx, input)
}
