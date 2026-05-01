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

func TestUserCore_ResendVerificationEmail(t *testing.T) {
	ctx := context.Background()

	t.Run("success issues verification token and queues email", func(t *testing.T) {
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
			Return("en", nil)
		tx.EXPECT().
			WithinTransaction(mock.Anything, mock.Anything).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})
		repo.EXPECT().
			DeleteUserActionTokensByPurpose(mock.Anything, "user-1", coreentity.UserActionTokenPurposeEmailVerification).
			Return(nil)
		repo.EXPECT().
			CreateUserActionToken(mock.Anything, mock.MatchedBy(func(token coreentity.UserActionToken) bool {
				return token.UserID == "user-1" &&
					token.Purpose == coreentity.UserActionTokenPurposeEmailVerification &&
					token.TokenHash != "" &&
					token.ExpiresAt.After(time.Now().UTC())
			})).
			Return(nil)
		bus.EXPECT().
			PublishJSON(mock.Anything, common.MessageTopicEmailVerification, mock.MatchedBy(func(payload any) bool {
				msg, ok := payload.(coreentity.QueuedEmailMessage)
				return ok &&
					msg.Email == "user@example.com" &&
					msg.UserName == "User Example" &&
					msg.TenantName == "Acme Corp" &&
					msg.PreferredLanguage == "en" &&
					msg.ActionLink != ""
			})).
			Return(nil)

		core := NewUserCore(Config{
			Repo: repo,
			Tx:   tx,
			Bus:  bus,
		})

		err := core.ResendVerificationEmail(ctx, coreentity.User{
			Email:      "user@example.com",
			TenantCode: "tenant-1",
		})

		require.NoError(t, err)
	})

	t.Run("rejects already verified email", func(t *testing.T) {
		repo := dbMocks.NewUserRepository(t)

		repo.EXPECT().
			FindActiveUserByEmailAndTenant(mock.Anything, "user@example.com", "tenant-1").
			Return(&coreentity.User{
				ID:              "user-1",
				Email:           "user@example.com",
				TenantID:        "tenant-1",
				EmailVerifiedAt: timePtr("2026-04-30T09:00:00Z"),
			}, nil)

		core := NewUserCore(Config{
			Repo: repo,
		})

		err := core.ResendVerificationEmail(ctx, coreentity.User{
			Email:      "user@example.com",
			TenantCode: "tenant-1",
		})

		require.Error(t, err)

		customErr, ok := err.(*errmsg.CustomError)
		require.True(t, ok)
		require.Equal(t, errmsg.MessageEmailIsAlreadyVerified, customErr.Msg)
	})

	t.Run("returns dependency error from user lookup", func(t *testing.T) {
		repo := dbMocks.NewUserRepository(t)

		repo.EXPECT().
			FindActiveUserByEmailAndTenant(mock.Anything, "user@example.com", "tenant-1").
			Return(nil, errmsg.NewCustomErrors(500).SetMessage("database error"))

		core := NewUserCore(Config{
			Repo: repo,
		})

		err := core.ResendVerificationEmail(ctx, coreentity.User{
			Email:      "user@example.com",
			TenantCode: "tenant-1",
		})

		require.Error(t, err)
	})
}
