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

func TestGetJobPosition_ReturnsRepositoryResult(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.JobPosition{TenantID: "tenant-1", ID: "job-1"}
	want := &coreentity.JobPosition{ID: "job-1", TenantID: "tenant-1", Name: "Engineer"}
	repo := dbmocks.NewJobPositionRepository(t)
	repo.EXPECT().GetJobPosition(mock.Anything, filter).Return(want, nil)
	core := NewJobPositionCore(Config{Repo: repo})

	got, err := core.GetJobPosition(ctx, filter)

	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestGetJobPosition_ReturnsRepositoryError(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.JobPosition{TenantID: "tenant-1", ID: "job-1"}
	wantErr := errors.New("repository failed")
	repo := dbmocks.NewJobPositionRepository(t)
	repo.EXPECT().GetJobPosition(mock.Anything, filter).Return(nil, wantErr)
	core := NewJobPositionCore(Config{Repo: repo})

	got, err := core.GetJobPosition(ctx, filter)

	require.ErrorIs(t, err, wantErr)
	require.Nil(t, got)
}
