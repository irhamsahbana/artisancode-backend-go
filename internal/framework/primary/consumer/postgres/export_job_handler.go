package consumer

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/coreentity"
	infraTracing "codebase-app/internal/infrastructure/tracing"
	integrationPorts "codebase-app/internal/ports/secondary/integration"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
)

func ExportJobRequestedHandler(ctx context.Context, core IntegrationExportJobProcessor) func(integrationPorts.MessageBusMessage) {
	return func(msg integrationPorts.MessageBusMessage) {
		msgCtx := otel.GetTextMapPropagator().Extract(ctx, MessageHeadersCarrier(msg.Headers()))
		msgCtx, span := infraTracing.StartSpan(msgCtx, "consumer.ExportJobRequestedHandler")
		defer span.End()

		payload := &coreentity.ExportJobRequestedEvent{}

		err := json.Unmarshal(msg.Data(), payload)
		if err != nil {
			infraTracing.RecordError(span, err)
			log.Ctx(msgCtx).Error().Err(err).Msg("consumer::ExportJobRequestedHandler error while unmarshalling payload")
			if ackErr := msg.Ack(); ackErr != nil {
				infraTracing.RecordError(span, ackErr)
				log.Ctx(msgCtx).Error().Err(ackErr).Msg("consumer::ExportJobRequestedHandler error while acknowledging invalid message")
			}
			return
		}

		err = core.ProcessExportJob(msgCtx, coreentity.ExportJobDetailFilter{
			TenantID: payload.TenantID,
			ID:       payload.JobID,
		})
		if err != nil {
			infraTracing.RecordError(span, err)
			log.Ctx(msgCtx).Error().Err(err).Any("payload", payload).Msg("consumer::ExportJobRequestedHandler error while processing export job")
			if nakErr := msg.Nak(err.Error()); nakErr != nil {
				infraTracing.RecordError(span, nakErr)
				log.Ctx(msgCtx).Error().Err(nakErr).Any("payload", payload).Msg("consumer::ExportJobRequestedHandler error while negative-acknowledging message")
			}
			return
		}

		if err := msg.Ack(); err != nil {
			infraTracing.RecordError(span, err)
			log.Ctx(msgCtx).Error().Err(err).Any("payload", payload).Msg("consumer::ExportJobRequestedHandler error while acknowledging message")
		}
	}
}

type IntegrationExportJobProcessor interface {
	ProcessExportJob(ctx context.Context, filter coreentity.ExportJobDetailFilter) error
}
