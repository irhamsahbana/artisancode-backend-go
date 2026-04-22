package coreentity

import (
	"time"

	"codebase-app/internal/entity/common"
)

const (
	UserInvitationRoleAdmin    = "admin"
	UserInvitationRoleEmployee = "employee"

	UserInvitationStatusPending  = "pending"
	UserInvitationStatusAccepted = "accepted"
	UserInvitationStatusExpired  = "expired"
	UserInvitationStatusRevoked  = "revoked"
)

type UserInvitation struct {
	UserCtx common.UserContext

	ID           string
	TenantID     string
	TenantCode   string
	TenantName   string
	EmployeeID   *string
	EmployeeNo   *string
	EmployeeName *string
	Email        string
	RoleCode     string
	Status       string
	TokenHash    string
	AcceptToken  string
	ExpiresAt    time.Time
	AcceptedAt   *time.Time
	RevokedAt    *time.Time
	InvitedBy    string
	LastSentAt   time.Time
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	EmailSent    bool
}

type UserInvitationListFilter struct {
	UserCtx     common.UserContext
	TenantID    string
	Q           string
	RoleCode    string
	Status      string
	EmployeeIDs []string
	Page        int
	Paginate    int
}

type UserInvitationAcceptPayload struct {
	Token    string
	Password string
	FullName string
}
