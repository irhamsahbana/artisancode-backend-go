package coreentity

type DOKUWebhookNotification struct {
	TargetPath string
	RawBody    []byte
	Headers    DOKUWebhookSignatureHeaders
	Event      DOKUWebhookEvent
}

type DOKUWebhookSignatureHeaders struct {
	ClientID         string
	RequestID        string
	RequestTimestamp string
	Signature        string
}

type DOKUWebhookEvent struct {
	Order       DOKUWebhookOrder
	Transaction DOKUWebhookTransaction
}

type DOKUWebhookOrder struct {
	InvoiceNumber string
	Amount        int64
	Status        string
}

type DOKUWebhookTransaction struct {
	Status            string
	Date              string
	OriginalRequestID string
}

type DOKUWebhookResult struct {
	Provider      string
	InvoiceNumber string
	OrderStatus   string
	PaymentStatus string
}
