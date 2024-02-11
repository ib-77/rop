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

		assert.False(t, output.IsSuccess())
		assert.NotEmpty(t, output.Err())
		count++
	}
	assert.Equal(t, totalElements, count)
}

func successConvertIntToStrResult(r int) rop.Rop[string] {
	return rop.Success(strconv.Itoa(r))
}

func successConvertIntToStr(r int) (string, error) {
	return strconv.Itoa(r), nil
}

func failConvertIntToStrResult(r int) rop.Rop[string] {
	return rop.Fail[string](fmt.Errorf("cannot convert %d", r))
}

func failConvertIntToStr(r int) (string, err error) {
	err = fmt.Errorf("cannot convert %d", r)
	return
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
