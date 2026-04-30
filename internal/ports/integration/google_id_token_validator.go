package integration

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type GoogleIDTokenValidator interface {
	Validate(ctx context.Context, idToken string, nonce string) (*coreentity.GoogleIdentity, error)
}
