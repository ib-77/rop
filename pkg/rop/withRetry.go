package rop

import (
	"context"
	"math"
	"time"
)

const (
	ExponentialFactor = 2.0
	RetryStrategyKey  = "retry-strategy"
)

type RetryStrategy interface {
	Attempts() int64
	Wait(attempt int64) time.Duration
}

func WithRetry(ctx context.Context, strategy RetryStrategy) context.Context {
	return context.WithValue(ctx, RetryStrategyKey, strategy)
}

func GetRetryFromCtx(ctx context.Context) (rs RetryStrategy, ok bool) {
	rs, ok = ctx.Value(RetryStrategyKey).(RetryStrategy)
	return
}

func GetRetryFromCtxDef(ctx context.Context, def RetryStrategy) RetryStrategy {
	rs, ok := ctx.Value(RetryStrategyKey).(RetryStrategy)
	if ok {
		return rs
	}
	return def
}

type FixedRetryStrategy struct {
	attempts int64
	delay    time.Duration
}

func NewFixedRetryStrategy(attempts int64, delay time.Duration) FixedRetryStrategy {
	return FixedRetryStrategy{
		attempts: attempts,
		delay:    delay,
	}
}

func (r FixedRetryStrategy) Attempts() int64 {
	return r.attempts
}

func (r FixedRetryStrategy) Wait(int64) time.Duration {
	return r.delay
}

type LinearRetryStrategy struct {
	ctx      FixedRetryStrategy
	maxDelay *time.Duration
}

func NewLinearRetryStrategy(attempts int64, delay time.Duration, maxDelay *time.Duration) LinearRetryStrategy {
	return LinearRetryStrategy{
		ctx: FixedRetryStrategy{
			attempts: attempts,
			delay:    delay,
		},
		maxDelay: maxDelay,
	}
}

func (r LinearRetryStrategy) Attempts() int64 {
	return r.ctx.attempts
}

func (r LinearRetryStrategy) Wait(attempt int64) time.Duration {

	if attempt > r.ctx.attempts {
		attempt = r.ctx.attempts
	}

	currentDelay := time.Duration(r.ctx.delay.Nanoseconds() * attempt)
	//TODO check overflow!
	if r.maxDelay != nil && currentDelay.Nanoseconds() >= r.maxDelay.Nanoseconds() {
		return *r.maxDelay
	} else {
		return currentDelay
	}
}

type ExponentialRetryStrategy struct {
	ctx      FixedRetryStrategy
	maxDelay *time.Duration
	factor   float64
}

func NewExponentialRetryStrategy(attempts int64, factor float64,
	delay time.Duration, maxDelay *time.Duration) ExponentialRetryStrategy {
	return ExponentialRetryStrategy{
		ctx: FixedRetryStrategy{
			attempts: attempts,
			delay:    delay,
		},
		maxDelay: maxDelay,
		factor:   factor,
	}
}

func (r ExponentialRetryStrategy) Attempts() int64 {
	return r.ctx.attempts
}

func (r ExponentialRetryStrategy) Wait(attempt int64) time.Duration {

	if attempt > r.ctx.attempts {
		attempt = r.ctx.attempts
	}

	currentDelay := time.Duration(r.ctx.delay.Nanoseconds() * int64(math.Pow(r.factor, float64(attempt))))
	//TODO check overflow!
	if r.maxDelay != nil && currentDelay.Nanoseconds() >= r.maxDelay.Nanoseconds() {
		return *r.maxDelay
	} else {
		return currentDelay
	}
}
