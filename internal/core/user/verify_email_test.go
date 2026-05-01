package core

import (
	"context"
	"errors"
	"testing"

	"codebase-app/internal/entity/coreentity"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserCore_VerifyEmail(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func(repo *dbMocks.UserRepository, tx *dbMocks.Transactor)
		wantErr error
	}{
		{
			name: "success marks email verified and token used",
			setup: func(repo *dbMocks.UserRepository, tx *dbMocks.Transactor) {
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				repo.EXPECT().
					GetValidUserActionToken(
						mock.Anything,
						hashUserActionToken("verify-token"),
						coreentity.UserActionTokenPurposeEmailVerification,
					).
					Return(&coreentity.UserActionToken{
						ID:     "token-1",
						UserID: "user-1",
					}, nil)
				repo.EXPECT().
					MarkUserEmailVerified(mock.Anything, "user-1").
					Return(nil)
				repo.EXPECT().
					MarkUserActionTokenUsed(mock.Anything, "token-1").
					Return(nil)
			},
		},
		{
			name: "returns dependency error when mark email verified fails",
			setup: func(repo *dbMocks.UserRepository, tx *dbMocks.Transactor) {
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.Anything).
					RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					})
				repo.EXPECT().
					GetValidUserActionToken(
						mock.Anything,
						hashUserActionToken("verify-token"),
						coreentity.UserActionTokenPurposeEmailVerification,
					).
					Return(&coreentity.UserActionToken{
						ID:     "token-1",
						UserID: "user-1",
					}, nil)
				repo.EXPECT().
					MarkUserEmailVerified(mock.Anything, "user-1").
					Return(errors.New("mark email verified failed"))
			},
			wantErr: errors.New("mark email verified failed"),
		},
		{
			name: "returns transaction error when commit fails",
			setup: func(repo *dbMocks.UserRepository, tx *dbMocks.Transactor) {
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
						hashUserActionToken("verify-token"),
						coreentity.UserActionTokenPurposeEmailVerification,
					).
					Return(&coreentity.UserActionToken{
						ID:     "token-1",
						UserID: "user-1",
					}, nil)
				repo.EXPECT().
					MarkUserEmailVerified(mock.Anything, "user-1").
					Return(nil)
				repo.EXPECT().
					MarkUserActionTokenUsed(mock.Anything, "token-1").
					Return(nil)
			},
			wantErr: errors.New("commit failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := dbMocks.NewUserRepository(t)
			tx := dbMocks.NewTransactor(t)

			tt.setup(repo, tx)

			core := NewUserCore(Config{
				Repo: repo,
				Tx:   tx,
			})

			err := core.VerifyEmail(ctx, coreentity.UserActionToken{
				Token: "verify-token",
			})

			if tt.wantErr != nil {
				require.EqualError(t, err, tt.wantErr.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}
