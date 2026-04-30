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

func TestDeleteJobPosition_ReturnsNilWhenRepositorySucceeds(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.JobPositionDeleteFilter{TenantID: "tenant-1", ID: "job-1"}
	repo := dbmocks.NewJobPositionRepository(t)
	repo.EXPECT().DeleteJobPosition(mock.Anything, filter).Return(nil)
	core := NewJobPositionCore(Config{Repo: repo})

	err := core.DeleteJobPosition(ctx, filter)

	require.NoError(t, err)
}

func TestDeleteJobPosition_ReturnsRepositoryError(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.JobPositionDeleteFilter{TenantID: "tenant-1", ID: "job-1"}
	wantErr := errors.New("repository failed")
	repo := dbmocks.NewJobPositionRepository(t)
	repo.EXPECT().DeleteJobPosition(mock.Anything, filter).Return(wantErr)
	core := NewJobPositionCore(Config{Repo: repo})

	err := core.DeleteJobPosition(ctx, filter)

	require.ErrorIs(t, err, wantErr)
}
