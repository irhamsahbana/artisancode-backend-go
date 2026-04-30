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

func TestGetOrgUnits_ReturnsRepositoryResult(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.OrgUnitListFilter{TenantID: "tenant-1", Page: 1, Paginate: 10}
	items := []coreentity.OrgUnit{{ID: "org-1", TenantID: "tenant-1", Name: "HQ"}}
	repo := dbmocks.NewOrgUnitRepository(t)
	repo.EXPECT().GetOrgUnits(mock.Anything, filter).Return(items, 1, nil)
	core := NewOrgUnitCore(Config{Repo: repo})

	got, total, err := core.GetOrgUnits(ctx, filter)

	require.NoError(t, err)
	require.Equal(t, items, got)
	require.Equal(t, 1, total)
}

func TestGetOrgUnits_ReturnsRepositoryError(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.OrgUnitListFilter{TenantID: "tenant-1"}
	wantErr := errors.New("repository failed")
	repo := dbmocks.NewOrgUnitRepository(t)
	repo.EXPECT().GetOrgUnits(mock.Anything, filter).Return(nil, 0, wantErr)
	core := NewOrgUnitCore(Config{Repo: repo})

	got, total, err := core.GetOrgUnits(ctx, filter)

	require.ErrorIs(t, err, wantErr)
	require.Nil(t, got)
	require.Zero(t, total)
}
