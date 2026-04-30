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

func TestMessageQueueCore_CleanupProcessedMessages(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		req       coreentity.CleanupProcessedMessageQueueReq
		setup     func(repo *dbmocks.MessageQueueRepository)
		want      *coreentity.CleanupProcessedMessageQueueResp
		wantError bool
	}{
		{
			name: "success with defaults",
			req:  coreentity.CleanupProcessedMessageQueueReq{},
			setup: func(repo *dbmocks.MessageQueueRepository) {
				repo.EXPECT().
					CleanupProcessedMessages(
						mock.Anything,
						mock.MatchedBy(func(req coreentity.CleanupProcessedMessageQueueReq) bool {
							return req.Before != "" && req.Limit == 100
						}),
					).
					Return(&coreentity.CleanupProcessedMessageQueueResp{Deleted: 3}, nil)
			},
			want: &coreentity.CleanupProcessedMessageQueueResp{Deleted: 3},
		},
		{
			name: "repository error",
			req: coreentity.CleanupProcessedMessageQueueReq{
				Before: "2026-04-30T00:00:00Z",
				Limit:  10,
			},
			setup: func(repo *dbmocks.MessageQueueRepository) {
				repo.EXPECT().
					CleanupProcessedMessages(
						mock.Anything,
						coreentity.CleanupProcessedMessageQueueReq{
							Before: "2026-04-30T00:00:00Z",
							Limit:  10,
						},
					).
					Return(nil, errors.New("repository failed"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := dbmocks.NewMessageQueueRepository(t)
			tt.setup(repo)

			core := NewMessageQueueCore(MessageQueueCoreConfig{Repo: repo})
			got, err := core.CleanupProcessedMessages(ctx, tt.req)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
