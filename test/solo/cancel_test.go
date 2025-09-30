package solo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ib-77/rop/pkg/rop"
	"github.com/ib-77/rop/pkg/rop/solo"
	"github.com/stretchr/testify/assert"
)

func TestValidateWithCtxWithCancel_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := "valid_input"
	errMsg := "validation failed"

	// Test case where validation returns true
	validateFn := func(ctx context.Context, in string) bool {
		return in == "valid_input"
	}

	resultCh := solo.ValidateWithCtxWithCancel(ctx, input, validateFn, errMsg)
	result := <-resultCh

	assert.True(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, input, result.Result())
	assert.Nil(t, result.Err())
}

func TestValidateWithCtxWithCancel_ValidationFails(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := "invalid_input"
	errMsg := "validation failed"

	// Test case where validation returns false
	validateFn := func(ctx context.Context, in string) bool {
		return in == "valid_input"
	}

	resultCh := solo.ValidateWithCtxWithCancel(ctx, input, validateFn, errMsg)
	result := <-resultCh

	assert.False(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, "", result.Result()) // zero value for string
	assert.NotNil(t, result.Err())
	assert.Equal(t, errMsg, result.Err().Error())
}

func TestValidateWithCtxWithCancel_ContextCancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	input := "test_input"
	errMsg := "validation failed"

	// Create a validation function that takes some time
	validateFn := func(ctx context.Context, in string) bool {
		time.Sleep(100 * time.Millisecond)
		return true
	}

	resultCh := solo.ValidateWithCtxWithCancel(ctx, input, validateFn, errMsg)

	// Cancel the context before validation completes
	cancel()

	result := <-resultCh

	assert.False(t, result.IsSuccess())
	assert.True(t, result.IsCancel())
	assert.Equal(t, "", result.Result()) // zero value for string
	assert.NotNil(t, result.Err())
	assert.Equal(t, "cancelled", result.Err().Error())
}

func TestValidateWithCtxWithCancel_ContextTimeout(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	input := 42
	errMsg := "number validation failed"

	// Create a validation function that takes longer than the timeout
	validateFn := func(ctx context.Context, in int) bool {
		// Check context periodically during long operation
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				return false // Context cancelled/timed out
			default:
				time.Sleep(10 * time.Millisecond)
			}
		}
		return in > 0
	}

	resultCh := solo.ValidateWithCtxWithCancel(ctx, input, validateFn, errMsg)
	result := <-resultCh

	assert.False(t, result.IsSuccess())
	assert.True(t, result.IsCancel())
	assert.Equal(t, 0, result.Result()) // zero value for int
	assert.NotNil(t, result.Err())
	assert.Equal(t, "cancelled", result.Err().Error())
}

func TestValidateWithCtxWithCancel_QuickValidation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	input := 42
	errMsg := "number validation failed"

	// Create a validation function that completes quickly
	validateFn := func(ctx context.Context, in int) bool {
		return in > 0
	}

	resultCh := solo.ValidateWithCtxWithCancel(ctx, input, validateFn, errMsg)
	result := <-resultCh

	// Should succeed because validation completes before timeout
	assert.True(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, input, result.Result())
	assert.Nil(t, result.Err())
}

func TestValidateWithCtxWithCancel_StructType(t *testing.T) {
	t.Parallel()

	type User struct {
		Name string
		Age  int
	}

	ctx := context.Background()
	input := User{Name: "John", Age: 25}
	errMsg := "user validation failed"

	// Test with custom struct type
	validateFn := func(ctx context.Context, user User) bool {
		return user.Name != "" && user.Age >= 18
	}

	resultCh := solo.ValidateWithCtxWithCancel(ctx, input, validateFn, errMsg)
	result := <-resultCh

	assert.True(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, input, result.Result())
	assert.Nil(t, result.Err())
}

func TestValidateWithCtxWithCancel_StructTypeValidationFails(t *testing.T) {
	t.Parallel()

	type User struct {
		Name string
		Age  int
	}

	ctx := context.Background()
	input := User{Name: "John", Age: 16} // Age below 18
	errMsg := "user must be 18 or older"

	// Test with custom struct type that fails validation
	validateFn := func(ctx context.Context, user User) bool {
		return user.Name != "" && user.Age >= 18
	}

	resultCh := solo.ValidateWithCtxWithCancel(ctx, input, validateFn, errMsg)
	result := <-resultCh

	assert.False(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, User{}, result.Result()) // zero value for User
	assert.NotNil(t, result.Err())
	assert.Equal(t, errMsg, result.Err().Error())
}

func TestValidateWithCtxWithCancel_ValidationUsesContext(t *testing.T) {
	t.Parallel()

	type contextKey string
	const testKey contextKey = "test"

	ctx := context.WithValue(context.Background(), testKey, "expected_value")
	input := "test_input"
	errMsg := "context validation failed"

	// Test that validation function receives and can use the context
	validateFn := func(ctx context.Context, in string) bool {
		value := ctx.Value(testKey)
		return value == "expected_value" && in == "test_input"
	}

	resultCh := solo.ValidateWithCtxWithCancel(ctx, input, validateFn, errMsg)
	result := <-resultCh

	assert.True(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, input, result.Result())
	assert.Nil(t, result.Err())
}

func TestValidateWithCtxWithCancel_ChannelIsClosed(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := "test"
	errMsg := "test error"

	validateFn := func(ctx context.Context, in string) bool {
		return true
	}

	resultCh := solo.ValidateWithCtxWithCancel(ctx, input, validateFn, errMsg)

	// Read the result
	result := <-resultCh
	assert.True(t, result.IsSuccess())

	// Try to read again - should get zero value since channel is closed
	select {
	case result2, ok := <-resultCh:
		assert.False(t, ok, "Channel should be closed")
		assert.True(t, result2.IsSuccess() == false && result2.IsCancel() == false) // zero value
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Should have received from closed channel immediately")
	}
}

func TestValidateWithCtxWithCancel_ConcurrentCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	input := "test_input"
	errMsg := "validation failed"

	// Validation function that checks if context is cancelled during execution
	validateFn := func(ctx context.Context, in string) bool {
		// Simulate some work with periodic context checking
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				return false // Context cancelled during validation
			default:
				time.Sleep(10 * time.Millisecond)
			}
		}
		return true
	}

	resultCh := solo.ValidateWithCtxWithCancel(ctx, input, validateFn, errMsg)

	// Cancel after a short delay
	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel()
	}()

	result := <-resultCh

	// Should be cancelled due to context cancellation
	assert.False(t, result.IsSuccess())
	assert.True(t, result.IsCancel())
	assert.Equal(t, "", result.Result())
	assert.NotNil(t, result.Err())
	assert.Equal(t, "cancelled", result.Err().Error())
}

// Tests for AndValidateWithCtxWithCancel
func TestAndValidateWithCtxWithCancel_SuccessInput_ValidateSuccess(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := rop.Success("valid_input")
	errMsg := "validation failed"

	validateFn := func(ctx context.Context, in string) bool {
		return in == "valid_input"
	}

	resultCh := solo.AndValidateWithCtxWithCancel(ctx, input, validateFn, errMsg)
	result := <-resultCh

	assert.True(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, "valid_input", result.Result())
	assert.Nil(t, result.Err())
}

func TestAndValidateWithCtxWithCancel_SuccessInput_ValidateFails(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := rop.Success("invalid_input")
	errMsg := "validation failed"

	validateFn := func(ctx context.Context, in string) bool {
		return in == "valid_input"
	}

	resultCh := solo.AndValidateWithCtxWithCancel(ctx, input, validateFn, errMsg)
	result := <-resultCh

	assert.False(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, "", result.Result())
	assert.NotNil(t, result.Err())
	assert.Equal(t, errMsg, result.Err().Error())
}

func TestAndValidateWithCtxWithCancel_FailInput(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := rop.Fail[string](errors.New("initial error"))
	errMsg := "validation failed"

	validateFn := func(ctx context.Context, in string) bool {
		return true
	}

	resultCh := solo.AndValidateWithCtxWithCancel(ctx, input, validateFn, errMsg)
	result := <-resultCh

	assert.False(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, "", result.Result())
	assert.NotNil(t, result.Err())
	assert.Equal(t, "initial error", result.Err().Error())
}

func TestAndValidateWithCtxWithCancel_ContextCancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	input := rop.Success("test_input")
	errMsg := "validation failed"

	validateFn := func(ctx context.Context, in string) bool {
		time.Sleep(100 * time.Millisecond)
		return true
	}

	resultCh := solo.AndValidateWithCtxWithCancel(ctx, input, validateFn, errMsg)
	cancel()

	result := <-resultCh

	assert.False(t, result.IsSuccess())
	assert.True(t, result.IsCancel())
	assert.Equal(t, "", result.Result())
	assert.NotNil(t, result.Err())
	assert.Equal(t, "cancelled", result.Err().Error())
}

// Tests for SwitchWithCtxWithCancel
func TestSwitchWithCtxWithCancel_SuccessInput_SwitchSuccess(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := rop.Success(42)

	switchFn := func(ctx context.Context, in int) rop.Result[string] {
		return rop.Success("converted: " + string(rune(in)))
	}

	resultCh := solo.SwitchWithCtxWithCancel(ctx, input, switchFn)
	result := <-resultCh

	assert.True(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, "converted: *", result.Result())
	assert.Nil(t, result.Err())
}

func TestSwitchWithCtxWithCancel_SuccessInput_SwitchFails(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := rop.Success(42)

	switchFn := func(ctx context.Context, in int) rop.Result[string] {
		return rop.Fail[string](errors.New("switch failed"))
	}

	resultCh := solo.SwitchWithCtxWithCancel(ctx, input, switchFn)
	result := <-resultCh

	assert.False(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, "", result.Result())
	assert.NotNil(t, result.Err())
	assert.Equal(t, "switch failed", result.Err().Error())
}

func TestSwitchWithCtxWithCancel_FailInput(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := rop.Fail[int](errors.New("initial error"))

	switchFn := func(ctx context.Context, in int) rop.Result[string] {
		return rop.Success("should not be called")
	}

	resultCh := solo.SwitchWithCtxWithCancel(ctx, input, switchFn)
	result := <-resultCh

	assert.False(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, "", result.Result())
	assert.NotNil(t, result.Err())
	assert.Equal(t, "initial error", result.Err().Error())
}

func TestSwitchWithCtxWithCancel_ContextCancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	input := rop.Success(42)

	switchFn := func(ctx context.Context, in int) rop.Result[string] {
		time.Sleep(100 * time.Millisecond)
		return rop.Success("should not complete")
	}

	resultCh := solo.SwitchWithCtxWithCancel(ctx, input, switchFn)
	cancel()

	result := <-resultCh

	assert.False(t, result.IsSuccess())
	assert.True(t, result.IsCancel())
	assert.Equal(t, "", result.Result())
	assert.NotNil(t, result.Err())
	assert.Equal(t, "cancelled", result.Err().Error())
}

// Tests for MapWithCtxWithCancel
func TestMapWithCtxWithCancel_SuccessInput(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := rop.Success(42)

	mapFn := func(ctx context.Context, in int) string {
		return "mapped: " + string(rune(in+48)) // Convert to ASCII digit
	}

	resultCh := solo.MapWithCtxWithCancel(ctx, input, mapFn)
	result := <-resultCh

	assert.True(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, "mapped: Z", result.Result())
	assert.Nil(t, result.Err())
}

func TestMapWithCtxWithCancel_FailInput(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := rop.Fail[int](errors.New("initial error"))

	mapFn := func(ctx context.Context, in int) string {
		return "should not be called"
	}

	resultCh := solo.MapWithCtxWithCancel(ctx, input, mapFn)
	result := <-resultCh

	assert.False(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, "", result.Result())
	assert.NotNil(t, result.Err())
	assert.Equal(t, "initial error", result.Err().Error())
}

func TestMapWithCtxWithCancel_CancelInput(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := rop.Cancel[int](errors.New("cancelled error"))

	mapFn := func(ctx context.Context, in int) string {
		return "should not be called"
	}

	resultCh := solo.MapWithCtxWithCancel(ctx, input, mapFn)
	result := <-resultCh

	assert.False(t, result.IsSuccess())
	assert.True(t, result.IsCancel())
	assert.Equal(t, "", result.Result())
	assert.NotNil(t, result.Err())
	assert.Equal(t, "cancelled error", result.Err().Error())
}

func TestMapWithCtxWithCancel_ContextCancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	input := rop.Success(42)

	mapFn := func(ctx context.Context, in int) string {
		time.Sleep(100 * time.Millisecond)
		return "should not complete"
	}

	resultCh := solo.MapWithCtxWithCancel(ctx, input, mapFn)
	cancel()

	result := <-resultCh

	assert.False(t, result.IsSuccess())
	assert.True(t, result.IsCancel())
	assert.Equal(t, "", result.Result())
	assert.NotNil(t, result.Err())
	assert.Equal(t, "cancelled", result.Err().Error())
}

// Tests for TeeWithCtxWithCancel
func TestTeeWithCtxWithCancel_SuccessInput(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := rop.Success(42)
	var sideEffectValue int

	deadEndFn := func(ctx context.Context, r rop.Result[int]) {
		if r.IsSuccess() {
			sideEffectValue = r.Result() * 2
		}
	}

	resultCh := solo.TeeWithCtxWithCancel(ctx, input, deadEndFn)
	result := <-resultCh

	assert.True(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, 42, result.Result())
	assert.Nil(t, result.Err())
	assert.Equal(t, 84, sideEffectValue) // Side effect should have occurred
}

func TestTeeWithCtxWithCancel_FailInput(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := rop.Fail[int](errors.New("initial error"))
	var sideEffectCalled bool

	deadEndFn := func(ctx context.Context, r rop.Result[int]) {
		sideEffectCalled = true
	}

	resultCh := solo.TeeWithCtxWithCancel(ctx, input, deadEndFn)
	result := <-resultCh

	assert.False(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, 0, result.Result())
	assert.NotNil(t, result.Err())
	assert.Equal(t, "initial error", result.Err().Error())
	assert.False(t, sideEffectCalled) // Side effect should not be called for failed input
}

// Tests for DoubleTeeWithCtxWithCancel
func TestDoubleTeeWithCtxWithCancel_SuccessInput(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := rop.Success(42)
	var successSideEffect int
	var errorSideEffect bool

	successFn := func(ctx context.Context, r int) {
		successSideEffect = r * 2
	}

	errorFn := func(ctx context.Context, err error) {
		errorSideEffect = true
	}

	resultCh := solo.DoubleTeeWithCtxWithCancel(ctx, input, successFn, errorFn)
	result := <-resultCh

	assert.True(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, 42, result.Result())
	assert.Nil(t, result.Err())
	assert.Equal(t, 84, successSideEffect)
	assert.False(t, errorSideEffect)
}

func TestDoubleTeeWithCtxWithCancel_FailInput(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := rop.Fail[int](errors.New("initial error"))
	var successSideEffect bool
	var errorMessage string

	successFn := func(ctx context.Context, r int) {
		successSideEffect = true
	}

	errorFn := func(ctx context.Context, err error) {
		errorMessage = err.Error()
	}

	resultCh := solo.DoubleTeeWithCtxWithCancel(ctx, input, successFn, errorFn)
	result := <-resultCh

	// Small delay to ensure side effect completes
	time.Sleep(10 * time.Millisecond)

	assert.False(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, 0, result.Result())
	assert.NotNil(t, result.Err())
	assert.Equal(t, "initial error", result.Err().Error())
	assert.False(t, successSideEffect)
	assert.Equal(t, "initial error", errorMessage)
}

// Tests for DoubleMapWithCtxWithCancel
func TestDoubleMapWithCtxWithCancel_SuccessInput(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := rop.Success(42)

	successFn := func(ctx context.Context, r int) string {
		return "success: " + string(rune(r+48))
	}

	failFn := func(ctx context.Context, err error) string {
		return "fail: " + err.Error()
	}

	cancelFn := func(ctx context.Context, err error) string {
		return "cancel: " + err.Error()
	}

	resultCh := solo.DoubleMapWithCtxWithCancel(ctx, input, successFn, failFn, cancelFn)
	result := <-resultCh

	assert.True(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, "success: Z", result.Result())
	assert.Nil(t, result.Err())
}

func TestDoubleMapWithCtxWithCancel_FailInput(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := rop.Fail[int](errors.New("initial error"))
	var failCalled bool

	successFn := func(ctx context.Context, r int) string {
		return "should not be called"
	}

	failFn := func(ctx context.Context, err error) string {
		failCalled = true
		return "fail handled"
	}

	cancelFn := func(ctx context.Context, err error) string {
		return "cancel handled"
	}

	resultCh := solo.DoubleMapWithCtxWithCancel(ctx, input, successFn, failFn, cancelFn)
	result := <-resultCh

	assert.False(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, "", result.Result())
	assert.NotNil(t, result.Err())
	assert.Equal(t, "initial error", result.Err().Error())
	assert.True(t, failCalled)
}

func TestDoubleMapWithCtxWithCancel_CancelInput(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := rop.Cancel[int](errors.New("cancelled error"))
	var cancelCalled bool

	successFn := func(ctx context.Context, r int) string {
		return "should not be called"
	}

	failFn := func(ctx context.Context, err error) string {
		return "fail handled"
	}

	cancelFn := func(ctx context.Context, err error) string {
		cancelCalled = true
		return "cancel handled"
	}

	resultCh := solo.DoubleMapWithCtxWithCancel(ctx, input, successFn, failFn, cancelFn)
	result := <-resultCh

	assert.False(t, result.IsSuccess())
	assert.True(t, result.IsCancel())
	assert.Equal(t, "", result.Result())
	assert.NotNil(t, result.Err())
	assert.Equal(t, "cancelled error", result.Err().Error())
	assert.True(t, cancelCalled)
}

func TestDoubleMapWithCtxWithCancel_ContextCancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	input := rop.Success(42)

	successFn := func(ctx context.Context, r int) string {
		time.Sleep(100 * time.Millisecond)
		return "should not complete"
	}

	failFn := func(ctx context.Context, err error) string {
		return "fail handled"
	}

	cancelFn := func(ctx context.Context, err error) string {
		return "cancel handled"
	}

	resultCh := solo.DoubleMapWithCtxWithCancel(ctx, input, successFn, failFn, cancelFn)
	cancel()

	result := <-resultCh

	assert.False(t, result.IsSuccess())
	assert.True(t, result.IsCancel())
	assert.Equal(t, "", result.Result())
	assert.NotNil(t, result.Err())
	assert.Equal(t, "cancelled", result.Err().Error())
}

// Tests for TryWithCtxWithCancel
func TestTryWithCtxWithCancel_SuccessInput_TrySuccess(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := rop.Success(42)

	tryFn := func(ctx context.Context, r int) (string, error) {
		if r > 0 {
			return "positive: " + string(rune(r+48)), nil
		}
		return "", errors.New("negative number")
	}

	resultCh := solo.TryWithCtxWithCancel(ctx, input, tryFn)
	result := <-resultCh

	assert.True(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, "positive: Z", result.Result())
	assert.Nil(t, result.Err())
}

func TestTryWithCtxWithCancel_SuccessInput_TryFails(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := rop.Success(-5)

	tryFn := func(ctx context.Context, r int) (string, error) {
		if r > 0 {
			return "positive", nil
		}
		return "", errors.New("negative number")
	}

	resultCh := solo.TryWithCtxWithCancel(ctx, input, tryFn)
	result := <-resultCh

	assert.False(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, "", result.Result())
	assert.NotNil(t, result.Err())
	assert.Equal(t, "negative number", result.Err().Error())
}

func TestTryWithCtxWithCancel_FailInput(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := rop.Fail[int](errors.New("initial error"))

	tryFn := func(ctx context.Context, r int) (string, error) {
		return "should not be called", nil
	}

	resultCh := solo.TryWithCtxWithCancel(ctx, input, tryFn)
	result := <-resultCh

	assert.False(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, "", result.Result())
	assert.NotNil(t, result.Err())
	assert.Equal(t, "initial error", result.Err().Error())
}

func TestTryWithCtxWithCancel_CancelInput(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := rop.Cancel[int](errors.New("cancelled error"))

	tryFn := func(ctx context.Context, r int) (string, error) {
		return "should not be called", nil
	}

	resultCh := solo.TryWithCtxWithCancel(ctx, input, tryFn)
	result := <-resultCh

	assert.False(t, result.IsSuccess())
	assert.True(t, result.IsCancel())
	assert.Equal(t, "", result.Result())
	assert.NotNil(t, result.Err())
	assert.Equal(t, "cancelled error", result.Err().Error())
}

func TestTryWithCtxWithCancel_ContextCancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	input := rop.Success(42)

	tryFn := func(ctx context.Context, r int) (string, error) {
		time.Sleep(100 * time.Millisecond)
		return "should not complete", nil
	}

	resultCh := solo.TryWithCtxWithCancel(ctx, input, tryFn)
	cancel()

	result := <-resultCh

	assert.False(t, result.IsSuccess())
	assert.True(t, result.IsCancel())
	assert.Equal(t, "", result.Result())
	assert.NotNil(t, result.Err())
	assert.Equal(t, "cancelled", result.Err().Error())
}

// Tests for FinallyWithCtxWithCancel
func TestFinallyWithCtxWithCancel_SuccessInput(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	inputCh := make(chan rop.Result[int], 1)
	inputCh <- rop.Success(42)
	close(inputCh)

	successFn := func(ctx context.Context, r int) string {
		return "success: " + string(rune(r+48)) // Convert to ASCII
	}

	failOrCancelFn := func(ctx context.Context, err error) string {
		return "error: " + err.Error()
	}

	result := solo.FinallyWithCtxWithCancel(ctx, inputCh, successFn, failOrCancelFn)

	assert.Equal(t, "success: Z", result)
}

func TestFinallyWithCtxWithCancel_FailInput(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	inputCh := make(chan rop.Result[int], 1)
	inputCh <- rop.Fail[int](errors.New("processing failed"))
	close(inputCh)

	successFn := func(ctx context.Context, r int) string {
		return "should not be called"
	}

	failOrCancelFn := func(ctx context.Context, err error) string {
		return "error: " + err.Error()
	}

	result := solo.FinallyWithCtxWithCancel(ctx, inputCh, successFn, failOrCancelFn)

	assert.Equal(t, "error: processing failed", result)
}

func TestFinallyWithCtxWithCancel_CancelInput(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	inputCh := make(chan rop.Result[int], 1)
	inputCh <- rop.Cancel[int](errors.New("operation cancelled"))
	close(inputCh)

	successFn := func(ctx context.Context, r int) string {
		return "should not be called"
	}

	failOrCancelFn := func(ctx context.Context, err error) string {
		return "cancelled: " + err.Error()
	}

	result := solo.FinallyWithCtxWithCancel(ctx, inputCh, successFn, failOrCancelFn)

	assert.Equal(t, "cancelled: operation cancelled", result)
}

func TestFinallyWithCtxWithCancel_ContextCancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	inputCh := make(chan rop.Result[int])

	// Cancel context immediately
	cancel()

	successFn := func(ctx context.Context, r int) string {
		return "should not be called"
	}

	failOrCancelFn := func(ctx context.Context, err error) string {
		return "context cancelled: " + err.Error()
	}

	result := solo.FinallyWithCtxWithCancel(ctx, inputCh, successFn, failOrCancelFn)

	assert.Equal(t, "context cancelled: cancelled", result)
}

func TestFinallyWithCtxWithCancel_ContextTimeout(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	inputCh := make(chan rop.Result[int])

	// Don't send anything to the channel, let it timeout
	go func() {
		time.Sleep(100 * time.Millisecond)
		inputCh <- rop.Success(42)
		close(inputCh)
	}()

	successFn := func(ctx context.Context, r int) string {
		return "should not be called"
	}

	failOrCancelFn := func(ctx context.Context, err error) string {
		return "timeout: " + err.Error()
	}

	result := solo.FinallyWithCtxWithCancel(ctx, inputCh, successFn, failOrCancelFn)

	assert.Equal(t, "timeout: cancelled", result)
}

func TestFinallyWithCtxWithCancel_ContextUsage(t *testing.T) {
	t.Parallel()

	type contextKey string
	const testKey contextKey = "test"

	ctx := context.WithValue(context.Background(), testKey, "expected_value")
	inputCh := make(chan rop.Result[string], 1)
	inputCh <- rop.Success("input_value")
	close(inputCh)

	successFn := func(ctx context.Context, r string) string {
		contextValue := ctx.Value(testKey)
		return "success with context: " + contextValue.(string) + " and input: " + r
	}

	failOrCancelFn := func(ctx context.Context, err error) string {
		contextValue := ctx.Value(testKey)
		return "error with context: " + contextValue.(string) + " - " + err.Error()
	}

	result := solo.FinallyWithCtxWithCancel(ctx, inputCh, successFn, failOrCancelFn)

	assert.Equal(t, "success with context: expected_value and input: input_value", result)
}

func TestFinallyWithCtxWithCancel_StructType(t *testing.T) {
	t.Parallel()

	type User struct {
		Name string
		Age  int
	}

	type UserReport struct {
		Status string
		User   User
		Error  string
	}

	ctx := context.Background()
	inputCh := make(chan rop.Result[User], 1)
	inputCh <- rop.Success(User{Name: "Alice", Age: 30})
	close(inputCh)

	successFn := func(ctx context.Context, user User) UserReport {
		return UserReport{
			Status: "processed",
			User:   user,
			Error:  "",
		}
	}

	failOrCancelFn := func(ctx context.Context, err error) UserReport {
		return UserReport{
			Status: "failed",
			User:   User{},
			Error:  err.Error(),
		}
	}

	result := solo.FinallyWithCtxWithCancel(ctx, inputCh, successFn, failOrCancelFn)

	expected := UserReport{
		Status: "processed",
		User:   User{Name: "Alice", Age: 30},
		Error:  "",
	}
	assert.Equal(t, expected, result)
}

func TestFinallyWithCtxWithCancel_StructTypeWithError(t *testing.T) {
	t.Parallel()

	type User struct {
		Name string
		Age  int
	}

	type UserReport struct {
		Status string
		User   User
		Error  string
	}

	ctx := context.Background()
	inputCh := make(chan rop.Result[User], 1)
	inputCh <- rop.Fail[User](errors.New("user validation failed"))
	close(inputCh)

	successFn := func(ctx context.Context, user User) UserReport {
		return UserReport{
			Status: "processed",
			User:   user,
			Error:  "",
		}
	}

	failOrCancelFn := func(ctx context.Context, err error) UserReport {
		return UserReport{
			Status: "failed",
			User:   User{},
			Error:  err.Error(),
		}
	}

	result := solo.FinallyWithCtxWithCancel(ctx, inputCh, successFn, failOrCancelFn)

	expected := UserReport{
		Status: "failed",
		User:   User{},
		Error:  "user validation failed",
	}
	assert.Equal(t, expected, result)
}

func TestFinallyWithCtxWithCancel_ConcurrentChannelOperation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	inputCh := make(chan rop.Result[int])

	// Send result asynchronously
	go func() {
		time.Sleep(10 * time.Millisecond)
		inputCh <- rop.Success(99)
		close(inputCh)
	}()

	successFn := func(ctx context.Context, r int) string {
		return "async success: " + string(rune(r+48))
	}

	failOrCancelFn := func(ctx context.Context, err error) string {
		return "async error: " + err.Error()
	}

	result := solo.FinallyWithCtxWithCancel(ctx, inputCh, successFn, failOrCancelFn)

	// 99 + 48 = 147, which corresponds to a specific character
	assert.Contains(t, result, "async success:")
}
