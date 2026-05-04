package restentity

import (
	"errors"
	"strings"
)

const (
	defaultExpiryMins = 10080
)

type DokuLineItem struct {
	Name     string `json:"name"`
	Price    int64  `json:"price"`
	Quantity int    `json:"quantity"`
}

type DokuCreatePaymentRequest struct {
	InvoiceNumber   string         `json:"invoice_number"`
	Amount          int64          `json:"amount"`
	Currency        string         `json:"currency"`
	CustomerEmail   string         `json:"customer_email"`
	CustomerName    string         `json:"customer_name"`
	CustomerPhone   string         `json:"customer_phone,omitempty"`
	CustomerAddress string         `json:"customer_address,omitempty"`
	CustomerCountry string         `json:"customer_country,omitempty"`
	LineItems       []DokuLineItem `json:"line_items,omitempty"`
	ExpiryMinutes   int            `json:"expiry_minutes,omitempty"`
	CallbackURL     string         `json:"callback_url,omitempty"`
	AutoRedirect    *bool          `json:"auto_redirect,omitempty"`
}

func (r DokuCreatePaymentRequest) validate() error {
	switch {
	case strings.TrimSpace(r.InvoiceNumber) == "":
		return errors.New("invoice number is required")
	case r.Amount <= 0:
		return errors.New("amount must be greater than zero")
	case len(strings.TrimSpace(r.Currency)) != 3 || strings.ToUpper(strings.TrimSpace(r.Currency)) != strings.TrimSpace(r.Currency):
		return errors.New("currency must be a 3-letter uppercase code")
	case strings.TrimSpace(r.CustomerEmail) == "":
		return errors.New("customer email is required")
	case strings.TrimSpace(r.CustomerName) == "":
		return errors.New("customer name is required")
	default:
		return nil
	}
}

func (r DokuCreatePaymentRequest) ExpiryMinutesOrDefault() int {
	if r.ExpiryMinutes > 0 {
		return r.ExpiryMinutes
	}

	return defaultExpiryMins
}

func (r DokuCreatePaymentRequest) CustomerCountryOrDefault() string {
	if strings.TrimSpace(r.CustomerCountry) != "" {
		return r.CustomerCountry
	}

	return "ID"
}

func (r DokuCreatePaymentRequest) AutoRedirectOrDefault() bool {
	if r.AutoRedirect != nil {
		return *r.AutoRedirect
	}

	return true
}

type DokuCreatePaymentResponse struct {
	InvoiceID     string   `json:"invoice_id"`
	InvoiceNumber string   `json:"invoice_number"`
	Amount        int64    `json:"amount"`
	PaymentURL    string   `json:"payment_url"`
	RequestID     string   `json:"request_id"`
	Message       []string `json:"message,omitempty"`
}

type DokuCheckStatusResponse struct {
	Order struct {
		InvoiceNumber string        `json:"invoice_number,omitempty"`
		Amount        FlexibleInt64 `json:"amount,omitempty"`
		Status        string        `json:"status,omitempty"`
	} `json:"order"`
	Transaction struct {
		Status            string `json:"status,omitempty"`
		Date              string `json:"date,omitempty"`
		OriginalRequestID string `json:"original_request_id,omitempty"`
	} `json:"transaction"`
}

type DokuCheckoutRequest struct {
	Order    DokuCheckoutOrder    `json:"order"`
	Payment  DokuCheckoutPayment  `json:"payment"`
	Customer DokuCheckoutCustomer `json:"customer"`
}

type DokuCheckoutOrder struct {
	Amount        int64          `json:"amount"`
	InvoiceNumber string         `json:"invoice_number"`
	Currency      string         `json:"currency"`
	CallbackURL   string         `json:"callback_url"`
	AutoRedirect  bool           `json:"auto_redirect"`
	LineItems     []DokuLineItem `json:"line_items,omitempty"`
}

type DokuCheckoutPayment struct {
	PaymentDueDate int `json:"payment_due_date"`
}

type DokuCheckoutCustomer struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone,omitempty"`
	Address string `json:"address,omitempty"`
	Country string `json:"country"`
}

type DokuCheckoutResponseEnvelope struct {
	Message  []string `json:"message"`
	Response struct {
		Order struct {
			InvoiceNumber string        `json:"invoice_number"`
			Amount        FlexibleInt64 `json:"amount"`
		} `json:"order"`
		Payment struct {
			URL string `json:"url"`
		} `json:"payment"`
	} `json:"response"`
}
