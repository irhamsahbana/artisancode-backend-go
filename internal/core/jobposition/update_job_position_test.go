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

func TestUpdateJobPosition_ReturnsNilWhenRepositorySucceeds(t *testing.T) {
	ctx := context.Background()
	data := coreentity.JobPosition{ID: "job-1", TenantID: "tenant-1", Name: "Lead Engineer"}
	repo := dbmocks.NewJobPositionRepository(t)
	repo.EXPECT().UpdateJobPosition(mock.Anything, data).Return(nil)
	core := NewJobPositionCore(Config{Repo: repo})

	err := core.UpdateJobPosition(ctx, data)

	require.NoError(t, err)
}

func TestUpdateJobPosition_ReturnsRepositoryError(t *testing.T) {
	ctx := context.Background()
	data := coreentity.JobPosition{ID: "job-1", TenantID: "tenant-1", Name: "Lead Engineer"}
	wantErr := errors.New("repository failed")
	repo := dbmocks.NewJobPositionRepository(t)
	repo.EXPECT().UpdateJobPosition(mock.Anything, data).Return(wantErr)
	core := NewJobPositionCore(Config{Repo: repo})

	err := core.UpdateJobPosition(ctx, data)

	require.ErrorIs(t, err, wantErr)
}
