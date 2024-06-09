package mass

import (
	"context"
	"fmt"
	"github.com/ib-77/rop/pkg/rop"
	"github.com/ib-77/rop/pkg/rop/mass"
	"github.com/ib-77/rop/test"
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
	t.Parallel()

	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range mass.Validate(ctx, inputs, allSuccess[int], CancelF[int], "error") {
		assert.True(t, output.IsSuccess())
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassValidate_True_BufferedInput(t *testing.T) {
	t.Parallel()

	totalElements := 5
	ctx := context.Background()
	inputs := generateBufferedChan(totalElements, totalElements)
	count := 0

	for output := range mass.Validate(ctx, inputs, allSuccess[int], CancelF[int], "error") {
		assert.True(t, output.IsSuccess())
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassValidate_True_BufferedInput_CloseInput(t *testing.T) {
	t.Parallel()

	totalElements := 5
	closeOn := 2
	ctx := context.Background()
	inputs := generateBufferedChanCloseOn(totalElements, totalElements, closeOn)
	count := 0

	for output := range mass.Validate(ctx, inputs, allSuccess[int], CancelF[int], "error") {
		assert.True(t, output.IsSuccess())
		count++
	}
	assert.Equal(t, closeOn, count)
}

func Test_MassValidate_False(t *testing.T) {
	t.Parallel()

	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range mass.Validate(ctx, inputs, allFail[int], CancelF[int], "error") {
		assert.False(t, output.IsSuccess())
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassValidate_WithCancel(t *testing.T) {
	t.Parallel()

	totalElements := 5
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inputs := generateUnbufferedChan(totalElements)
	count := 0
	cancelOnCount := 1
	//cancel()
	for output := range mass.Validate(ctx, inputs, allSuccess[int], CancelF[int], "error") {

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
	t.Parallel()

	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range mass.AndValidate(ctx,
		mass.Validate(ctx, inputs, allSuccess[int], CancelF[int], "error"),
		allSuccess[int], CancelRopF[int], "and error") {
		assert.True(t, output.IsSuccess())
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassValidateAndValidate_False_AtFirst(t *testing.T) {
	t.Parallel()

	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range mass.AndValidate(ctx,
		mass.Validate(ctx, inputs, allFail[int], CancelF[int], "error"),
		allSuccess[int], CancelRopF[int], "and error") {
		assert.False(t, output.IsSuccess())
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassValidate_AndValidate_False_AtSecond(t *testing.T) {
	t.Parallel()

	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range mass.AndValidate(ctx,
		mass.Validate(ctx, inputs, allSuccess[int], CancelF[int], "error"),
		allFail[int], CancelRopF[int], "and error") {

		assert.False(t, output.IsSuccess())
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassValidate_AndValidate_WithCancel(t *testing.T) {
	t.Parallel()

	totalElements := 5
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	inputs := generateUnbufferedChan(totalElements)
	count := 0
	cancelOnCount := 2

	for output := range mass.AndValidate(ctx,
		mass.Validate(ctx, inputs, allSuccess[int], CancelF[int], "error"),
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
	t.Parallel()

	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range mass.Switch(ctx,
		mass.Validate(ctx, inputs, allSuccess[int], CancelF[int], "error"),
		successConvertIntToStrResult, CancelRopF[int]) {

		assert.True(t, output.IsSuccess())
		assert.Equal(t, reflect.TypeOf(output.Result()).Kind(), reflect.String)
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassSwitch_Fail(t *testing.T) {
	t.Parallel()

	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range mass.Switch(ctx,
		mass.Validate(ctx, inputs, allSuccess[int], CancelF[int], "error"),
		failConvertIntToStrResult, CancelRopF[int]) {

		assert.False(t, output.IsSuccess())
		assert.NotEmpty(t, output.Err())
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassMap_Success(t *testing.T) {
	t.Parallel()

	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range mass.Map(ctx,
		mass.Validate(ctx, inputs, allSuccess[int], CancelF[int], "error"),
		successConvertIntToStr, CancelRopF[int]) {

		assert.True(t, output.IsSuccess())
		assert.Equal(t, reflect.TypeOf(output.Result()).Kind(), reflect.String)
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassMap_Fail(t *testing.T) {
	t.Parallel()

	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range mass.Map(ctx,
		mass.Validate(ctx, inputs, allSuccess[int], CancelF[int], "error"),
		failConvertIntToStr, CancelRopF[int]) {

		assert.True(t, output.IsSuccess())
		assert.Equal(t, output.Result(), "cannot convert")
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassTry_Success(t *testing.T) {
	t.Parallel()

	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range mass.Try(ctx,
		mass.Validate(ctx, inputs, allSuccess[int], CancelF[int], "error"),
		successConvertIntToStrWithErr, CancelRopF[int]) {

		assert.True(t, output.IsSuccess())
		assert.Equal(t, reflect.TypeOf(output.Result()).Kind(), reflect.String)
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassTry_Fail(t *testing.T) {
	t.Parallel()

	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range mass.Try(ctx,
		mass.Validate(ctx, inputs, allSuccess[int], CancelF[int], "error"),
		failConvertIntToStrWithErr, CancelRopF[int]) {

		assert.False(t, output.IsSuccess())
		assert.Empty(t, output.Result())
		assert.NotEmpty(t, output.Err())
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassDoubleMap_Success(t *testing.T) {
	t.Parallel()

	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range mass.DoubleMap(ctx,
		mass.Validate(ctx, inputs, allSuccess[int], CancelF[int], "error"),
		successConvertIntToStr, failConvertIntToStrProcessError,
		cancelConvertIntToStrProcessError, CancelRopF[int]) {

		assert.True(t, output.IsSuccess())
		assert.Equal(t, reflect.TypeOf(output.Result()).Kind(), reflect.String)
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassDoubleMap_Fail(t *testing.T) {
	t.Parallel()

	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range mass.DoubleMap(ctx,
		mass.Validate(ctx, inputs, allFail[int], CancelF[int], "error"),
		successConvertIntToStr, failConvertIntToStrProcessError, cancelConvertIntToStrProcessError,
		CancelRopF[int]) {

		assert.False(t, output.IsSuccess())
		assert.NotEmpty(t, output.Err())
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassFinally_Success(t *testing.T) {
	t.Parallel()

	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range mass.Finally(ctx,
		mass.Validate(ctx, inputs, allSuccess[int], CancelF[int], "error"),
		successFinally, failConvertIntToStrProcessError, CancelStrF[int]) {

		assert.Equal(t, "ok", output)
		assert.Equal(t, reflect.TypeOf(output).Kind(), reflect.String)
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassFinally_Fail(t *testing.T) {
	t.Parallel()

	totalElements := 5
	ctx := context.Background()
	inputs := generateUnbufferedChan(totalElements)
	count := 0

	for output := range mass.Finally(ctx,
		mass.Validate(ctx, inputs, allFail[int], CancelF[int], "error"),
		successFinally, failConvertIntToStrProcessError, CancelStrF[int]) {

		assert.Equal(t, "error", output)
		assert.Equal(t, reflect.TypeOf(output).Kind(), reflect.String)
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassCase01_Success(t *testing.T) {
	t.Parallel()

	totalElements := 5
	ctx := context.Background()
	inputs := generateFixedValueUnbufferedChan(totalElements, 1)
	count := 0

	for output := range test.MassRopCase01(ctx, inputs) {

		assert.Equal(t, "all ok", output)
		assert.Equal(t, reflect.TypeOf(output).Kind(), reflect.String)
		count++
	}

	assert.Equal(t, totalElements, count)
}

func Test_MassCase01_Fail(t *testing.T) {
	t.Parallel()

	totalElements := 5
	ctx := context.Background()
	inputs := generateFixedValueUnbufferedChan(totalElements, 3)
	count := 0

	for output := range test.MassRopCase01(ctx, inputs) {

		assert.Equal(t, "error: value more than 2", output)
		assert.Equal(t, reflect.TypeOf(output).Kind(), reflect.String)
		count++
	}

	assert.Equal(t, totalElements, count)
}

func Test_MassCase01_FailZero(t *testing.T) {
	t.Parallel()

	totalElements := 5
	ctx := context.Background()
	inputs := generateFixedValueUnbufferedChan(totalElements, 0)
	count := 0

	for output := range test.MassRopCase01(ctx, inputs) {

		assert.Equal(t, "error: a is less or 0!", output)
		assert.Equal(t, reflect.TypeOf(output).Kind(), reflect.String)
		count++
	}

	assert.Equal(t, totalElements, count)
}

func successConvertIntToStrResult(_ context.Context, r int) rop.Result[string] {
	return rop.Success(strconv.Itoa(r))
}

func successConvertIntToStr(_ context.Context, r int) string {
	return strconv.Itoa(r)
}

func successFinally(_ context.Context, r int) string {
	return "ok"
}

func successConvertIntToStrWithErr(_ context.Context, r int) (string, error) {
	return strconv.Itoa(r), nil
}

func failConvertIntToStrWithErr(_ context.Context, r int) (re string, err error) {
	err = fmt.Errorf("cannot convert %d", r)
	return
}

func failConvertIntToStrResult(_ context.Context, r int) rop.Result[string] {
	return rop.Fail[string](fmt.Errorf("cannot convert %d", r))
}

func failConvertIntToStr(_ context.Context, r int) string {
	return "cannot convert"
}

func failConvertIntToStrProcessError(_ context.Context, err error) string {
	return "error"
}

func cancelConvertIntToStrProcessError(_ context.Context, err error) string {
	return "cancelled before"
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

func generateBufferedChanCloseOn(amount int, bufferSize int, closeOn int) chan int {
	inputs := make(chan int, bufferSize)
	closed := &struct{}{}
	defer func() {
		if closed != nil {
			close(inputs)
		}
	}()

	for i := 0; i < amount; i++ {

		if i == closeOn {
			close(inputs)
			closed = nil
			break
		}
		inputs <- i
	}

	return inputs
}

func allSuccess[T any](ctx context.Context, i T) bool {
	return true
}

func allFail[T any](ctx context.Context, i T) bool {
	return false
}

func CancelF[T any](ctx context.Context, in T) error {
	return fmt.Errorf("---- processing of value %v was cancelled", in)
}

func CancelRopF[T any](ctx context.Context, in T) error {
	return fmt.Errorf("---- processing of value %v was cancelled", in)
}

func CancelStrF[In any](_ context.Context, in rop.Result[In]) string {
	return fmt.Sprintf("---- processing of value %v was cancelled", in)
}
