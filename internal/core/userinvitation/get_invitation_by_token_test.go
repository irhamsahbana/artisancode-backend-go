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
)

func TestUserInvitationCore_GetInvitationByToken(t *testing.T) {
	ctx := context.Background()
	token := "raw-token"

	tests := []struct {
		name      string
		setup     func(repo *dbmocks.UserInvitationRepository)
		wantError bool
	}{
		{
			name: "success",
			setup: func(repo *dbmocks.UserInvitationRepository) {
				repo.EXPECT().
					GetInvitationByTokenHash(mock.Anything, hashInvitationToken(token)).
					Return(&coreentity.UserInvitation{
						ID:        "invitation-1",
						Status:    coreentity.UserInvitationStatusPending,
						ExpiresAt: time.Now().Add(time.Hour),
					}, nil)
			},
		},
		{
			name: "dependency error from repository",
			setup: func(repo *dbmocks.UserInvitationRepository) {
				repo.EXPECT().
					GetInvitationByTokenHash(mock.Anything, hashInvitationToken(token)).
					Return(nil, errors.New("repository failed"))
			},
			wantError: true,
		},
		{
			name: "business error when invitation is accepted",
			setup: func(repo *dbmocks.UserInvitationRepository) {
				repo.EXPECT().
					GetInvitationByTokenHash(mock.Anything, hashInvitationToken(token)).
					Return(&coreentity.UserInvitation{
						ID:        "invitation-1",
						Status:    coreentity.UserInvitationStatusAccepted,
						ExpiresAt: time.Now().Add(time.Hour),
					}, nil)
			},
			wantError: true,
		},
		{
			name: "business error when invitation is expired",
			setup: func(repo *dbmocks.UserInvitationRepository) {
				repo.EXPECT().
					GetInvitationByTokenHash(mock.Anything, hashInvitationToken(token)).
					Return(&coreentity.UserInvitation{
						ID:        "invitation-1",
						Status:    coreentity.UserInvitationStatusPending,
						ExpiresAt: time.Now().Add(-time.Hour),
					}, nil)
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

			got, err := core.GetInvitationByToken(ctx, token)

			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, "invitation-1", got.ID)
		})
	}
}
