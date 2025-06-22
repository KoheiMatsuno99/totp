package internal

import (
	"context"
	"time"
)

var Now = DefaultNow

func DefaultNow(_ context.Context) time.Time {
	return time.Now()
}
