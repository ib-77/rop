package solo

import (
	"context"
	"errors"
	"github.com/ib-77/rop/pkg/rop"
	"github.com/ib-77/rop/pkg/rop/solo"
	"github.com/ib-77/rop/test"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
	"time"
)

func init() {

}

func TestMain(t *testing.M) {
	setupAll()
	code := t.Run()
	tearDownAll()
	os.Exit(code)
}

func setupAll() {

}

func tearDownAll() {

}

func Test_ValidateTrue(t *testing.T) {
	t.Parallel()

	value := 1
	result := solo.Validate(value, func(a int) bool {
		if a < 2 {
			return true
		}
		return false
	}, "value more than 2")

	assert.True(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, value, result.Result())
	assert.Nil(t, result.Err())
}

func Test_ValidateFalse(t *testing.T) {
	t.Parallel()

	value := 7
	errMsg := "value more than 2"
	result := solo.Validate(value, func(a int) bool {
		if a < 2 {
			return true
		}
		return false
	}, "value more than 2")

	assert.False(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Empty(t, result.Result())
	assert.NotEmpty(t, result.Err())
	assert.Equal(t, errMsg, result.Err().Error())
}

func Test_OnSuccessSwitch_StaySuccess(t *testing.T) {
	t.Parallel()

	value := 100
	okValue := "ok"
	input := rop.Success(value)
	result := solo.Switch(input, func(a int) rop.Result[string] {
		if a == value {
			return rop.Success(okValue)
		}
		return rop.Fail[string](errors.New("fail"))
	})

	assert.True(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, okValue, result.Result())
	assert.Nil(t, result.Err())
}

func Test_OnSuccessSwitch_StayFail(t *testing.T) {
	t.Parallel()

	value := 100
	okValue := "ok"
	input := rop.Success(value)
	result := solo.Switch(input, func(a int) rop.Result[string] {
		if a != value {
			return rop.Success(okValue)
		}
		return rop.Fail[string](errors.New("fail"))
	})

	assert.False(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Empty(t, result.Result())
	assert.NotEmpty(t, result.Err())
	assert.Equal(t, "fail", result.Err().Error())
}

func Test_OnSuccessSwitch_Fail(t *testing.T) {
	t.Parallel()

	value := 100
	okValue := "ok"
	input := rop.Fail[int](errors.New("fail2"))
	result := solo.Switch(input, func(a int) rop.Result[string] {
		if a != value {
			return rop.Success(okValue)
		}
		return rop.Fail[string](errors.New("fail"))
	})

	assert.False(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Empty(t, result.Result())
	assert.NotEmpty(t, result.Err())
	assert.Equal(t, "fail2", result.Err().Error())
}

func Test_OnSuccessMap_StaySuccess(t *testing.T) {
	t.Parallel()

	value := 100
	okValue := "ok"
	input := rop.Success(value)
	result := solo.Map(input, func(a int) string {
		return okValue
	})

	assert.True(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, okValue, result.Result())
	assert.Nil(t, result.Err())
}

func Test_OnSuccessMap_StayFail(t *testing.T) {
	t.Parallel()

	okValue := "ok"
	input := rop.Fail[int](errors.New("fail3"))
	result := solo.Map(input, func(a int) string {
		return okValue
	})

	assert.False(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Empty(t, result.Result())
	assert.NotEmpty(t, result.Err())
	assert.Equal(t, "fail3", result.Err().Error())
}

type Value struct {
	v int
}

func Test_OnSuccessDo_Fail(t *testing.T) {
	t.Parallel()

	input := rop.Fail[Value](errors.New("fail4"))
	result := solo.Tee(input, func(a rop.Result[Value]) {
		// any
	})

	assert.False(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Empty(t, result.Result())
	assert.NotEmpty(t, result.Err())
	assert.Equal(t, "fail4", result.Err().Error())
}

func Test_OnBothMap_Success(t *testing.T) {
	t.Parallel()

	value := 100
	input := rop.Success(value)
	result := solo.DoubleMap(input, func(a int) string {
		return strconv.Itoa(a)
	},
		func(e error) string {
			return e.Error()
		}, func(e error) string {
			return e.Error()
		})

	assert.True(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, strconv.Itoa(value), result.Result())
	assert.Nil(t, result.Err())
}

func Test_Retry_Success(t *testing.T) {
	t.Parallel()

	value := 100
	input := rop.Success(value)
	ctx := rop.WithRetry(context.Background(), rop.NewFixedRetryStrategy(7, time.Second))

	result := solo.RetryWithCtx(ctx, input, throwError)

	assert.True(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Equal(t, "OK", result.Result())
	assert.Nil(t, result.Err())
}

func Test_Retry_Fail(t *testing.T) {
	t.Parallel()

	value := 55
	input := rop.Success(value)
	ctx := rop.WithRetry(context.Background(), rop.NewFixedRetryStrategy(7, time.Second))

	result := solo.RetryWithCtx(ctx, input, throwError)

	assert.False(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	//assert.Equal(t, strconv.Itoa(value), result.Result())
	assert.NotNil(t, result.Err())
}

func throwError(_ context.Context, r int) (string, error) {
	if r == 100 {
		return "OK", nil
	}
	return "ER", errors.New("! 100")
}

func Test_OnBothMap_Fail(t *testing.T) {
	t.Parallel()

	input := rop.Fail[int](errors.New("fails"))
	errMsg := "erroo"
	result := solo.DoubleMap(input, func(a int) string {
		return strconv.Itoa(a)
	}, func(e error) string {
		errMsg += e.Error()
		return e.Error()
	}, func(e error) string {
		errMsg += e.Error()
		return e.Error()
	})

	assert.False(t, result.IsSuccess())
	assert.False(t, result.IsCancel())
	assert.Empty(t, result.Result())
	assert.NotEmpty(t, result.Err())
	assert.Equal(t, "erroo"+"fails", errMsg)
}

func TestCase01(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		result := test.RopCase01(1)
		assert.Equal(t, "all ok", result)
	})

	t.Run("fail at start", func(t *testing.T) {
		t.Parallel()

		result := test.RopCase01(3)
		assert.Equal(t, "error: value more than 2", result)
	})

	t.Run("fail zero", func(t *testing.T) {
		t.Parallel()

		result := test.RopCase01(0)
		assert.Equal(t, "error: a is less or 0!", result)
	})
}

func TestCase02(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		result := test.RopCase02(1)
		assert.Equal(t, "all ok", result)
	})

	t.Run("fail at start", func(t *testing.T) {
		t.Parallel()

		result := test.RopCase02(3)
		assert.Equal(t, "error: value more than 2", result)
	})

	t.Run("fail zero", func(t *testing.T) {
		t.Parallel()

		result := test.RopCase02(0)
		assert.Equal(t, "error: a is less or 0!", result)
	})
}

func TestCase03(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		result := test.RopCase03(ctx, 1)
		assert.Equal(t, "all ok", result)
	})

	t.Run("fail at start", func(t *testing.T) {
		t.Parallel()

		result := test.RopCase03(ctx, 3)
		assert.Equal(t, "error: value more than 2", result)
	})

	t.Run("fail zero", func(t *testing.T) {
		t.Parallel()

		result := test.RopCase03(ctx, 0)
		assert.Equal(t, "error: a is less or 0!", result)
	})
}

func TestCase04(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		result := test.RopCase04(ctx, 1)
		assert.Equal(t, "all ok 100", result)
	})

	t.Run("success default 3", func(t *testing.T) {
		t.Parallel()

		result := test.RopCase04(ctx, 3)
		assert.Equal(t, "all ok 3", result)
	})

	t.Run("fail zero", func(t *testing.T) {
		t.Parallel()

		result := test.RopCase04(ctx, 0)
		assert.Equal(t, "error: a is less or 0!", result)
	})

	t.Run("success default 5", func(t *testing.T) {
		t.Parallel()

		result := test.RopCase04(ctx, 5)
		assert.Equal(t, "all ok 5", result)
	})
}

var benchGlobalRes string

func BenchmarkCase01(b *testing.B) {

	for i := 0; i < b.N; i++ {
		result := test.RopBenchCase01(1)
		benchGlobalRes = result
	}
}
