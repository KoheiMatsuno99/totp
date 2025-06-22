package ctxtimetest

import (
	"context"
	"testing"
	"time"
	"totp/ctxtime/internal"
)

func init() {
	if testing.Testing() {
		internal.Now = nowForTest
	}
}

func nowForTest(ctx context.Context) time.Time {
	now, ok := nowFromContext(ctx)
	if ok {
		return now
	}
	return internal.DefaultNow(ctx)
}

type ctxkey struct{}

func WithFixedNow(t *testing.T, ctx context.Context, tm time.Time) context.Context {
	t.Helper()
	return context.WithValue(ctx, ctxkey{}, tm)
}

func nowFromContext(ctx context.Context) (time.Time, bool) {
	tm, ok := ctx.Value(ctxkey{}).(time.Time)
	return tm, ok
}
