package restentity

type TenantBillingPlan struct {
	ID           string                   `json:"id"`
	Name         string                   `json:"name"`
	Description  string                   `json:"description"`
	PricingID    string                   `json:"pricing_id,omitempty"`
	BillingCycle string                   `json:"billing_cycle,omitempty"`
	Amount       string                   `json:"amount,omitempty"`
	Currency     string                   `json:"currency,omitempty"`
	Features     []string                 `json:"features"`
	Prices       []TenantBillingPlanPrice `json:"prices"`
	AddOns       []TenantBillingAddOn     `json:"add_ons"`
}

type TenantBillingPlanPrice struct {
	ID                    string `json:"id"`
	PricingID             string `json:"pricing_id"`
	BillingCycle          string `json:"billing_cycle"`
	Amount                string `json:"amount"`
	Currency              string `json:"currency"`
	CurrencySymbol        string `json:"currency_symbol"`
	CurrencyDecimalPlaces int    `json:"currency_decimal_places"`
	IsDefaultCurrency     bool   `json:"is_default_currency"`
}

type TenantBillingAddOn struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	Amount       string `json:"amount,omitempty"`
	Currency     string `json:"currency,omitempty"`
	BillingCycle string `json:"billing_cycle,omitempty"`
}

type TenantBillingSubscription struct {
	ID                 string               `json:"id,omitempty"`
	Status             string               `json:"status"`
	PlanName           string               `json:"plan_name,omitempty"`
	BillingCycle       string               `json:"billing_cycle,omitempty"`
	CurrentPeriodStart *string              `json:"current_period_start,omitempty"`
	CurrentPeriodEnd   *string              `json:"current_period_end,omitempty"`
	NextRenewalAt      *string              `json:"next_renewal_at,omitempty"`
	AddOns             []TenantBillingAddOn `json:"add_ons"`
}

type TenantBillingEntitlements struct {
	SubscriptionState string           `json:"subscription_state"`
	Features          []string         `json:"features"`
	UsageLimits       map[string]int64 `json:"usage_limits"`
}

type TenantBillingInvoice struct {
	ID               string  `json:"id"`
	InvoiceNumber    string  `json:"invoice_number"`
	Status           string  `json:"status"`
	Amount           string  `json:"amount"`
	Currency         string  `json:"currency"`
	PaymentAttemptID *string `json:"payment_attempt_id,omitempty"`
	PaymentURL       *string `json:"payment_url,omitempty"`
	DueAt            *string `json:"due_at,omitempty"`
	PaidAt           *string `json:"paid_at,omitempty"`
	CreatedAt        string  `json:"created_at"`
}

type GetTenantBillingInvoicesReq struct {
	Page     int `query:"page"`
	Paginate int `query:"paginate"`
	Limit    int `query:"limit"`
}

func (r *GetTenantBillingInvoicesReq) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}
	if r.Paginate < 1 {
		r.Paginate = r.Limit
	}
	if r.Paginate < 1 {
		r.Paginate = 20
	}
	if r.Paginate > 100 {
		r.Paginate = 100
	}
}

type GetTenantBillingInvoicesResp struct {
	Items []TenantBillingInvoice `json:"items"`
}

type CreateTenantBillingCheckoutReq struct {
	PriceID               string   `json:"price_id" validate:"required"`
	AddOnIDs              []string `json:"add_on_ids"`
	ReplaceActiveCheckout bool     `json:"replace_active_checkout"`
	CallbackURL           string   `json:"callback_url"`
}

type CreateTenantBillingCheckoutResp struct {
	InvoiceID               string `json:"invoice_id"`
	PaymentAttemptID        string `json:"payment_attempt_id"`
	PaymentURL              string `json:"payment_url"`
	TargetSubscriptionState string `json:"target_subscription_state"`
	CheckoutReused          bool   `json:"checkout_reused"`
}
