package restentity

import (
	"time"

	"codebase-app/pkg/types"
)

type UserInvitationResource struct {
	ID           string     `json:"id"`
	TenantCode   string     `json:"tenant_code"`
	TenantName   string     `json:"tenant_name"`
	EmployeeID   *string    `json:"employee_id"`
	EmployeeNo   *string    `json:"employee_no"`
	EmployeeName *string    `json:"employee_name"`
	Email        string     `json:"email"`
	RoleCode     string     `json:"role_code"`
	Status       string     `json:"status"`
	ExpiresAt    time.Time  `json:"expires_at"`
	AcceptedAt   *time.Time `json:"accepted_at"`
	RevokedAt    *time.Time `json:"revoked_at"`
	LastSentAt   time.Time  `json:"last_sent_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

type CreateUserInvitationReq struct {
	Email      string  `json:"email" validate:"required,email"`
	RoleCode   string  `json:"role_code" validate:"required,oneof=admin employee"`
	EmployeeID *string `json:"employee_id" validate:"omitempty,uuidv7"`
}

type CreateUserInvitationResp struct {
	ID          string    `json:"id"`
	AcceptToken string    `json:"accept_token"`
	ExpiresAt   time.Time `json:"expires_at"`
	EmailSent   bool      `json:"email_sent"`
}

type GetUserInvitationsReq struct {
	Q           string `query:"q" validate:"omitempty,min=2"`
	RoleCode    string `query:"role_code" validate:"omitempty,oneof=admin employee"`
	Status      string `query:"status" validate:"omitempty,oneof=pending accepted expired revoked"`
	EmployeeIDs string `query:"employee_ids"`
	types.MetaQuery
}

func (r *GetUserInvitationsReq) SetDefault() {
	r.MetaQuery.SetDefault()
}

type GetUserInvitationsResp struct {
	Items []UserInvitationResource `json:"items"`
	Meta  types.Meta               `json:"meta"`
}

type ResendUserInvitationReq struct {
	ID string `params:"id" validate:"required,uuidv7"`
}

type ResendUserInvitationResp struct {
	ID          string    `json:"id"`
	AcceptToken string    `json:"accept_token"`
	ExpiresAt   time.Time `json:"expires_at"`
	EmailSent   bool      `json:"email_sent"`
}

type RevokeUserInvitationReq struct {
	ID string `params:"id" validate:"required,uuidv7"`
}

type AcceptUserInvitationPreviewReq struct {
	Token string `query:"token" validate:"required"`
}

type AcceptUserInvitationReq struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"omitempty,min=3"`
}

type AcceptUserInvitationResp struct {
	UserID string `json:"user_id"`
}
