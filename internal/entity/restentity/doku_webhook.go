package restentity

import (
	"encoding/json"
	"strconv"
	"strings"
)

type DOKUWebhookNotification struct {
	Order       DOKUWebhookOrder       `json:"order"`
	Transaction DOKUWebhookTransaction `json:"transaction"`
}

type DOKUWebhookOrder struct {
	InvoiceNumber string        `json:"invoice_number"`
	Amount        FlexibleInt64 `json:"amount"`
	Status        string        `json:"status"`
}

type DOKUWebhookTransaction struct {
	Status            string `json:"status"`
	Date              string `json:"date"`
	OriginalRequestID string `json:"original_request_id"`
}

type FlexibleInt64 int64

func (v *FlexibleInt64) UnmarshalJSON(data []byte) error {
	raw := strings.TrimSpace(string(data))
	if raw == "" || raw == "null" {
		*v = 0
		return nil
	}

	var intValue int64
	if err := json.Unmarshal(data, &intValue); err == nil {
		*v = FlexibleInt64(intValue)
		return nil
	}

	var strValue string
	if err := json.Unmarshal(data, &strValue); err != nil {
		return err
	}

	strValue = strings.TrimSpace(strValue)
	if strValue == "" {
		*v = 0
		return nil
	}

	parsed, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		return err
	}

	*v = FlexibleInt64(parsed)
	return nil
}
