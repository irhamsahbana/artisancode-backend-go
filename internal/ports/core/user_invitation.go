package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type UserInvitationCore interface {
	CreateInvitation(ctx context.Context, data coreentity.UserInvitation) (*coreentity.UserInvitation, error)
	GetInvitations(ctx context.Context, filter coreentity.UserInvitationListFilter) ([]coreentity.UserInvitation, int, error)
	ResendInvitation(ctx context.Context, data coreentity.UserInvitation) (*coreentity.UserInvitation, error)
	RevokeInvitation(ctx context.Context, data coreentity.UserInvitation) error
	GetInvitationByToken(ctx context.Context, token string) (*coreentity.UserInvitation, error)
	AcceptInvitation(ctx context.Context, data coreentity.UserInvitationAcceptPayload) (*coreentity.User, error)
}
