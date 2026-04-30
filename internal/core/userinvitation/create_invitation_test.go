package core

import (
	"context"
	"errors"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/config"
	integrationmocks "codebase-app/internal/ports/integration/mocks"
	dbmocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserInvitationCore_CreateInvitation(t *testing.T) {
	ctx := context.Background()
	ownerCtx := common.UserContext{
		UserID:   "owner-1",
		TenantID: "tenant-1",
		Roles:    []string{"owner"},
	}
	adminCtx := common.UserContext{
		UserID:   "admin-1",
		TenantID: "tenant-1",
		Roles:    []string{"admin"},
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
			name: "success for owner inviting admin",
			input: coreentity.UserInvitation{
				UserCtx:  ownerCtx,
				Email:    " New.Admin@Example.COM ",
				RoleCode: " ADMIN ",
			},
			setup: func(repo *dbmocks.UserInvitationRepository, userRepo *dbmocks.UserRepository, tx *dbmocks.Transactor, bus *integrationmocks.MessagePublisher) {
				userRepo.EXPECT().
					FindActiveUserByEmailAndTenantID(mock.Anything, "new.admin@example.com", "tenant-1").
					Return(nil, nil)
				userRepo.EXPECT().
					GetRoleByName(mock.Anything, "admin", "tenant-1").
					Return(&coreentity.Role{ID: "role-admin", Name: "admin"}, nil)
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(txCtx context.Context, fn func(context.Context) error) error {
						return fn(txCtx)
					})
				repo.EXPECT().
					ExistsActiveInvitation(mock.Anything, "tenant-1", "new.admin@example.com", "admin", (*string)(nil)).
					Return(false, nil)
				repo.EXPECT().
					CreateInvitation(mock.Anything, mock.MatchedBy(func(data coreentity.UserInvitation) bool {
						return data.Email == "new.admin@example.com" &&
							data.RoleCode == "admin" &&
							data.TenantID == "tenant-1" &&
							data.InvitedBy == "owner-1" &&
							data.Status == coreentity.UserInvitationStatusPending &&
							data.TokenHash != "" &&
							data.AcceptToken != ""
					})).
					RunAndReturn(func(ctx context.Context, data coreentity.UserInvitation) (*coreentity.UserInvitation, error) {
						data.ID = "invitation-1"
						data.TenantName = "Tenant One"
						return &data, nil
					})
				userRepo.EXPECT().
					GetTenantPreferredLanguage(mock.Anything, "tenant-1").
					Return("id", nil)
				bus.EXPECT().
					PublishJSON(mock.Anything, common.MessageTopicEmailInvitation, mock.Anything).
					Return(nil)
			},
		},
		{
			name: "authorization error when requester cannot invite users",
			input: coreentity.UserInvitation{
				UserCtx:  memberCtx,
				Email:    "person@example.com",
				RoleCode: "employee",
			},
			wantError: true,
		},
		{
			name: "authorization error when admin invites admin",
			input: coreentity.UserInvitation{
				UserCtx:  adminCtx,
				Email:    "new.admin@example.com",
				RoleCode: "admin",
			},
			wantError: true,
		},
		{
			name: "dependency error from user lookup",
			input: coreentity.UserInvitation{
				UserCtx:  ownerCtx,
				Email:    "person@example.com",
				RoleCode: "admin",
			},
			setup: func(repo *dbmocks.UserInvitationRepository, userRepo *dbmocks.UserRepository, tx *dbmocks.Transactor, bus *integrationmocks.MessagePublisher) {
				userRepo.EXPECT().
					FindActiveUserByEmailAndTenantID(mock.Anything, "person@example.com", "tenant-1").
					Return(nil, errors.New("user lookup failed"))
			},
			wantError: true,
		},
		{
			name: "business error when active invitation already exists",
			input: coreentity.UserInvitation{
				UserCtx:  ownerCtx,
				Email:    "person@example.com",
				RoleCode: "admin",
			},
			setup: func(repo *dbmocks.UserInvitationRepository, userRepo *dbmocks.UserRepository, tx *dbmocks.Transactor, bus *integrationmocks.MessagePublisher) {
				userRepo.EXPECT().
					FindActiveUserByEmailAndTenantID(mock.Anything, "person@example.com", "tenant-1").
					Return(nil, nil)
				userRepo.EXPECT().
					GetRoleByName(mock.Anything, "admin", "tenant-1").
					Return(&coreentity.Role{ID: "role-admin", Name: "admin"}, nil)
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(txCtx context.Context, fn func(context.Context) error) error {
						return fn(txCtx)
					})
				repo.EXPECT().
					ExistsActiveInvitation(mock.Anything, "tenant-1", "person@example.com", "admin", (*string)(nil)).
					Return(true, nil)
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := dbmocks.NewUserInvitationRepository(t)
			userRepo := dbmocks.NewUserRepository(t)
			employeeRepo := dbmocks.NewEmployeeRepository(t)
			tx := dbmocks.NewTransactor(t)
			bus := integrationmocks.NewMessagePublisher(t)

			if tt.setup != nil {
				tt.setup(repo, userRepo, tx, bus)
			}

			core := NewUserInvitationCore(Config{
				Repo:         repo,
				UserRepo:     userRepo,
				EmployeeRepo: employeeRepo,
				Tx:           tx,
				Bus:          bus,
			})

			got, err := core.CreateInvitation(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, "invitation-1", got.ID)
			require.Equal(t, "new.admin@example.com", got.Email)
			require.Equal(t, "admin", got.RoleCode)
			require.True(t, got.EmailSent)
			require.NotEmpty(t, got.AcceptToken)
		})
	}
}
