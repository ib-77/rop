package bridge

import (
	"context"
	"errors"
	"fmt"
	"github.com/ib-77/rop/pkg/rop"
	"github.com/ib-77/rop/pkg/rop/bridge"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestValidate(t *testing.T) {
	inputs := make(chan chan int, 2)
	input1 := make(chan int, 2)
	input1 <- 1
	input1 <- 2
	input2 := make(chan int, 3)
	input2 <- 5
	input2 <- 6
	input2 <- 7
	inputs <- input1
	inputs <- input2

	ctx := context.Background()
	outputs := bridge.Validate(ctx, inputs, validate, cancel, "err")

	assert.NotEmpty(t, outputs)
	output1 := <-outputs
	//	assert.NotEmpty(t, output1)
	assert.Equal(t, rop.Success(1), <-output1)
	assert.Equal(t, rop.Success(2), <-output1)

	output2 := <-outputs
	//	assert.NotEmpty(t, output2)
	assert.Equal(t, rop.Success(5), <-output2)
	assert.Equal(t, rop.Success(6), <-output2)
	assert.Equal(t, rop.Success(7), <-output2)
}

func TestValidateAndValidate(t *testing.T) {
	inputs := make(chan chan int, 2)
	input1 := make(chan int, 2)
	input1 <- 1
	input1 <- 2
	input2 := make(chan int, 3)
	input2 <- 5
	input2 <- 6
	input2 <- 7
	inputs <- input1
	inputs <- input2

	ctx := context.Background()
	outputs := bridge.AndValidate(ctx,
		bridge.Validate(ctx, inputs, validate, cancel, "err"),
		validateT, cancelT, "and err")

	assert.NotEmpty(t, outputs)
	output1 := <-outputs
	//	assert.NotEmpty(t, output1)
	assert.Equal(t, rop.Success(1), <-output1)
	assert.Equal(t, rop.Success(2), <-output1)

	output2 := <-outputs
	//	assert.NotEmpty(t, output2)
	assert.Equal(t, rop.Success(5), <-output2)
	assert.Equal(t, rop.Success(6), <-output2)
	assert.Equal(t, rop.Success(7), <-output2)
}

func TestValidateWithErrAndValidate(t *testing.T) {
	inputs := make(chan chan int, 2)
	input1 := make(chan int, 2)
	input1 <- 1
	input1 <- 2
	input2 := make(chan int, 3)
	input2 <- 5
	input2 <- 6
	input2 <- 7
	inputs <- input1
	inputs <- input2

	ctx := context.Background()

	outputs := bridge.AndValidate(ctx,
		bridge.Validate(ctx, inputs, validateWithError, cancel, "err"),
		validateT, cancelT, "and err")

	assert.NotEmpty(t, outputs)
	output1 := <-outputs
	//	assert.NotEmpty(t, output1)
	assert.Equal(t, rop.Success(1), <-output1)
	assert.Equal(t, rop.Fail[int](errors.New("err")), <-output1)

	output2 := <-outputs
	//	assert.NotEmpty(t, output2)
	assert.Equal(t, rop.Success(5), <-output2)
	assert.Equal(t, rop.Success(6), <-output2)
	assert.Equal(t, rop.Success(7), <-output2)
}

func TestValidateWithCancelAndValidate(t *testing.T) {
	inputs := make(chan chan int, 2)
	input1 := make(chan int, 2)
	input1 <- 1
	input1 <- 2
	input2 := make(chan int, 3)
	input2 <- 5
	input2 <- 6
	input2 <- 7
	inputs <- input1
	inputs <- input2

	ctx, cancellation := context.WithTimeout(context.Background(), time.Second*3)
	defer cancellation()

	outputs := bridge.AndValidate(ctx,
		bridge.Validate(ctx, inputs, validateLongRunning, cancel, "err"),
		validateT, cancelT, "and err")

	assert.NotEmpty(t, outputs)
	output1 := <-outputs
	assert.Equal(t, rop.Success(1), <-output1)
	assert.Equal(t, rop.Cancel[int](errors.New("and operation was cancelled 2")), <-output1)

	output2 := <-outputs
	assert.Equal(t, rop.Success(5), <-output2)
	assert.Equal(t, rop.Cancel[int](errors.New("and operation was cancelled 6")), <-output2)
	assert.Equal(t, rop.Cancel[int](errors.New("operation was cancelled 7")), <-output2)
}

func TestValidateAndValidateAndFinallySuccess(t *testing.T) {
	inputs := make(chan chan int, 2)
	input1 := make(chan int, 2)
	input1 <- 1
	input1 <- 2
	input2 := make(chan int, 3)
	input2 <- 5
	input2 <- 6
	input2 <- 7
	inputs <- input1
	inputs <- input2

	ctx := context.Background()

	outputs := bridge.Finally(ctx,
		bridge.AndValidate(ctx,
			bridge.Validate(ctx, inputs, validateNoDelay, cancel, "err"),
			validateNoDelay, cancelT, "and err"), successFinally, failFinally, cancelFinally)

	assert.NotEmpty(t, outputs)
	output1 := <-outputs
	assert.Equal(t, "success 1", <-output1)
	assert.Equal(t, "success 2", <-output1)

	output2 := <-outputs
	assert.Equal(t, "success 5", <-output2)
	assert.Equal(t, "success 6", <-output2)
	assert.Equal(t, "success 7", <-output2)
}

func validate(ctx context.Context, in int) bool {
	time.Sleep(time.Second * 5)
	//fmt.Println(" >>", in)
	return true
}

func validateNoDelay(ctx context.Context, in int) bool {
	return true
}

func validateLongRunning(ctx context.Context, in int) bool {

	if in == 1 || in == 5 {
		time.Sleep(time.Second)
	} else {
		time.Sleep(time.Second * 5)
	}
	return true
}

func validateT(ctx context.Context, in int) bool {
	//time.Sleep(time.Second)
	//fmt.Println("and >>", in)
	return true
}

func validateWithError(ctx context.Context, in int) bool {
	time.Sleep(time.Second)
	if in == 2 {
		return false
	}
	//fmt.Println("and >>", in)
	return true
}

func cancel(ctx context.Context, in int) error {
	return fmt.Errorf("operation was cancelled %d", in)
}

func cancelT(ctx context.Context, in int) error {
	return fmt.Errorf("and operation was cancelled %v", in)
}

func successFinally(ctx context.Context, in int) string {
	return fmt.Sprintf("success %d", in)
}

func failFinally(ctx context.Context, err error) string {
	return fmt.Sprintf("fail %v", err)
}

func cancelFinally(ctx context.Context, r rop.Result[int]) string {
	return fmt.Sprintf("cancel %v", r)
}
