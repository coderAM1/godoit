package godoit

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type IdMaker func(ctx context.Context, taskName string, when time.Time) (string, error)

func DefaultIdMaker(ctx context.Context, taskName string, when time.Time) (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func ExampleIdMaker(ctx context.Context, taskName string, when time.Time) (string, error) {
	id := fmt.Sprintf("%s~%s", taskName, when)
	return id, nil
}
