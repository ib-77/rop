package test

import (
	"errors"
	"github.com/ib-77/rop/pkg/rop"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
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
	result := rop.Validate(value, func(a int) bool {
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
	result := rop.Validate(value, func(a int) bool {
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
	result := rop.Switch(input, func(a int) rop.Rop[string] {
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
	result := rop.Switch(input, func(a int) rop.Rop[string] {
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
	result := rop.Switch(input, func(a int) rop.Rop[string] {
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
	result := rop.Map(input, func(a int) string {
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
	result := rop.Map(input, func(a int) string {
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
	result := rop.Tee(input, func(a rop.Rop[Value]) {
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
	result := rop.DoubleMap(input, func(a int) string {
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

func Test_OnBothMap_Fail(t *testing.T) {
	t.Parallel()

	input := rop.Fail[int](errors.New("fails"))
	errMsg := "erroo"
	result := rop.DoubleMap(input, func(a int) string {
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
	assert.Equal(t, errMsg, "erroo"+"fails")
}

func TestCase01(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		result := RopCase01(1)
		assert.Equal(t, result, "all ok")
	})

	t.Run("fail at start", func(t *testing.T) {
		t.Parallel()

		result := RopCase01(3)
		assert.Equal(t, result, "error: value more than 2")
	})

	t.Run("fail zero", func(t *testing.T) {
		t.Parallel()

		result := RopCase01(0)
		assert.Equal(t, result, "error: a is less or 0!")
	})
}
