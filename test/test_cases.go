package test

import (
	"context"
	"errors"
	"fmt"
	"github.com/ib-77/rop/pkg/rop"
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

func MassRopCase01(ctx context.Context, inputs <-chan int) <-chan string {
	return mass.Finally(ctx,
		mass.DoubleMap(ctx,
			mass.Map(ctx,
				mass.Tee(ctx,
					mass.Try(ctx,
						mass.Switch(ctx,
							mass.AndValidate(ctx,
								mass.Validate(ctx, inputs,
									lessTwo, cancelF[int], "value more than 2"),
								notFive, cancelRopF[int], "value is 5"),
							greaterThanZero, cancelRopF[int]),
						equalHundredOrThrowError, cancelRopF[int]),
					doAndForget, cancelRopF[string]),
				addChars, cancelRopF[string]),
			logSuccess, logFail, logCancel, cancelRopF[string]),
		returnSuccessResult, returnFailResult, cancelResultF[string])
}

func cancelF[T any](in T) error {
	return errors.New("some error")
}

func cancelRopF[T any](in rop.Result[T]) error {
	return errors.New("some error")
}
func cancelResultF[T any](in rop.Result[T]) string {
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

func greaterThanZero(a int) rop.Result[int] {
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

func doAndForget(r rop.Result[string]) {
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

func logCancel(er error) string {
	fmt.Printf("cancel: %s\n", er.Error())
	return er.Error()
}

func returnSuccessResult(r string) string {
	return "all ok"
}

func returnFailResult(er error) string {
	return fmt.Sprintf("error: %s", er)
}
