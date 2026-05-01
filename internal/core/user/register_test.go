package core

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	integrationMocks "codebase-app/internal/ports/integration/mocks"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"
	"codebase-app/pkg/errmsg"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserCore_RegisterOwner(t *testing.T) {
	ctx := context.Background()

	user := coreentity.User{
		Name:       "Owner Example",
		UserName:   "owner",
		Email:      "owner@example.com",
		Password:   "password123",
		TenantName: "PT Contoh",
	}
	tenant := coreentity.Tenant{
		Name:              "PT Contoh",
		Code:              "TCO",
		PreferredLanguage: "id",
	}

	t.Run("success creates owner and issues email verification", func(t *testing.T) {
		repo := dbMocks.NewUserRepository(t)
		tx := dbMocks.NewTransactor(t)
		bus := integrationMocks.NewMessagePublisher(t)

		tx.EXPECT().
			WithinTransaction(mock.Anything, mock.Anything).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			}).
			Twice()
		repo.EXPECT().
			ExistsTenantByCode(mock.Anything, tenant.Code).
			Return(false, nil)
		repo.EXPECT().
			InsertTenant(mock.Anything, coreentity.Tenant{
				Name: tenant.Name,
				Code: tenant.Code,
			}).
			Return("tenant-1", nil)
		repo.EXPECT().
			InitializeTenant(mock.Anything, "tenant-1", tenant.Name, tenant.PreferredLanguage).
			Return("company-1", nil)
		repo.EXPECT().
			GetRoleByName(mock.Anything, "owner", "tenant-1").
			Return(&coreentity.Role{ID: "role-owner", Name: "owner"}, nil)
		repo.EXPECT().
			InsertUser(mock.Anything, mock.MatchedBy(func(got coreentity.User) bool {
				return got.Name == user.Name &&
					got.UserName == user.UserName &&
					got.Email == user.Email &&
					got.TenantID == "tenant-1" &&
					got.TenantName == user.TenantName &&
					len(got.RoleIDs) == 1 &&
					got.RoleIDs[0] == "role-owner" &&
					len(got.RoleNames) == 1 &&
					got.RoleNames[0] == "owner" &&
					got.Password != "" &&
					got.Password != user.Password
			})).
			Return("user-1", nil)
		repo.EXPECT().
			DeleteUserActionTokensByPurpose(
				mock.Anything,
				"user-1",
				coreentity.UserActionTokenPurposeEmailVerification,
			).
			Return(nil)
		repo.EXPECT().
			CreateUserActionToken(mock.Anything, mock.MatchedBy(func(token coreentity.UserActionToken) bool {
				return token.UserID == "user-1" &&
					token.Purpose == coreentity.UserActionTokenPurposeEmailVerification &&
					token.TokenHash != "" &&
					token.ExpiresAt.After(time.Now().UTC().Add(6*24*time.Hour))
			})).
			Return(nil)
		bus.EXPECT().
			PublishJSON(
				mock.Anything,
				common.MessageTopicEmailVerification,
				mock.MatchedBy(func(payload any) bool {
					msg, ok := payload.(coreentity.QueuedEmailMessage)
					return ok &&
						msg.Email == user.Email &&
						msg.UserName == user.Name &&
						msg.TenantName == user.TenantName &&
						msg.PreferredLanguage == tenant.PreferredLanguage &&
						strings.Contains(msg.ActionLink, "token=")
				}),
			).
			Return(nil)

		core := NewUserCore(Config{
			Repo: repo,
			Tx:   tx,
			Bus:  bus,
		})

		got, err := core.RegisterOwner(ctx, user, tenant)

		require.NoError(t, err)
		require.Equal(t, &coreentity.RegisterResult{
			Email:                user.Email,
			VerificationRequired: true,
		}, got)
	})

	tests := []struct {
		name      string
		setup     func(repo *dbMocks.UserRepository, tx *dbMocks.Transactor, bus *integrationMocks.MessagePublisher) error
		wantError func(t *testing.T, err error)
	}{
		{
			name: "rejects duplicate tenant code",
			setup: func(repo *dbMocks.UserRepository, tx *dbMocks.Transactor, bus *integrationMocks.MessagePublisher) error {
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()
				repo.EXPECT().
					ExistsTenantByCode(mock.Anything, tenant.Code).
					Return(true, nil)
				return nil
			},
			wantError: func(t *testing.T, err error) {
				var customErr *errmsg.CustomError
				require.ErrorAs(t, err, &customErr)
				require.Equal(t, 400, customErr.Code)
				require.Equal(t, errmsg.MessageTenantCodeIsAlreadyRegistered, customErr.Msg)
			},
		},
		{
			name: "returns outer transaction error",
			setup: func(repo *dbMocks.UserRepository, tx *dbMocks.Transactor, bus *integrationMocks.MessagePublisher) error {
				txErr := errors.New("transaction failed")
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					Return(txErr).
					Once()
				return txErr
			},
			wantError: func(t *testing.T, err error) {
				require.EqualError(t, err, "transaction failed")
			},
		},
		{
			name: "returns tenant exists check error",
			setup: func(repo *dbMocks.UserRepository, tx *dbMocks.Transactor, bus *integrationMocks.MessagePublisher) error {
				depErr := errors.New("exists check failed")
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()
				repo.EXPECT().
					ExistsTenantByCode(mock.Anything, tenant.Code).
					Return(false, depErr)
				return depErr
			},
			wantError: func(t *testing.T, err error) {
				require.EqualError(t, err, "exists check failed")
			},
		},
		{
			name: "returns owner insert error",
			setup: func(repo *dbMocks.UserRepository, tx *dbMocks.Transactor, bus *integrationMocks.MessagePublisher) error {
				depErr := errors.New("insert user failed")
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Once()
				repo.EXPECT().
					ExistsTenantByCode(mock.Anything, tenant.Code).
					Return(false, nil)
				repo.EXPECT().
					InsertTenant(mock.Anything, coreentity.Tenant{
						Name: tenant.Name,
						Code: tenant.Code,
					}).
					Return("tenant-1", nil)
				repo.EXPECT().
					InitializeTenant(mock.Anything, "tenant-1", tenant.Name, tenant.PreferredLanguage).
					Return("company-1", nil)
				repo.EXPECT().
					GetRoleByName(mock.Anything, "owner", "tenant-1").
					Return(&coreentity.Role{ID: "role-owner", Name: "owner"}, nil)
				repo.EXPECT().
					InsertUser(mock.Anything, mock.Anything).
					Return("", depErr)
				return depErr
			},
			wantError: func(t *testing.T, err error) {
				require.EqualError(t, err, "insert user failed")
			},
		},
		{
			name: "returns email verification delivery error",
			setup: func(repo *dbMocks.UserRepository, tx *dbMocks.Transactor, bus *integrationMocks.MessagePublisher) error {
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).
					Twice()
				repo.EXPECT().
					ExistsTenantByCode(mock.Anything, tenant.Code).
					Return(false, nil)
				repo.EXPECT().
					InsertTenant(mock.Anything, coreentity.Tenant{
						Name: tenant.Name,
						Code: tenant.Code,
					}).
					Return("tenant-1", nil)
				repo.EXPECT().
					InitializeTenant(mock.Anything, "tenant-1", tenant.Name, tenant.PreferredLanguage).
					Return("company-1", nil)
				repo.EXPECT().
					GetRoleByName(mock.Anything, "owner", "tenant-1").
					Return(&coreentity.Role{ID: "role-owner", Name: "owner"}, nil)
				repo.EXPECT().
					InsertUser(mock.Anything, mock.Anything).
					Return("user-1", nil)
				repo.EXPECT().
					DeleteUserActionTokensByPurpose(
						mock.Anything,
						"user-1",
						coreentity.UserActionTokenPurposeEmailVerification,
					).
					Return(nil)
				repo.EXPECT().
					CreateUserActionToken(mock.Anything, mock.Anything).
					Return(nil)
				return nil
			},
			wantError: func(t *testing.T, err error) {
				var customErr *errmsg.CustomError
				require.ErrorAs(t, err, &customErr)
				require.Equal(t, 500, customErr.Code)
				require.Equal(t, errmsg.MessageEmailMessageBusIsNotConfigured, customErr.Msg)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := dbMocks.NewUserRepository(t)
			tx := dbMocks.NewTransactor(t)
			bus := integrationMocks.NewMessagePublisher(t)

			if tt.setup != nil {
				tt.setup(repo, tx, bus)
			}

			core := NewUserCore(Config{
				Repo: repo,
				Tx:   tx,
			})

			got, err := core.RegisterOwner(ctx, user, tenant)

			require.Nil(t, got)
			require.Error(t, err)
			tt.wantError(t, err)
		})
	}
}
