package core

import (
	"context"

	"codebase-app/internal/infrastructure/tracing"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (c *storageCore) ListFiles(ctx context.Context) ([]types.Object, error) {
	ctx, span := tracing.StartSpan(ctx, "core.ListFiles")
	defer span.End()

	return c.s3.ListFiles(ctx)
}
