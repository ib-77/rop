package group

import (
	"context"
	"errors"
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

func AndTee[In any](input rop.Result[In],
	fs ...func(in rop.Result[In]) rop.Result[In]) rop.Result[In] {

	if !input.IsSuccess() {
		return input
	}

	res := input
	for _, f := range fs {

		res = f(res)

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

func AndSwitch[In, Out any](input rop.Result[In],
	switchF func(r rop.Result[In]) rop.Result[Out],
	fs ...func(in rop.Result[In]) rop.Result[In]) rop.Result[Out] {

	if !input.IsSuccess() {
		return switchF(input)
	}

	res := input
	for _, f := range fs {

		res = f(res)

		if !res.IsSuccess() {
			return rop.Fail[Out](res.Err())
		}
	}

	return switchF(res)
}

func OrTeeWithCtx[In any](ctx context.Context, input rop.Result[In],
	fs ...func(ctx context.Context, in rop.Result[In]) rop.Result[In]) rop.Result[In] {

	if !input.IsSuccess() {
		return input
	}

	var err error
	for _, f := range fs {

		r := f(ctx, input)

		if r.IsSuccess() {
			return r
		}

		if r.IsAccepted() {
			return r // error or cancel
		}

		err = rop.Iif(err == nil, r.Err(), errors.Join(err, r.Err()))
	}

	return rop.Fail[In](err)
}

func OrTee[In any](input rop.Result[In],
	fs ...func(in rop.Result[In]) rop.Result[In]) rop.Result[In] {

	if !input.IsSuccess() {
		return input
	}

	var err error
	for _, f := range fs {

		r := f(input)

		if r.IsSuccess() {
			return r
		}

		if r.IsAccepted() {
			return r // error or cancel
		}

		err = rop.Iif(err == nil, r.Err(), errors.Join(err, r.Err()))
	}

	return rop.Fail[In](err)
}

func OrSwitchWithCtx[In, Out any](ctx context.Context, input rop.Result[In],
	switchF func(ctx context.Context, r rop.Result[In]) rop.Result[Out],
	fs ...func(ctx context.Context, in rop.Result[In]) rop.Result[In]) rop.Result[Out] {

	if !input.IsSuccess() {
		return switchF(ctx, input)
	}

	var err error
	for _, f := range fs {

		r := f(ctx, input)

		if r.IsSuccess() {
			return switchF(ctx, r)
		}

		if r.IsAccepted() { // error or cancel

			if r.IsCancel() {
				return rop.Cancel[Out](r.Err())
			}

			return rop.Fail[Out](r.Err())
		}

		err = rop.Iif(err == nil, r.Err(), errors.Join(err, r.Err()))
	}

	return rop.Fail[Out](err)
}

func OrSwitch[In, Out any](input rop.Result[In],
	switchF func(r rop.Result[In]) rop.Result[Out],
	fs ...func(in rop.Result[In]) rop.Result[In]) rop.Result[Out] {

	if !input.IsSuccess() {
		return switchF(input)
	}

	var err error
	for _, f := range fs {

		r := f(input)

		if r.IsSuccess() {
			return switchF(r)
		}

		if r.IsAccepted() { // error or cancel

			if r.IsCancel() {
				return rop.Cancel[Out](r.Err())
			}

			return rop.Fail[Out](r.Err())
		}

		err = rop.Iif(err == nil, r.Err(), errors.Join(err, r.Err()))
	}

	return rop.Fail[Out](err)
}

func AcceptWithCtx[In any](ctx context.Context, input rop.Result[In],
	accept func(ctx context.Context, in rop.Result[In]) bool) rop.Result[In] {
	if accept(ctx, input) {
		return rop.Accept[In](input)
	}
	return input
}

func Accept[In any](input rop.Result[In], accept func(in rop.Result[In]) bool) rop.Result[In] {
	if accept(input) {
		return rop.Accept[In](input)
	}
	return input
}

func AcceptSuccessWithCtx[In any](ctx context.Context,
	input rop.Result[In], accept func(ctx context.Context, in In) bool) rop.Result[In] {

	if input.IsSuccess() {
		if accept(ctx, input.Result()) {
			return rop.Accept[In](input)
		}
	}
	return input
}

func AcceptSuccess[In any](input rop.Result[In], accept func(in In) bool) rop.Result[In] {

	if input.IsSuccess() {
		if accept(input.Result()) {
			return rop.Accept[In](input)
		}
	}
	return input
}

func AcceptFail[In any](input rop.Result[In], accept func(err error) bool) rop.Result[In] {

	if !input.IsSuccess() {
		if accept(input.Err()) {
			return rop.Accept[In](input)
		}
	}
	return input
}

func AcceptFailWithCtx[In any](ctx context.Context,
	input rop.Result[In], accept func(ctx context.Context, err error) bool) rop.Result[In] {

	if !input.IsSuccess() {
		if accept(ctx, input.Err()) {
			return rop.Accept[In](input)
		}
	}
	return input
}
