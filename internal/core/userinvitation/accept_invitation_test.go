package core

import (
	"context"
	"errors"
	"testing"
	"time"

	"codebase-app/internal/entity/coreentity"
	dbmocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUserInvitationCore_AcceptInvitation(t *testing.T) {
	ctx := context.Background()
	token := "raw-token"

	tests := []struct {
		name      string
		input     coreentity.UserInvitationAcceptPayload
		setup     func(repo *dbmocks.UserInvitationRepository, userRepo *dbmocks.UserRepository, employeeRepo *dbmocks.EmployeeRepository, tx *dbmocks.Transactor)
		wantError bool
	}{
		{
			name: "success for admin invitation",
			input: coreentity.UserInvitationAcceptPayload{
				Token:    token,
				Password: "secret123",
				FullName: "New Admin",
			},
			setup: func(repo *dbmocks.UserInvitationRepository, userRepo *dbmocks.UserRepository, employeeRepo *dbmocks.EmployeeRepository, tx *dbmocks.Transactor) {
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(txCtx context.Context, fn func(context.Context) error) error {
						return fn(txCtx)
					})
				repo.EXPECT().
					GetInvitationByTokenHash(mock.Anything, hashInvitationToken(token)).
					Return(&coreentity.UserInvitation{
						ID:        "invitation-1",
						TenantID:  "tenant-1",
						Email:     "new.admin@example.com",
						RoleCode:  coreentity.UserInvitationRoleAdmin,
						Status:    coreentity.UserInvitationStatusPending,
						ExpiresAt: time.Now().Add(time.Hour),
					}, nil)
				userRepo.EXPECT().
					FindActiveUserByEmailAndTenantID(mock.Anything, "new.admin@example.com", "tenant-1").
					Return(nil, nil)
				userRepo.EXPECT().
					GetRoleByName(mock.Anything, coreentity.UserInvitationRoleAdmin, "tenant-1").
					Return(&coreentity.Role{ID: "role-admin", Name: "admin"}, nil)
				userRepo.EXPECT().
					InsertUser(mock.Anything, mock.MatchedBy(func(user coreentity.User) bool {
						return user.Email == "new.admin@example.com" &&
							user.TenantID == "tenant-1" &&
							user.Name == "New Admin" &&
							user.UserName == "new.admin" &&
							len(user.RoleIDs) == 1 &&
							user.RoleIDs[0] == "role-admin" &&
							bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("secret123")) == nil
					})).
					Return("user-1", nil)
				userRepo.EXPECT().
					MarkUserEmailVerified(mock.Anything, "user-1").
					Return(nil)
				repo.EXPECT().
					MarkInvitationAccepted(mock.Anything, "invitation-1").
					Return(nil)
				userRepo.EXPECT().
					GetUser(mock.Anything, coreentity.User{
						ID:       "user-1",
						TenantID: "tenant-1",
					}).
					Return(&coreentity.User{
						ID:       "user-1",
						Email:    "new.admin@example.com",
						TenantID: "tenant-1",
						Name:     "New Admin",
					}, nil)
			},
		},
		{
			name: "dependency error from invitation lookup",
			input: coreentity.UserInvitationAcceptPayload{
				Token:    token,
				Password: "secret123",
				FullName: "New Admin",
			},
			setup: func(repo *dbmocks.UserInvitationRepository, userRepo *dbmocks.UserRepository, employeeRepo *dbmocks.EmployeeRepository, tx *dbmocks.Transactor) {
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(txCtx context.Context, fn func(context.Context) error) error {
						return fn(txCtx)
					})
				repo.EXPECT().
					GetInvitationByTokenHash(mock.Anything, hashInvitationToken(token)).
					Return(nil, errors.New("repository failed"))
			},
			wantError: true,
		},
		{
			name: "business error when email is already registered",
			input: coreentity.UserInvitationAcceptPayload{
				Token:    token,
				Password: "secret123",
				FullName: "New Admin",
			},
			setup: func(repo *dbmocks.UserInvitationRepository, userRepo *dbmocks.UserRepository, employeeRepo *dbmocks.EmployeeRepository, tx *dbmocks.Transactor) {
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(txCtx context.Context, fn func(context.Context) error) error {
						return fn(txCtx)
					})
				repo.EXPECT().
					GetInvitationByTokenHash(mock.Anything, hashInvitationToken(token)).
					Return(&coreentity.UserInvitation{
						ID:        "invitation-1",
						TenantID:  "tenant-1",
						Email:     "new.admin@example.com",
						RoleCode:  coreentity.UserInvitationRoleAdmin,
						Status:    coreentity.UserInvitationStatusPending,
						ExpiresAt: time.Now().Add(time.Hour),
					}, nil)
				userRepo.EXPECT().
					FindActiveUserByEmailAndTenantID(mock.Anything, "new.admin@example.com", "tenant-1").
					Return(&coreentity.User{ID: "existing-user"}, nil)
			},
			wantError: true,
		},
		{
			name: "business error when full name is empty for admin invitation",
			input: coreentity.UserInvitationAcceptPayload{
				Token:    token,
				Password: "secret123",
				FullName: " ",
			},
			setup: func(repo *dbmocks.UserInvitationRepository, userRepo *dbmocks.UserRepository, employeeRepo *dbmocks.EmployeeRepository, tx *dbmocks.Transactor) {
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(txCtx context.Context, fn func(context.Context) error) error {
						return fn(txCtx)
					})
				repo.EXPECT().
					GetInvitationByTokenHash(mock.Anything, hashInvitationToken(token)).
					Return(&coreentity.UserInvitation{
						ID:        "invitation-1",
						TenantID:  "tenant-1",
						Email:     "new.admin@example.com",
						RoleCode:  coreentity.UserInvitationRoleAdmin,
						Status:    coreentity.UserInvitationStatusPending,
						ExpiresAt: time.Now().Add(time.Hour),
					}, nil)
				userRepo.EXPECT().
					FindActiveUserByEmailAndTenantID(mock.Anything, "new.admin@example.com", "tenant-1").
					Return(nil, nil)
				userRepo.EXPECT().
					GetRoleByName(mock.Anything, coreentity.UserInvitationRoleAdmin, "tenant-1").
					Return(&coreentity.Role{ID: "role-admin", Name: "admin"}, nil)
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

			if tt.setup != nil {
				tt.setup(repo, userRepo, employeeRepo, tx)
			}

			core := NewUserInvitationCore(Config{
				Repo:         repo,
				UserRepo:     userRepo,
				EmployeeRepo: employeeRepo,
				Tx:           tx,
			})

			got, err := core.AcceptInvitation(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, "user-1", got.ID)
			require.Equal(t, "new.admin@example.com", got.Email)
		})
	}
}
