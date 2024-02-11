package test

import (
	"context"
	"errors"
	"fmt"
	"github.com/ib-77/rop/pkg/rop"
)

func ropCase01(input int) string {
	return rop.Finally(
		rop.DoubleMap(
			rop.Map(
				rop.Tee(
					rop.Try(
						rop.Switch(
							rop.AndValidate(
								rop.Validate(input,
									lessTwo, "value more than 2"),
								notFive, "value is 5"),
							greaterThanZero),
						equalHundredOrThrowError),
					doAndForget),
				addChars),
			logSuccess, logFail),
		returnSuccessResult, returnFailResult)
}

func massRopCase01(ctx context.Context, inputs <-chan int) <-chan string {
	return rop.MassFinally(ctx,
		rop.MassDoubleMap(ctx,
			rop.MassMap(ctx,
				rop.MassTee(ctx,
					rop.MassTry(ctx,
						rop.MassSwitch(ctx,
							rop.MassAndValidate(ctx,
								rop.MassValidate(ctx, inputs,
									lessTwo, cancelF[int], "value more than 2"),
								notFive, cancelRopF[int], "value is 5"),
							greaterThanZero, cancelRopF[int]),
						equalHundredOrThrowError, cancelRopF[int]),
					doAndForget, cancelRopF[string]),
				addChars, cancelRopF[string]),
			logSuccess, logFail, cancelRopF[string]),
		returnSuccessResult, returnFailResult, cancelResultF[string])
}

func cancelF[T any](in T) error {
	return errors.New("some error")
}

func cancelRopF[T any](in rop.Rop[T]) error {
	return errors.New("some error")
}
func cancelResultF[T any](in rop.Rop[T]) string {
	return "some error"
}
func lessTwo(a int) bool {
	if a < 2 {
		return true
	}
	return false
}

func notFive(a int) bool {
	if a != 5 {
		return true
	}
	return false
}

func greaterThanZero(a int) rop.Rop[int] {
	if a > 0 {
		return rop.Success(100)
	}
	return rop.Fail[int](errors.New("a is less or 0!"))
}

func addChars(r string) string {
	return r + "fff"
}

func equalHundredOrThrowError(r int) (string, error) {
	if r == 100 {
		return "OK", nil
	}
	return "ER", errors.New("! 100")
}

func doAndForget(r rop.Rop[string]) {
	fmt.Printf("do something with 100!\n")
}

func logSuccess(r string) string {
	fmt.Printf("string: %s\n", r)
	return r
}

func logFail(er error) string {
	fmt.Printf("error: %s\n", er.Error())
	return er.Error()
}

func returnSuccessResult(r string) string {
	return "all ok"
}

func returnFailResult(er error) string {
	return fmt.Sprintf("error: %s", er)
}
