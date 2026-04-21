package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type UserInvitationRepository interface {
	CreateInvitation(ctx context.Context, data coreentity.UserInvitation) (*coreentity.UserInvitation, error)
	GetInvitations(ctx context.Context, filter coreentity.UserInvitationListFilter) ([]coreentity.UserInvitation, int, error)
	GetInvitationByID(ctx context.Context, tenantID, invitationID string) (*coreentity.UserInvitation, error)
	GetInvitationByTokenHash(ctx context.Context, tokenHash string) (*coreentity.UserInvitation, error)
	ExistsActiveInvitation(ctx context.Context, tenantID, email, roleCode string, employeeID *string) (bool, error)
	ResendInvitation(ctx context.Context, data coreentity.UserInvitation) error
	RevokeInvitation(ctx context.Context, tenantID, invitationID string) error
	MarkInvitationAccepted(ctx context.Context, invitationID string) error
}
