package coreentity

import "time"

const (
	UserActionTokenPurposeEmailVerification = "email_verification"
	UserActionTokenPurposePasswordReset     = "password_reset"
)

type UserActionToken struct {
	ID         string
	UserID     string
	TenantID   string
	Email      string
	UserName   string
	TenantName string
	Purpose    string
	Token      string
	TokenHash  string
	ExpiresAt  time.Time
	UsedAt     *time.Time
}

type RegisterResult struct {
	Email                string
	VerificationRequired bool
}
