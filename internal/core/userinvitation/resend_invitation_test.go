package core

import (
	"context"
	"errors"
	"testing"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/config"
	integrationmocks "codebase-app/internal/ports/integration/mocks"
	dbmocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserInvitationCore_ResendInvitation(t *testing.T) {
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

	config.Envs = &config.Config{}
	config.Envs.FrontendURL.ClientBaseURL = "http://frontend.test"
	config.Envs.FrontendURL.Invitation = "/invite"

	tests := []struct {
		name      string
		input     coreentity.UserInvitation
		setup     func(repo *dbmocks.UserInvitationRepository, userRepo *dbmocks.UserRepository, tx *dbmocks.Transactor, bus *integrationmocks.MessagePublisher)
		wantError bool
	}{
		{
			name: "success",
			input: coreentity.UserInvitation{
				UserCtx: ownerCtx,
				ID:      "invitation-1",
			},
			setup: func(repo *dbmocks.UserInvitationRepository, userRepo *dbmocks.UserRepository, tx *dbmocks.Transactor, bus *integrationmocks.MessagePublisher) {
				repo.EXPECT().
					GetInvitationByID(mock.Anything, "tenant-1", "invitation-1").
					Return(&coreentity.UserInvitation{
						ID:        "invitation-1",
						TenantID:  "tenant-1",
						Email:     "person@example.com",
						RoleCode:  coreentity.UserInvitationRoleAdmin,
						Status:    coreentity.UserInvitationStatusPending,
						ExpiresAt: time.Now().Add(time.Hour),
					}, nil).
					Once()
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(txCtx context.Context, fn func(context.Context) error) error {
						return fn(txCtx)
					})
				repo.EXPECT().
					ResendInvitation(mock.Anything, mock.MatchedBy(func(data coreentity.UserInvitation) bool {
						return data.ID == "invitation-1" &&
							data.TokenHash != "" &&
							data.AcceptToken != "" &&
							data.ExpiresAt.After(time.Now()) &&
							!data.LastSentAt.IsZero()
					})).
					Return(nil)
				repo.EXPECT().
					GetInvitationByID(mock.Anything, "tenant-1", "invitation-1").
					Return(&coreentity.UserInvitation{
						ID:         "invitation-1",
						TenantID:   "tenant-1",
						TenantName: "Tenant One",
						Email:      "person@example.com",
						RoleCode:   coreentity.UserInvitationRoleAdmin,
						Status:     coreentity.UserInvitationStatusPending,
						ExpiresAt:  time.Now().Add(invitationExpiryDuration),
					}, nil).
					Once()
				userRepo.EXPECT().
					GetTenantPreferredLanguage(mock.Anything, "tenant-1").
					Return("id", nil)
				bus.EXPECT().
					PublishJSON(mock.Anything, common.MessageTopicEmailInvitation, mock.Anything).
					Return(nil)
			},
		},
		{
			name: "authorization error",
			input: coreentity.UserInvitation{
				UserCtx: memberCtx,
				ID:      "invitation-1",
			},
			setup: func(repo *dbmocks.UserInvitationRepository, userRepo *dbmocks.UserRepository, tx *dbmocks.Transactor, bus *integrationmocks.MessagePublisher) {
				repo.EXPECT().
					GetInvitationByID(mock.Anything, "tenant-1", "invitation-1").
					Return(&coreentity.UserInvitation{
						ID:       "invitation-1",
						TenantID: "tenant-1",
						Email:    "person@example.com",
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
			setup: func(repo *dbmocks.UserInvitationRepository, userRepo *dbmocks.UserRepository, tx *dbmocks.Transactor, bus *integrationmocks.MessagePublisher) {
				repo.EXPECT().
					GetInvitationByID(mock.Anything, "tenant-1", "invitation-1").
					Return(&coreentity.UserInvitation{
						ID:       "invitation-1",
						TenantID: "tenant-1",
						Email:    "person@example.com",
						RoleCode: coreentity.UserInvitationRoleAdmin,
						Status:   coreentity.UserInvitationStatusAccepted,
					}, nil)
			},
			wantError: true,
		},
		{
			name: "dependency error from invitation lookup",
			input: coreentity.UserInvitation{
				UserCtx: ownerCtx,
				ID:      "invitation-1",
			},
			setup: func(repo *dbmocks.UserInvitationRepository, userRepo *dbmocks.UserRepository, tx *dbmocks.Transactor, bus *integrationmocks.MessagePublisher) {
				repo.EXPECT().
					GetInvitationByID(mock.Anything, "tenant-1", "invitation-1").
					Return(nil, errors.New("repository failed"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := dbmocks.NewUserInvitationRepository(t)
			userRepo := dbmocks.NewUserRepository(t)
			tx := dbmocks.NewTransactor(t)
			bus := integrationmocks.NewMessagePublisher(t)

			if tt.setup != nil {
				tt.setup(repo, userRepo, tx, bus)
			}

			core := NewUserInvitationCore(Config{
				Repo:     repo,
				UserRepo: userRepo,
				Tx:       tx,
				Bus:      bus,
			})

			got, err := core.ResendInvitation(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, "invitation-1", got.ID)
			require.True(t, got.EmailSent)
			require.NotEmpty(t, got.AcceptToken)
		})
	}
}
