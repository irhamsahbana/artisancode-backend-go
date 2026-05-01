package core

import (
	"context"
	"errors"
	"testing"
	"time"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/integration/tokencache"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUserCore_ResetPassword(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		setup      func(repo *dbMocks.UserRepository, tx *dbMocks.Transactor, cache *tokencache.TokenCache)
		wantErr    error
		assertPost func(t *testing.T, cache *tokencache.TokenCache)
	}{
		{
			name: "success updates password marks token used and clears refresh tokens",
			setup: func(repo *dbMocks.UserRepository, tx *dbMocks.Transactor, cache *tokencache.TokenCache) {
				cache.SetRefreshToken(ctx, "reset-user-token", tokencache.RefreshTokenData{
					UserID:   "user-1",
					TenantID: "tenant-1",
				}, time.Hour)
				cache.SetRefreshToken(ctx, "other-user-token", tokencache.RefreshTokenData{
					UserID:   "user-2",
					TenantID: "tenant-2",
				}, time.Hour)

				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				repo.EXPECT().
					GetValidUserActionToken(
						mock.Anything,
						hashUserActionToken("reset-token"),
						coreentity.UserActionTokenPurposePasswordReset,
					).
					Return(&coreentity.UserActionToken{
						ID:       "token-1",
						UserID:   "user-1",
						TenantID: "tenant-1",
					}, nil)
				repo.EXPECT().
					UpdateUserPassword(
						mock.Anything,
						"user-1",
						"tenant-1",
						mock.MatchedBy(func(hashed string) bool {
							return bcrypt.CompareHashAndPassword([]byte(hashed), []byte("new-password")) == nil
						}),
					).
					Return(nil)
				repo.EXPECT().
					MarkUserActionTokenUsed(mock.Anything, "token-1").
					Return(nil)
			},
			assertPost: func(t *testing.T, cache *tokencache.TokenCache) {
				_, found := cache.GetRefreshToken(ctx, "reset-user-token")
				require.False(t, found)

				other, found := cache.GetRefreshToken(ctx, "other-user-token")
				require.True(t, found)
				require.Equal(t, "user-2", other.UserID)
			},
		},
		{
			name: "returns dependency error when token lookup fails",
			setup: func(repo *dbMocks.UserRepository, tx *dbMocks.Transactor, cache *tokencache.TokenCache) {
				cache.SetRefreshToken(ctx, "reset-user-token", tokencache.RefreshTokenData{
					UserID:   "user-1",
					TenantID: "tenant-1",
				}, time.Hour)

				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				repo.EXPECT().
					GetValidUserActionToken(
						mock.Anything,
						hashUserActionToken("reset-token"),
						coreentity.UserActionTokenPurposePasswordReset,
					).
					Return(nil, errors.New("token lookup failed"))
			},
			wantErr: errors.New("token lookup failed"),
			assertPost: func(t *testing.T, cache *tokencache.TokenCache) {
				_, found := cache.GetRefreshToken(ctx, "reset-user-token")
				require.True(t, found)
			},
		},
		{
			name: "returns transaction error and keeps refresh tokens when commit fails",
			setup: func(repo *dbMocks.UserRepository, tx *dbMocks.Transactor, cache *tokencache.TokenCache) {
				cache.SetRefreshToken(ctx, "reset-user-token", tokencache.RefreshTokenData{
					UserID:   "user-1",
					TenantID: "tenant-1",
				}, time.Hour)

				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						if err := fn(ctx); err != nil {
							return err
						}
						return errors.New("commit failed")
					})
				repo.EXPECT().
					GetValidUserActionToken(
						mock.Anything,
						hashUserActionToken("reset-token"),
						coreentity.UserActionTokenPurposePasswordReset,
					).
					Return(&coreentity.UserActionToken{
						ID:       "token-1",
						UserID:   "user-1",
						TenantID: "tenant-1",
					}, nil)
				repo.EXPECT().
					UpdateUserPassword(
						mock.Anything,
						"user-1",
						"tenant-1",
						mock.Anything,
					).
					Return(nil)
				repo.EXPECT().
					MarkUserActionTokenUsed(mock.Anything, "token-1").
					Return(nil)
			},
			wantErr: errors.New("commit failed"),
			assertPost: func(t *testing.T, cache *tokencache.TokenCache) {
				_, found := cache.GetRefreshToken(ctx, "reset-user-token")
				require.True(t, found)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := dbMocks.NewUserRepository(t)
			tx := dbMocks.NewTransactor(t)
			cache := tokencache.NewTokenCache(time.Hour, time.Minute)

			tt.setup(repo, tx, cache)

			core := NewUserCore(Config{
				Repo:       repo,
				Tx:         tx,
				TokenCache: cache,
			})

			err := core.ResetPassword(ctx, coreentity.UserActionToken{
				Token: "reset-token",
			}, coreentity.User{
				Password: "new-password",
			})

			if tt.wantErr != nil {
				require.EqualError(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}

			if tt.assertPost != nil {
				tt.assertPost(t, cache)
			}
		})
	}
}
