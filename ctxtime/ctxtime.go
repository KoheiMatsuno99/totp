package ctxtime

import (
	"context"
	"time"
	"totp/ctxtime/internal"
)

func Now(ctx context.Context) time.Time {
	return internal.Now(ctx)
}
