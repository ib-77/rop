package solo

import (
	"context"
	"errors"
	"github.com/ib-77/rop/pkg/rop"
	"github.com/ib-77/rop/pkg/rop/solo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTeeWithErr(t *testing.T) {
	// Test success case with no error from deadEndF
	result := rop.Success(5)
	called := false
	deadEndF := func(r rop.Result[int]) error {
		called = true
		assert.Equal(t, 5, r.Result())
		return nil
	}

	output := solo.TeeWithErr(result, deadEndF)
	assert.True(t, called)
	assert.True(t, output.IsSuccess())
	assert.Equal(t, 5, output.Result())

	// Test success case with error from deadEndF
	result = rop.Success(10)
	called = false
	expectedErr := errors.New("deadend error")
	deadEndF = func(r rop.Result[int]) error {
		called = true
		return expectedErr
	}

	output = solo.TeeWithErr(result, deadEndF)
	assert.True(t, called)
	assert.False(t, output.IsSuccess())
	assert.Equal(t, expectedErr, output.Err())

	// Test fail case
	failErr := errors.New("fail error")
	result = rop.Fail[int](failErr)
	called = false
	deadEndF = func(r rop.Result[int]) error {
		called = true
		return nil
	}

	output = solo.TeeWithErr(result, deadEndF)
	assert.False(t, called)
	assert.False(t, output.IsSuccess())
	assert.Equal(t, failErr, output.Err())
}

func TestDoubleTee(t *testing.T) {
	// Test success case
	result := rop.Success(5)
	successCalled := false
	errorCalled := false

	deadEndF := func(r int) {
		successCalled = true
		assert.Equal(t, 5, r)
	}

	deadEndWithErrF := func(err error) {
		errorCalled = true
	}

	output := solo.DoubleTee(result, deadEndF, deadEndWithErrF)
	assert.True(t, successCalled)
	assert.False(t, errorCalled)
	assert.True(t, output.IsSuccess())
	assert.Equal(t, 5, output.Result())

	// Test fail case
	failErr := errors.New("fail error")
	result = rop.Fail[int](failErr)
	successCalled = false
	errorCalled = false

	deadEndF = func(r int) {
		successCalled = true
	}

	deadEndWithErrF = func(err error) {
		errorCalled = true
		assert.Equal(t, failErr, err)
	}

	output = solo.DoubleTee(result, deadEndF, deadEndWithErrF)
	assert.False(t, successCalled)
	assert.True(t, errorCalled)
	assert.False(t, output.IsSuccess())
	assert.Equal(t, failErr, output.Err())
}

func TestCheck(t *testing.T) {
	// Test success case with true result
	result := rop.Success(5)
	boolF := func(r int) bool {
		return r > 3
	}

	output := solo.Check(result, boolF, "value not greater than 3")
	assert.True(t, output.IsSuccess())
	assert.Equal(t, true, output.Result())

	// Test success case with false result
	result = rop.Success(2)
	output = solo.Check(result, boolF, "value not greater than 3")
	assert.False(t, output.IsSuccess())
	assert.Equal(t, "value not greater than 3", output.Err().Error())

	// Test fail case
	failErr := errors.New("fail error")
	result = rop.Fail[int](failErr)
	output = solo.Check(result, boolF, "value not greater than 3")
	assert.False(t, output.IsSuccess())
	assert.Equal(t, failErr, output.Err())

	// Test cancel case
	cancelErr := errors.New("cancel error")
	result = rop.Cancel[int](cancelErr)
	output = solo.Check(result, boolF, "value not greater than 3")
	assert.False(t, output.IsSuccess())
	assert.True(t, output.IsCancel())
	assert.Equal(t, cancelErr, output.Err())
}

func TestCheckCancel(t *testing.T) {
	// Test success case with true result
	result := rop.Success(5)
	boolF := func(r int) bool {
		return r > 3
	}

	output := solo.CheckCancel(result, boolF, "value not greater than 3")
	assert.True(t, output.IsSuccess())
	assert.Equal(t, true, output.Result())

	// Test success case with false result (should cancel)
	result = rop.Success(2)
	output = solo.CheckCancel(result, boolF, "value not greater than 3")
	assert.False(t, output.IsSuccess())
	assert.True(t, output.IsCancel())
	assert.Equal(t, "value not greater than 3", output.Err().Error())

	// Test fail case
	failErr := errors.New("fail error")
	result = rop.Fail[int](failErr)
	output = solo.CheckCancel(result, boolF, "value not greater than 3")
	assert.False(t, output.IsSuccess())
	assert.False(t, output.IsCancel())
	assert.Equal(t, failErr, output.Err())

	// Test cancel case
	cancelErr := errors.New("cancel error")
	result = rop.Cancel[int](cancelErr)
	output = solo.CheckCancel(result, boolF, "value not greater than 3")
	assert.False(t, output.IsSuccess())
	assert.True(t, output.IsCancel())
	assert.Equal(t, cancelErr, output.Err())
}

func TestTeeWithErrWithCtx(t *testing.T) {
	ctx := context.Background()
	
	// Test success case with no error from deadEndF
	result := rop.Success(5)
	called := false
	deadEndF := func(ctx context.Context, r rop.Result[int]) error {
		called = true
		assert.Equal(t, 5, r.Result())
		return nil
	}

	output := solo.TeeWithErrWithCtx(ctx, result, deadEndF)
	assert.True(t, called)
	assert.True(t, output.IsSuccess())
	assert.Equal(t, 5, output.Result())

	// Test success case with error from deadEndF
	result = rop.Success(10)
	called = false
	expectedErr := errors.New("deadend error")
	deadEndF = func(ctx context.Context, r rop.Result[int]) error {
		called = true
		return expectedErr
	}

	output = solo.TeeWithErrWithCtx(ctx, result, deadEndF)
	assert.True(t, called)
	assert.False(t, output.IsSuccess())
	assert.Equal(t, expectedErr, output.Err())
}