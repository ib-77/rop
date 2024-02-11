package test

import (
	"context"
	"fmt"
	"github.com/ib-77/rop/pkg/rop"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"strconv"
	"testing"
)

func init() {
	log.SetFlags(0)
}

func Test_MassValidate_True_UnbufferedInput(t *testing.T) {
	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range rop.MassValidate(ctx, inputs, allSuccess[int], CancelF[int], "error") {
		assert.True(t, output.IsSuccess())
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassValidate_True_BufferedInput(t *testing.T) {
	totalElements := 5
	ctx := context.Background()
	inputs := generateBufferedChan(totalElements, totalElements)
	count := 0

	for output := range rop.MassValidate(ctx, inputs, allSuccess[int], CancelF[int], "error") {
		assert.True(t, output.IsSuccess())
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassValidate_False(t *testing.T) {
	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range rop.MassValidate(ctx, inputs, allFail[int], CancelF[int], "error") {
		assert.False(t, output.IsSuccess())
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassValidate_WithCancel(t *testing.T) {
	totalElements := 5
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inputs := generateUnbufferedChan(totalElements)
	count := 0
	cancelOnCount := 1
	//cancel()
	for output := range rop.MassValidate(ctx, inputs, allSuccess[int], CancelF[int], "error") {

		if count <= cancelOnCount {
			assert.True(t, output.IsSuccess())
		} else {
			log.Println(fmt.Sprintf("is output (%d) cancelled: %v", count, output.IsCancel()))
		}

		if count == cancelOnCount {
			cancel()
		}

		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassValidateAndValidate_True(t *testing.T) {
	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range rop.MassAndValidate(ctx,
		rop.MassValidate(ctx, inputs, allSuccess[int], CancelF[int], "error"),
		allSuccess[int], CancelRopF[int], "and error") {
		assert.True(t, output.IsSuccess())
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassValidateAndValidate_False_AtFirst(t *testing.T) {
	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range rop.MassAndValidate(ctx,
		rop.MassValidate(ctx, inputs, allFail[int], CancelF[int], "error"),
		allSuccess[int], CancelRopF[int], "and error") {
		assert.False(t, output.IsSuccess())
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassValidate_AndValidate_False_AtSecond(t *testing.T) {
	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range rop.MassAndValidate(ctx,
		rop.MassValidate(ctx, inputs, allSuccess[int], CancelF[int], "error"),
		allFail[int], CancelRopF[int], "and error") {

		assert.False(t, output.IsSuccess())
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassValidate_AndValidate_WithCancel(t *testing.T) {
	totalElements := 5
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	inputs := generateUnbufferedChan(totalElements)
	count := 0
	cancelOnCount := 2

	for output := range rop.MassAndValidate(ctx,
		rop.MassValidate(ctx, inputs, allSuccess[int], CancelF[int], "error"),
		allFail[int], CancelRopF[int], "and error") {

		if count == cancelOnCount {
			cancel()
		}
		log.Println(fmt.Sprintf("is output (%d) cancelled: %v", count, output.IsCancel()))
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassSwitch_Success(t *testing.T) {
	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range rop.MassSwitch(ctx,
		rop.MassValidate(ctx, inputs, allSuccess[int], CancelF[int], "error"),
		successConvertIntToStrResult, CancelRopF[int]) {

		assert.True(t, output.IsSuccess())
		assert.Equal(t, reflect.TypeOf(output.Result()).Kind(), reflect.String)
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassSwitch_Fail(t *testing.T) {
	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range rop.MassSwitch(ctx,
		rop.MassValidate(ctx, inputs, allSuccess[int], CancelF[int], "error"),
		failConvertIntToStrResult, CancelRopF[int]) {

		assert.False(t, output.IsSuccess())
		assert.NotEmpty(t, output.Err())
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassMap_Success(t *testing.T) {
	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range rop.MassMap(ctx,
		rop.MassValidate(ctx, inputs, allSuccess[int], CancelF[int], "error"),
		successConvertIntToStr, CancelRopF[int]) {

		assert.True(t, output.IsSuccess())
		assert.Equal(t, reflect.TypeOf(output.Result()).Kind(), reflect.String)
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassMap_Fail(t *testing.T) {
	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range rop.MassMap(ctx,
		rop.MassValidate(ctx, inputs, allSuccess[int], CancelF[int], "error"),
		failConvertIntToStr, CancelRopF[int]) {

		assert.True(t, output.IsSuccess())
		assert.Equal(t, output.Result(), "cannot convert")
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassTry_Success(t *testing.T) {
	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range rop.MassTry(ctx,
		rop.MassValidate(ctx, inputs, allSuccess[int], CancelF[int], "error"),
		successConvertIntToStrWithErr, CancelRopF[int]) {

		assert.True(t, output.IsSuccess())
		assert.Equal(t, reflect.TypeOf(output.Result()).Kind(), reflect.String)
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassTry_Fail(t *testing.T) {
	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range rop.MassTry(ctx,
		rop.MassValidate(ctx, inputs, allSuccess[int], CancelF[int], "error"),
		failConvertIntToStrWithErr, CancelRopF[int]) {

		assert.False(t, output.IsSuccess())
		assert.Empty(t, output.Result())
		assert.NotEmpty(t, output.Err())
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassDoubleMap_Success(t *testing.T) {
	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range rop.MassDoubleMap(ctx,
		rop.MassValidate(ctx, inputs, allSuccess[int], CancelF[int], "error"),
		successConvertIntToStr, failConvertIntToStrProcessError, CancelRopF[int]) {

		assert.True(t, output.IsSuccess())
		assert.Equal(t, reflect.TypeOf(output.Result()).Kind(), reflect.String)
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassDoubleMap_Fail(t *testing.T) {
	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range rop.MassDoubleMap(ctx,
		rop.MassValidate(ctx, inputs, allFail[int], CancelF[int], "error"),
		successConvertIntToStr, failConvertIntToStrProcessError, CancelRopF[int]) {

		assert.False(t, output.IsSuccess())
		assert.NotEmpty(t, output.Err())
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassFinally_Success(t *testing.T) {
	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range rop.MassFinally(ctx,
		rop.MassValidate(ctx, inputs, allSuccess[int], CancelF[int], "error"),
		successFinally, failConvertIntToStrProcessError, CancelStrF[int]) {

		assert.Equal(t, "ok", output)
		assert.Equal(t, reflect.TypeOf(output).Kind(), reflect.String)
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassFinally_Fail(t *testing.T) {
	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range rop.MassFinally(ctx,
		rop.MassValidate(ctx, inputs, allFail[int], CancelF[int], "error"),
		successFinally, failConvertIntToStrProcessError, CancelStrF[int]) {

		assert.Equal(t, "error", output)
		assert.Equal(t, reflect.TypeOf(output).Kind(), reflect.String)
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassCase01_Success(t *testing.T) {
	totalElements := 5
	ctx := context.Background()
	inputs := generateFixedValueUnbufferedChan(totalElements, 1)
	count := 0

	for output := range massRopCase01(ctx, inputs) {

		assert.Equal(t, "all ok", output)
		assert.Equal(t, reflect.TypeOf(output).Kind(), reflect.String)
		count++
	}

	assert.Equal(t, totalElements, count)
}

func Test_MassCase01_Fail(t *testing.T) {
	totalElements := 5
	ctx := context.Background()
	inputs := generateFixedValueUnbufferedChan(totalElements, 3)
	count := 0

	for output := range massRopCase01(ctx, inputs) {

		assert.Equal(t, "error: value more than 2", output)
		assert.Equal(t, reflect.TypeOf(output).Kind(), reflect.String)
		count++
	}

	assert.Equal(t, totalElements, count)
}

func Test_MassCase01_FailZero(t *testing.T) {
	totalElements := 5
	ctx := context.Background()
	inputs := generateFixedValueUnbufferedChan(totalElements, 0)
	count := 0

	for output := range massRopCase01(ctx, inputs) {

		assert.Equal(t, "error: a is less or 0!", output)
		assert.Equal(t, reflect.TypeOf(output).Kind(), reflect.String)
		count++
	}

	assert.Equal(t, totalElements, count)
}

func successConvertIntToStrResult(r int) rop.Rop[string] {
	return rop.Success(strconv.Itoa(r))
}

func successConvertIntToStr(r int) string {
	return strconv.Itoa(r)
}

func successFinally(r int) string {
	return "ok"
}

func successConvertIntToStrWithErr(r int) (string, error) {
	return strconv.Itoa(r), nil
}

func failConvertIntToStrWithErr(r int) (re string, err error) {
	err = fmt.Errorf("cannot convert %d", r)
	return
}

func failConvertIntToStrResult(r int) rop.Rop[string] {
	return rop.Fail[string](fmt.Errorf("cannot convert %d", r))
}

func failConvertIntToStr(r int) string {
	return "cannot convert"
}

func failConvertIntToStrProcessError(err error) string {
	return "error"
}

// helpers
func generateUnbufferedChan(amount int) chan int {
	inputs := make(chan int)
	go func() {
		defer close(inputs)
		for i := 0; i < amount; i++ {
			inputs <- i
		}
	}()
	return inputs
}

func generateFixedValueUnbufferedChan(amount int, fixedValue int) chan int {
	inputs := make(chan int)
	go func() {
		defer close(inputs)
		for i := 0; i < amount; i++ {
			inputs <- fixedValue
		}
	}()
	return inputs
}

func generateBufferedChan(amount int, bufferSize int) chan int {
	inputs := make(chan int, bufferSize)
	defer close(inputs)

	for i := 0; i < amount; i++ {
		inputs <- i
	}

	return inputs
}

func allSuccess[T any](i T) bool {
	return true
}

func allFail[T any](i T) bool {
	return false
}

func CancelF[T any](in T) error {
	return fmt.Errorf("---- processing of value %v was cancelled", in)
}

func CancelRopF[T any](in rop.Rop[T]) error {
	return fmt.Errorf("---- processing of value %v was cancelled", in)
}

func CancelStrF[In any](in rop.Rop[In]) string {
	return fmt.Sprintf("---- processing of value %v was cancelled", in)
}
