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

func TestGetAllOrgUnitsByCompany_ReturnsRepositoryResult(t *testing.T) {
	ctx := context.Background()
	items := []coreentity.OrgUnit{{ID: "org-1", TenantID: "tenant-1", Name: "HQ"}}
	repo := dbmocks.NewOrgUnitRepository(t)
	repo.EXPECT().
		GetAllOrgUnitsByCompany(mock.Anything, "tenant-1", "company-1").
		Return(items, nil)
	core := NewOrgUnitCore(Config{Repo: repo})

	got, err := core.GetAllOrgUnitsByCompany(ctx, "tenant-1", "company-1")

	require.NoError(t, err)
	require.Equal(t, items, got)
}

func TestGetAllOrgUnitsByCompany_ReturnsRepositoryError(t *testing.T) {
	ctx := context.Background()
	wantErr := errors.New("repository failed")
	repo := dbmocks.NewOrgUnitRepository(t)
	repo.EXPECT().
		GetAllOrgUnitsByCompany(mock.Anything, "tenant-1", "company-1").
		Return(nil, wantErr)
	core := NewOrgUnitCore(Config{Repo: repo})

	got, err := core.GetAllOrgUnitsByCompany(ctx, "tenant-1", "company-1")

	require.ErrorIs(t, err, wantErr)
	require.Nil(t, got)
}
