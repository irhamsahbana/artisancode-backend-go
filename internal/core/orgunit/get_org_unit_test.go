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

func TestGetOrgUnit_ReturnsRepositoryResult(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.OrgUnit{TenantID: "tenant-1", ID: "org-1"}
	want := &coreentity.OrgUnit{ID: "org-1", TenantID: "tenant-1", Name: "HQ"}
	repo := dbmocks.NewOrgUnitRepository(t)
	repo.EXPECT().GetOrgUnit(mock.Anything, filter).Return(want, nil)
	core := NewOrgUnitCore(Config{Repo: repo})

	got, err := core.GetOrgUnit(ctx, filter)

	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestGetOrgUnit_ReturnsRepositoryError(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.OrgUnit{TenantID: "tenant-1", ID: "org-1"}
	wantErr := errors.New("repository failed")
	repo := dbmocks.NewOrgUnitRepository(t)
	repo.EXPECT().GetOrgUnit(mock.Anything, filter).Return(nil, wantErr)
	core := NewOrgUnitCore(Config{Repo: repo})

	got, err := core.GetOrgUnit(ctx, filter)

	require.ErrorIs(t, err, wantErr)
	require.Nil(t, got)
}
