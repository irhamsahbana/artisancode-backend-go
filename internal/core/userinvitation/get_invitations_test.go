package core

import (
	"context"
	"errors"
	"testing"

	"codebase-app/internal/entity/coreentity"
	dbmocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserInvitationCore_GetInvitations(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		filter    coreentity.UserInvitationListFilter
		setup     func(repo *dbmocks.UserInvitationRepository)
		want      []coreentity.UserInvitation
		wantCount int
		wantError bool
	}{
		{
			name: "success for admin",
			filter: coreentity.UserInvitationListFilter{
				UserCtx: invitationUserContext("admin"),
				Page:    1,
			},
			setup: func(repo *dbmocks.UserInvitationRepository) {
				repo.EXPECT().
					GetInvitations(mock.Anything, coreentity.UserInvitationListFilter{
						UserCtx: invitationUserContext("admin"),
						Page:    1,
					}).
					Return([]coreentity.UserInvitation{{ID: "invitation-1"}}, 1, nil)
			},
			want:      []coreentity.UserInvitation{{ID: "invitation-1"}},
			wantCount: 1,
		},
		{
			name: "rejects unauthorized user",
			filter: coreentity.UserInvitationListFilter{
				UserCtx: invitationUserContext("employee"),
			},
			wantError: true,
		},
		{
			name: "repository error",
			filter: coreentity.UserInvitationListFilter{
				UserCtx: invitationUserContext("owner"),
			},
			setup: func(repo *dbmocks.UserInvitationRepository) {
				repo.EXPECT().
					GetInvitations(mock.Anything, coreentity.UserInvitationListFilter{
						UserCtx: invitationUserContext("owner"),
					}).
					Return(nil, 0, errors.New("repository failed"))
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

			core := NewUserInvitationCore(Config{Repo: repo})
			got, count, err := core.GetInvitations(ctx, tt.filter)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantCount, count)
		})
	}
}
