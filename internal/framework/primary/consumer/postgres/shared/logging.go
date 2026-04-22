package shared

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LogEmailMessage(ctx context.Context, subject string, recipientCount int) zerolog.Logger {
	return log.Ctx(ctx).With().
		Str("subject", subject).
		Int("recipient_count", recipientCount).
		Logger()
}

func CoalesceEmailIdentity(primary, fallback string) string {
	if primary != "" {
		return primary
	}

	return fallback
}
