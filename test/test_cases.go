package test

import (
	"context"
	"errors"
	"fmt"
	"github.com/ib-77/rop/pkg/rop"
	"github.com/ib-77/rop/pkg/rop/group"
	"github.com/ib-77/rop/pkg/rop/mass"
	"github.com/ib-77/rop/pkg/rop/solo"
)

func RopCase01(input int) string {
	return solo.Finally(
		solo.DoubleMap(
			solo.Map(
				solo.Tee(
					solo.Try(
						solo.Switch(
							solo.AndValidate(
								solo.Validate(input,
									lessTwo, "value more than 2"),
								notFive, "value is 5"),
							greaterThanZero),
						equalHundredOrThrowError),
					doAndForget),
				addChars),
			logSuccess, logFail, logCancel),
		returnSuccessResult, returnFailResult)
}

func RopCase02(input int) string {
	return solo.Finally(
		solo.DoubleMap(
			solo.Map(
				solo.Tee(
					solo.Try(
						solo.Switch(
							solo.AndValidateWithErr(
								solo.ValidateWithErr(input,
									lessTwoErr),
								notFiveErr),
							greaterThanZero),
						equalHundredOrThrowError),
					doAndForget),
				addChars),
			logSuccess, logFail, logCancel),
		returnSuccessResult, returnFailResult)
}

func RopCase03(ctx context.Context, input int) string {

	return solo.FinallyWithCtx(ctx,
		solo.DoubleMapWithCtx(ctx,
			solo.MapWithCtx(ctx,
				solo.TeeWithCtx(ctx,
					solo.TryWithCtx(ctx,
						solo.SwitchWithCtx(ctx,
							solo.AndValidateWithErrWithCtx(ctx,
								solo.ValidateWithErrWithCtx(ctx, input,
									lessTwoErrCtx),
								notFiveErrCtx),
							greaterThanZeroCtx),
						equalHundredOrThrowErrorCtx),
					doAndForgetCtx),
				addCharsCtx),
			logSuccessCtx, logFailCtx, logCancelCtx),
		returnSuccessResultCtx, returnFailResultCtx)
}

// TODO review test!
func RopCase04(ctx context.Context, input int) string {

	return solo.FinallyWithCtx(ctx,
		group.OrTeeWithCtx(ctx,
			rop.Success(input),
			func(ctx context.Context, in rop.Result[int]) (bool, rop.Result[int]) {

				accepted, err := lessTwoErrCtx(ctx, in.Result())
				if err != nil {
					return false, rop.Fail[int](err)
				}
				if !accepted {
					return false, rop.Cancel[int](fmt.Errorf("canceled 1"))
				}

				return true, solo.SwitchWithCtx(ctx, in, greaterThanZeroCtx)
			},
			func(ctx context.Context, in rop.Result[int]) (bool, rop.Result[int]) {

				accepted, err := lessTwoErrCtx(ctx, in.Result())
				if err != nil {
					return false, rop.Fail[int](err)
				}
				if !accepted {
					return false, rop.Cancel[int](fmt.Errorf("canceled 2"))
				}

				return true, solo.SwitchWithCtx(ctx, in, greaterThanZeroCtx)
			}),
		returnSuccessResultIntValueCtx, returnFailResultCtx)
}

func RopBenchCase01(input int) string {
	return solo.Finally(
		solo.DoubleMap(
			solo.Map(
				solo.Tee(
					solo.Try(
						solo.Switch(
							solo.AndValidate(
								solo.Validate(input,
									lessTwo, "value more than 2"),
								notFive, "value is 5"),
							greaterThanZero),
						equalHundredOrThrowError),
					doAndForgetNoFormat),
				addChars),
			logSuccessNoFormat, logFailNoFormat, logCancelNoFormat),
		returnSuccessResult, returnFailResultNoFormat)
}

func MassRopCase01(ctx context.Context, inputs <-chan int) <-chan string {
	return mass.Finally(ctx,
		mass.DoubleMap(ctx,
			mass.Map(ctx,
				mass.Tee(ctx,
					mass.Try(ctx,
						mass.Switch(ctx,
							mass.AndValidate(ctx,
								mass.Validate(ctx, inputs,
									lessTwoCtx, cancelF[int], "value more than 2"),
								notFiveCtx, cancelRopF[int], "value is 5"),
							greaterThanZeroCtx, cancelRopF[int]),
						equalHundredOrThrowErrorCtx, cancelRopF[int]),
					doAndForgetCtx, cancelRopF[string]),
				addCharsCtx, cancelRopF[string]),
			logSuccessCtx, logFailCtx, logCancelCtx, cancelRopF[string]),
		returnSuccessResultCtx, returnFailResultCtx, cancelResultF[string])
}

func cancelF[T any](_ context.Context, in T) error {
	return errors.New("some error")
}

func cancelRopF[T any](_ context.Context, in T) error {
	return errors.New("some error")
}

func cancelResultF[T any](_ context.Context, in rop.Result[T]) string {
	return "some error"
}

func lessTwo(a int) bool {
	if a < 2 {
		return true
	}
	return false
}

func lessTwoCtx(_ context.Context, a int) bool {
	if a < 2 {
		return true
	}
	return false
}
func lessTwoErr(a int) (bool, error) {
	if a < 2 {
		return true, nil
	}
	return false, errors.New("value more than 2")
}
func lessTwoErrCtx(_ context.Context, a int) (bool, error) {
	if a < 2 {
		return true, nil
	}
	return false, errors.New("value more than 2")
}
func notFiveErr(a int) (bool, error) {
	if a != 5 {
		return true, nil
	}
	return false, errors.New("value is 5")
}

func notFiveErrCtx(_ context.Context, a int) (bool, error) {
	if a != 5 {
		return true, nil
	}
	return false, errors.New("value is 5")
}

func notFive(a int) bool {
	if a != 5 {
		return true
	}
	return false
}
func notFiveCtx(_ context.Context, a int) bool {
	if a != 5 {
		return true
	}
	return false
}
func greaterThanZero(a int) rop.Result[int] {
	if a > 0 {
		return rop.Success(100)
	}
	return rop.Fail[int](errors.New("a is less or 0!"))
}

func greaterThanZeroCtx(_ context.Context, a int) rop.Result[int] {
	if a > 0 {
		return rop.Success(100)
	}
	return rop.Fail[int](errors.New("a is less or 0!"))
}

func sumInt(rs ...rop.Result[int]) rop.Result[int] {
	sum := 0
	for _, s := range rs {
		sum += s.Result()
	}
	return rop.Success(sum)
}

func addChars(r string) string {
	return r + "fff"
}
func addCharsCtx(_ context.Context, r string) string {
	return r + "fff"
}
func equalHundredOrThrowError(r int) (string, error) {
	if r == 100 {
		return "OK", nil
	}
	return "ER", errors.New("! 100")
}
func equalHundredOrThrowErrorCtx(_ context.Context, r int) (string, error) {
	if r == 100 {
		return "OK", nil
	}
	return "ER", errors.New("! 100")
}
func doAndForget(r rop.Result[string]) {
	fmt.Printf("do something with 100!\n")
}

func doAndForgetCtx(_ context.Context, r rop.Result[string]) {
	fmt.Printf("do something with 100!\n")
}

func doAndForgetIntCtx(_ context.Context, r rop.Result[int]) {
	fmt.Printf("do something with 100!\n")
}

func doAndForgetNoFormat(r rop.Result[string]) {
}

func logSuccess(r string) string {
	fmt.Printf("string: %s\n", r)
	return r
}
func logSuccessCtx(_ context.Context, r string) string {
	fmt.Printf("string: %s\n", r)
	return r
}
func logSuccessNoFormat(r string) string {
	return r
}

func logFail(er error) string {
	fmt.Printf("error: %s\n", er.Error())
	return er.Error()
}
func logFailCtx(_ context.Context, er error) string {
	fmt.Printf("error: %s\n", er.Error())
	return er.Error()
}
func logFailNoFormat(er error) string {
	return er.Error()
}
func logFailNoFormatCtx(_ context.Context, er error) string {
	return er.Error()
}
func logCancel(er error) string {
	fmt.Printf("cancel: %s\n", er.Error())
	return er.Error()
}
func logCancelCtx(_ context.Context, er error) string {
	fmt.Printf("cancel: %s\n", er.Error())
	return er.Error()
}

func logCancelNoFormat(er error) string {
	return er.Error()
}

func returnSuccessResult(r string) string {
	return "all ok"
}
func returnSuccessResultCtx(_ context.Context, r string) string {
	return "all ok"
}

func returnSuccessResultIntCtx(_ context.Context, r int) string {
	return "all ok"
}
func returnSuccessResultIntValueCtx(_ context.Context, r int) string {
	return fmt.Sprintf("all ok %d", r)
}

func returnFailResult(er error) string {
	return fmt.Sprintf("error: %s", er.Error())
}
func returnFailResultCtx(_ context.Context, er error) string {
	return fmt.Sprintf("error: %s", er.Error())
}

func returnFailResultNoFormat(er error) string {
	return er.Error()
}
func returnFailResultNoFormatCtx(_ context.Context, er error) string {
	return er.Error()
}
