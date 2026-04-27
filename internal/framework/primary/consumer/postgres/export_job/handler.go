package exportjob

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/coreentity"
	sharedconsumer "codebase-app/internal/framework/primary/consumer/postgres/shared"
	infraTracing "codebase-app/internal/infrastructure/tracing"

	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
)

func (a *Adapter) RequestedHandler(msg *watermillMessage.Message) error {
	ctx := otel.GetTextMapPropagator().Extract(
		msg.Context(),
		sharedconsumer.WatermillMetadataCarrier(msg.Metadata),
	)
	ctx, span := infraTracing.StartSpan(ctx, "internal:framework:primary:consumer:postgres:export_job:handler:RequestedHandler")
	defer span.End()

	payload := &coreentity.ExportJobRequestedEvent{}

	err := json.Unmarshal(msg.Payload, payload)
	if err != nil {
		infraTracing.RecordError(span, err)
		log.Ctx(ctx).Error().Err(err).Str("subject", msg.Metadata.Get("topic")).Msg("invalid export job payload")
		msg.Ack()
		return nil
	}

	return a.processRequested(ctx, msg, payload)
}

func (a *Adapter) processRequested(
	ctx context.Context,
	msg *watermillMessage.Message,
	payload *coreentity.ExportJobRequestedEvent,
) error {
	log.Ctx(ctx).Debug().
		Str("subject", msg.Metadata.Get("topic")).
		Str("job_id", payload.JobID).
		Str("tenant_id", payload.TenantID).
		Msg("received message")

	err := a.core.ProcessExportJob(ctx, coreentity.ExportJobDetailFilter{
		TenantID: payload.TenantID,
		ID:       payload.JobID,
	})
	if err != nil {
		log.Ctx(ctx).Error().
			Err(err).
			Str("subject", msg.Metadata.Get("topic")).
			Str("job_id", payload.JobID).
			Str("tenant_id", payload.TenantID).
			Msg("failed to process export job")
		return err
	}

	log.Ctx(ctx).Info().
		Str("subject", msg.Metadata.Get("topic")).
		Str("job_id", payload.JobID).
		Str("tenant_id", payload.TenantID).
		Msg("finished processing message")

	return nil
}
