package godoit

import (
	"context"
	"time"
)

type RetryConfig interface {
	RetryTime(ctx context.Context, task Task) time.Time
}
