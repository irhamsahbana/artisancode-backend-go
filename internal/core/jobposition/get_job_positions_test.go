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

func TestGetJobPositions_ReturnsRepositoryResult(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.JobPositionListFilter{TenantID: "tenant-1", Page: 1, Paginate: 10}
	items := []coreentity.JobPosition{{ID: "job-1", TenantID: "tenant-1", Name: "Engineer"}}
	repo := dbmocks.NewJobPositionRepository(t)
	repo.EXPECT().GetJobPositions(mock.Anything, filter).Return(items, 1, nil)
	core := NewJobPositionCore(Config{Repo: repo})

	got, total, err := core.GetJobPositions(ctx, filter)

	require.NoError(t, err)
	require.Equal(t, items, got)
	require.Equal(t, 1, total)
}

func TestGetJobPositions_ReturnsRepositoryError(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.JobPositionListFilter{TenantID: "tenant-1"}
	wantErr := errors.New("repository failed")
	repo := dbmocks.NewJobPositionRepository(t)
	repo.EXPECT().GetJobPositions(mock.Anything, filter).Return(nil, 0, wantErr)
	core := NewJobPositionCore(Config{Repo: repo})

	got, total, err := core.GetJobPositions(ctx, filter)

	require.ErrorIs(t, err, wantErr)
	require.Nil(t, got)
	require.Zero(t, total)
}
