package test

import (
	"context"
	"fmt"
	"github.com/ib-77/rop/pkg/rop"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func Test_MassValidateTrueUnbufferedInput(t *testing.T) {
	totalElements := 5
	ctx := context.Background()

	inputs := make(chan int)
	go func() {
		defer close(inputs)
		for i := 0; i < totalElements; i++ {
			inputs <- i
		}
	}()

	//outputs := rop.MassValidate(ctx, inputs, allSuccessInt, CancelInt, "error")
	//for {
	//
	//	validatedInt, ok := <-outputs
	//	if ok {
	//		assert.True(t, validatedInt.IsSuccess())
	//	} else {
	//		break
	//	}
	//}

	count := 0
	for output := range rop.MassValidate(ctx, inputs, allSuccessInt, CancelInt, "error") {
		assert.True(t, output.IsSuccess())
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassValidateTrueBufferedInput(t *testing.T) {
	totalElements := 5
	ctx := context.Background()

	inputs := make(chan int, totalElements)
	for i := 0; i < totalElements; i++ {
		inputs <- i
	}
	close(inputs) // close buffered after fill )

	count := 0
	for output := range rop.MassValidate(ctx, inputs, allSuccessInt, CancelInt, "error") {
		assert.True(t, output.IsSuccess())
		count++
	}
	assert.Equal(t, totalElements, count)
}

func Test_MassValidateWithCancel(t *testing.T) {
	log.SetFlags(0)

	totalElements := 5
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inputs := make(chan int)
	go func() {
		defer close(inputs)
		for i := 0; i < totalElements; i++ {
			inputs <- i
		}
	}()

	count := 0
	cancelOnCount := 1
	//cancel()
	for output := range rop.MassValidate(ctx, inputs, allSuccessInt, CancelInt, "error") {

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

func allSuccessInt(i int) bool {
	return true
}

func CancelInt(in int) error {
	return fmt.Errorf("processing of value %d was cancelled", in)
}
