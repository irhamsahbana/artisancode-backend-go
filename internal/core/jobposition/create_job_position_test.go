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

func TestCreateJobPosition_ReturnsCreatedJobPosition(t *testing.T) {
	ctx := context.Background()
	data := coreentity.JobPosition{TenantID: "tenant-1", Name: "Engineer"}
	want := &coreentity.JobPosition{ID: "job-1", TenantID: "tenant-1", Name: "Engineer"}
	repo := dbmocks.NewJobPositionRepository(t)
	repo.EXPECT().CreateJobPosition(mock.Anything, data).Return(want, nil)
	core := NewJobPositionCore(Config{Repo: repo})

	got, err := core.CreateJobPosition(ctx, data)

	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestCreateJobPosition_ReturnsRepositoryError(t *testing.T) {
	ctx := context.Background()
	data := coreentity.JobPosition{TenantID: "tenant-1", Name: "Engineer"}
	wantErr := errors.New("repository failed")
	repo := dbmocks.NewJobPositionRepository(t)
	repo.EXPECT().CreateJobPosition(mock.Anything, data).Return(nil, wantErr)
	core := NewJobPositionCore(Config{Repo: repo})

	got, err := core.CreateJobPosition(ctx, data)

	require.ErrorIs(t, err, wantErr)
	require.Nil(t, got)
}
