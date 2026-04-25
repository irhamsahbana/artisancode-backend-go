package exportjob

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/coreentity"
	infraTracing "codebase-app/internal/infrastructure/tracing"
	corePorts "codebase-app/internal/ports/core"
	integrationPorts "codebase-app/internal/ports/integration"

	"github.com/rs/zerolog/log"
)

func RequestedHandler(core corePorts.ExportJobCore) func(context.Context, integrationPorts.MessageBusMessage) {
	return func(ctx context.Context, msg integrationPorts.MessageBusMessage) {
		ctx, span := infraTracing.StartSpan(ctx, "internal:framework:primary:consumer:postgres:export_job:handler:RequestedHandler")
		defer span.End()

		payload := &coreentity.ExportJobRequestedEvent{}

		err := json.Unmarshal(msg.Data(), payload)
		if err != nil {
			infraTracing.RecordError(span, err)
			log.Ctx(ctx).Error().Err(err).Str("subject", msg.Subject()).Msg("invalid export job payload")
			if ackErr := msg.Ack(ctx); ackErr != nil {
				infraTracing.RecordError(span, ackErr)
				log.Ctx(ctx).Error().Err(ackErr).Str("subject", msg.Subject()).Msg("failed to acknowledge invalid message")
			}
			return
		}

		log.Ctx(ctx).Debug().
			Str("subject", msg.Subject()).
			Str("job_id", payload.JobID).
			Str("tenant_id", payload.TenantID).
			Msg("received message")

		err = core.ProcessExportJob(ctx, coreentity.ExportJobDetailFilter{
			TenantID: payload.TenantID,
			ID:       payload.JobID,
		})
		if err != nil {
			infraTracing.RecordError(span, err)
			log.Ctx(ctx).Error().
				Err(err).
				Str("subject", msg.Subject()).
				Str("job_id", payload.JobID).
				Str("tenant_id", payload.TenantID).
				Msg("failed to process export job")
			if nakErr := msg.Nak(ctx, err.Error()); nakErr != nil {
				infraTracing.RecordError(span, nakErr)
				log.Ctx(ctx).Error().
					Err(nakErr).
					Str("subject", msg.Subject()).
					Str("job_id", payload.JobID).
					Str("tenant_id", payload.TenantID).
					Msg("failed to negative-acknowledge message")
			}
			return
		}

		if err := msg.Ack(ctx); err != nil {
			infraTracing.RecordError(span, err)
			log.Ctx(ctx).Error().
				Err(err).
				Str("subject", msg.Subject()).
				Str("job_id", payload.JobID).
				Str("tenant_id", payload.TenantID).
				Msg("failed to acknowledge message")
			return
		}

		log.Ctx(ctx).Info().
			Str("subject", msg.Subject()).
			Str("job_id", payload.JobID).
			Str("tenant_id", payload.TenantID).
			Msg("finished processing message")
	}
}
