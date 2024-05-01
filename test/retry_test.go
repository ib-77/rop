package test

import (
	"context"
	"github.com/ib-77/rop/pkg/rop"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_FixedRetryStrategy(t *testing.T) {
	t.Parallel()
	du := 2 * time.Second
	rs := rop.NewFixedRetryStrategy(10, du)
	assert.Equal(t, int64(10), rs.Attempts())
	assert.Equal(t, du, rs.Wait(1))
	assert.Equal(t, du, rs.Wait(5))
}

func Test_LinearRetryStrategy(t *testing.T) {
	t.Parallel()
	du := 2 * time.Second
	maxDu := 6 * time.Second
	rs := rop.NewLinearRetryStrategy(10, du, &maxDu)
	assert.Equal(t, int64(10), rs.Attempts())
	assert.Equal(t, du, rs.Wait(1))
	assert.Equal(t, du*2, rs.Wait(2))
	assert.Equal(t, du*3, rs.Wait(3))
	assert.Equal(t, maxDu, rs.Wait(4))
	assert.Equal(t, maxDu, rs.Wait(5))
}

func Test_LinearRetryStrategyMaxAttempts(t *testing.T) {
	t.Parallel()
	du := 2 * time.Second
	maxDu := 6 * time.Second
	rs := rop.NewLinearRetryStrategy(2, du, &maxDu)
	assert.Equal(t, int64(2), rs.Attempts())
	assert.Equal(t, du, rs.Wait(1))
	assert.Equal(t, du*2, rs.Wait(2))
	assert.Equal(t, du*2, rs.Wait(3))
	assert.Equal(t, du*2, rs.Wait(4))
	assert.Equal(t, du*2, rs.Wait(5))
}

func Test_LinearRetryStrategyMaxAttemptsWithoutMaxDuration(t *testing.T) {
	t.Parallel()
	du := 2 * time.Second
	rs := rop.NewLinearRetryStrategy(2, du, nil)
	assert.Equal(t, int64(2), rs.Attempts())
	assert.Equal(t, du, rs.Wait(1))
	assert.Equal(t, du*2, rs.Wait(2))
	assert.Equal(t, du*2, rs.Wait(3))
	assert.Equal(t, du*2, rs.Wait(4))
	assert.Equal(t, du*2, rs.Wait(5))
}

func Test_ExpRetryStrategy(t *testing.T) {
	t.Parallel()
	du := 2 * time.Second
	maxDu := 8 * time.Second
	rs := rop.NewExponentialRetryStrategy(10, 2.0, du, &maxDu)
	assert.Equal(t, int64(10), rs.Attempts())
	assert.Equal(t, du*2, rs.Wait(1))
	assert.Equal(t, du*4, rs.Wait(2))
	assert.Equal(t, maxDu, rs.Wait(3))
	assert.Equal(t, maxDu, rs.Wait(4))
}

func Test_ExpRetryStrategyMaxAttempts(t *testing.T) {
	t.Parallel()
	du := 2 * time.Second
	maxDu := 8 * time.Second
	rs := rop.NewExponentialRetryStrategy(2, 2.0, du, &maxDu)
	assert.Equal(t, int64(2), rs.Attempts())
	assert.Equal(t, du*2, rs.Wait(1))
	assert.Equal(t, du*4, rs.Wait(2))
	assert.Equal(t, du*4, rs.Wait(3))
	assert.Equal(t, du*4, rs.Wait(4))
}

func Test_ExpRetryStrategyMaxAttemptsWitoutMaxDuration(t *testing.T) {
	t.Parallel()
	du := 2 * time.Second
	rs := rop.NewExponentialRetryStrategy(2, 2.0, du, nil)
	assert.Equal(t, int64(2), rs.Attempts())
	assert.Equal(t, du*2, rs.Wait(1))
	assert.Equal(t, du*4, rs.Wait(2))
	assert.Equal(t, du*4, rs.Wait(3))
	assert.Equal(t, du*4, rs.Wait(4))
}

func Test_WithRetryStrategyContext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	du := 2 * time.Second
	maxDu := 6 * time.Second
	ctx = rop.WithRetry(ctx, rop.NewLinearRetryStrategy(2, du, &maxDu))

	rs := rop.GetRetryFromCtx(ctx)
	assert.Equal(t, int64(2), rs.Attempts())
	assert.Equal(t, du, rs.Wait(1))
	assert.Equal(t, du*2, rs.Wait(2))
	assert.Equal(t, du*2, rs.Wait(3))
}
