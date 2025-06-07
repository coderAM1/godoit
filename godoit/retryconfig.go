package godoit

import (
	"context"
	"time"
)

type RetryConfig func(ctx context.Context, task Task) time.Time
