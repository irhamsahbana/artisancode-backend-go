package coreentity

type QueuedEmailMessage struct {
	Email             string `json:"email"`
	UserName          string `json:"user_name"`
	TenantName        string `json:"tenant_name"`
	ActionLink        string `json:"action_link"`
	PreferredLanguage string `json:"preferred_language"`
}
