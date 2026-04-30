package core

import (
	"context"
	"errors"
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	infraConfig "codebase-app/internal/infrastructure/config"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func init() {
	infraConfig.Envs = &infraConfig.Config{}
	infraConfig.Envs.App.Name = "internalclient-core-test"
}

func TestGetInternalClients(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.UserContextKeyClaims, common.UserContext{
		Roles: []string{coreentity.InternalUserRoleOperator},
	})
	filter := coreentity.InternalClientListFilter{Q: "acme"}
	want := []coreentity.InternalClient{{ID: "client-1", Name: "Acme"}}

	tests := []struct {
		name      string
		ctx       context.Context
		setup     func(repo *dbMocks.InternalClientRepository)
		want      []coreentity.InternalClient
		wantTotal int
		wantError bool
	}{
		{
			name: "returns clients for authorized operator",
			ctx:  ctx,
			setup: func(repo *dbMocks.InternalClientRepository) {
				repo.EXPECT().
					GetInternalClients(mock.Anything, filter).
					Return(want, 1, nil)
			},
			want:      want,
			wantTotal: 1,
		},
		{
			name:      "rejects unauthorized user before repository call",
			ctx:       context.Background(),
			wantError: true,
		},
		{
			name: "returns repository error",
			ctx:  ctx,
			setup: func(repo *dbMocks.InternalClientRepository) {
				repo.EXPECT().
					GetInternalClients(mock.Anything, filter).
					Return(nil, 0, errors.New("repo failed"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := dbMocks.NewInternalClientRepository(t)
			if tt.setup != nil {
				tt.setup(repo)
			}
			core := NewInternalClientCore(Config{Repo: repo})

			got, total, err := core.GetInternalClients(tt.ctx, filter)

			if tt.wantError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantTotal, total)
		})
	}
}
