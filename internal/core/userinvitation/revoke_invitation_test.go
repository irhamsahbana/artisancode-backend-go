package core

import (
	"context"
	"errors"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	dbmocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserInvitationCore_RevokeInvitation(t *testing.T) {
	ctx := context.Background()
	ownerCtx := common.UserContext{
		UserID:   "owner-1",
		TenantID: "tenant-1",
		Roles:    []string{"owner"},
	}
	memberCtx := common.UserContext{
		UserID:   "member-1",
		TenantID: "tenant-1",
		Roles:    []string{"employee"},
	}

	tests := []struct {
		name      string
		input     coreentity.UserInvitation
		setup     func(repo *dbmocks.UserInvitationRepository)
		wantError bool
	}{
		{
			name: "success",
			input: coreentity.UserInvitation{
				UserCtx: ownerCtx,
				ID:      "invitation-1",
			},
			setup: func(repo *dbmocks.UserInvitationRepository) {
				repo.EXPECT().
					GetInvitationByID(mock.Anything, "tenant-1", "invitation-1").
					Return(&coreentity.UserInvitation{
						ID:       "invitation-1",
						TenantID: "tenant-1",
						RoleCode: coreentity.UserInvitationRoleAdmin,
						Status:   coreentity.UserInvitationStatusPending,
					}, nil)
				repo.EXPECT().
					RevokeInvitation(mock.Anything, "tenant-1", "invitation-1").
					Return(nil)
			},
		},
		{
			name: "authorization error",
			input: coreentity.UserInvitation{
				UserCtx: memberCtx,
				ID:      "invitation-1",
			},
			setup: func(repo *dbmocks.UserInvitationRepository) {
				repo.EXPECT().
					GetInvitationByID(mock.Anything, "tenant-1", "invitation-1").
					Return(&coreentity.UserInvitation{
						ID:       "invitation-1",
						TenantID: "tenant-1",
						RoleCode: coreentity.UserInvitationRoleEmployee,
						Status:   coreentity.UserInvitationStatusPending,
					}, nil)
			},
			wantError: true,
		},
		{
			name: "business error when invitation is accepted",
			input: coreentity.UserInvitation{
				UserCtx: ownerCtx,
				ID:      "invitation-1",
			},
			setup: func(repo *dbmocks.UserInvitationRepository) {
				repo.EXPECT().
					GetInvitationByID(mock.Anything, "tenant-1", "invitation-1").
					Return(&coreentity.UserInvitation{
						ID:       "invitation-1",
						TenantID: "tenant-1",
						RoleCode: coreentity.UserInvitationRoleAdmin,
						Status:   coreentity.UserInvitationStatusAccepted,
					}, nil)
			},
			wantError: true,
		},
		{
			name: "dependency error from revoke",
			input: coreentity.UserInvitation{
				UserCtx: ownerCtx,
				ID:      "invitation-1",
			},
			setup: func(repo *dbmocks.UserInvitationRepository) {
				repo.EXPECT().
					GetInvitationByID(mock.Anything, "tenant-1", "invitation-1").
					Return(&coreentity.UserInvitation{
						ID:       "invitation-1",
						TenantID: "tenant-1",
						RoleCode: coreentity.UserInvitationRoleAdmin,
						Status:   coreentity.UserInvitationStatusPending,
					}, nil)
				repo.EXPECT().
					RevokeInvitation(mock.Anything, "tenant-1", "invitation-1").
					Return(errors.New("repository failed"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := dbmocks.NewUserInvitationRepository(t)
			if tt.setup != nil {
				tt.setup(repo)
			}

			core := NewUserInvitationCore(Config{
				Repo: repo,
			})

			err := core.RevokeInvitation(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
