package core

import (
	"context"
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

func TestUserCore_ForgotPassword(t *testing.T) {
	ctx := context.Background()

	t.Run("success issues reset token and queues reset email", func(t *testing.T) {
		repo := dbMocks.NewUserRepository(t)
		tx := dbMocks.NewTransactor(t)
		bus := integrationMocks.NewMessagePublisher(t)

		repo.EXPECT().
			FindActiveUserByEmailAndTenant(mock.Anything, "user@example.com", "tenant-1").
			Return(&coreentity.User{
				ID:         "user-1",
				Name:       "User Example",
				Email:      "user@example.com",
				TenantID:   "tenant-1",
				TenantName: "Acme Corp",
			}, nil)
		repo.EXPECT().
			GetTenantPreferredLanguage(mock.Anything, "tenant-1").
			Return("id", nil)
		tx.EXPECT().
			WithinTransaction(mock.Anything, mock.Anything).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})
		repo.EXPECT().
			DeleteUserActionTokensByPurpose(mock.Anything, "user-1", coreentity.UserActionTokenPurposePasswordReset).
			Return(nil)
		repo.EXPECT().
			CreateUserActionToken(mock.Anything, mock.MatchedBy(func(token coreentity.UserActionToken) bool {
				return token.UserID == "user-1" &&
					token.Purpose == coreentity.UserActionTokenPurposePasswordReset &&
					token.TokenHash != "" &&
					token.ExpiresAt.After(time.Now().UTC())
			})).
			Return(nil)
		bus.EXPECT().
			PublishJSON(mock.Anything, common.MessageTopicEmailForgotPassword, mock.MatchedBy(func(payload any) bool {
				msg, ok := payload.(coreentity.QueuedEmailMessage)
				return ok &&
					msg.Email == "user@example.com" &&
					msg.UserName == "User Example" &&
					msg.TenantName == "Acme Corp" &&
					msg.PreferredLanguage == "id" &&
					msg.ActionLink != ""
			})).
			Return(nil)

		core := NewUserCore(Config{
			Repo: repo,
			Tx:   tx,
			Bus:  bus,
		})

		err := core.ForgotPassword(ctx, coreentity.User{
			Email:      "user@example.com",
			TenantCode: "tenant-1",
		})

		require.NoError(t, err)
	})

	t.Run("user not found becomes nil", func(t *testing.T) {
		repo := dbMocks.NewUserRepository(t)

		repo.EXPECT().
			FindActiveUserByEmailAndTenant(mock.Anything, "missing@example.com", "tenant-1").
			Return(nil, errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageUserNotFound))

		core := NewUserCore(Config{
			Repo: repo,
		})

		err := core.ForgotPassword(ctx, coreentity.User{
			Email:      "missing@example.com",
			TenantCode: "tenant-1",
		})

		require.NoError(t, err)
	})

	t.Run("returns dependency error from preferred language lookup", func(t *testing.T) {
		repo := dbMocks.NewUserRepository(t)

		repo.EXPECT().
			FindActiveUserByEmailAndTenant(mock.Anything, "user@example.com", "tenant-1").
			Return(&coreentity.User{
				ID:       "user-1",
				Email:    "user@example.com",
				TenantID: "tenant-1",
			}, nil)
		repo.EXPECT().
			GetTenantPreferredLanguage(mock.Anything, "tenant-1").
			Return("", errmsg.NewCustomErrors(500).SetMessage("database error"))

		core := NewUserCore(Config{
			Repo: repo,
		})

		err := core.ForgotPassword(ctx, coreentity.User{
			Email:      "user@example.com",
			TenantCode: "tenant-1",
		})

		require.Error(t, err)
	})
}
