package mapper

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
)

func DOKUWebhookNotificationToCore(
	targetPath string,
	rawBody []byte,
	headers coreentity.DOKUWebhookSignatureHeaders,
	event restentity.DOKUWebhookNotification,
) coreentity.DOKUWebhookNotification {
	return coreentity.DOKUWebhookNotification{
		TargetPath: targetPath,
		RawBody:    rawBody,
		Headers:    headers,
		Event: coreentity.DOKUWebhookEvent{
			Order: coreentity.DOKUWebhookOrder{
				InvoiceNumber: event.Order.InvoiceNumber,
				Amount:        int64(event.Order.Amount),
				Status:        event.Order.Status,
			},
			Transaction: coreentity.DOKUWebhookTransaction{
				Status:            event.Transaction.Status,
				Date:              event.Transaction.Date,
				OriginalRequestID: event.Transaction.OriginalRequestID,
			},
		},
	}
}
