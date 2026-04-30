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

func TestUpdateOrgUnit_UpdatesNonCompanyOrgUnit(t *testing.T) {
	ctx := context.Background()
	data := coreentity.OrgUnit{ID: "org-1", TenantID: "tenant-1", Name: "Sales"}
	repo := dbmocks.NewOrgUnitRepository(t)
	repo.EXPECT().
		GetOrgUnit(mock.Anything, coreentity.OrgUnit{TenantID: "tenant-1", ID: "org-1"}).
		Return(&coreentity.OrgUnit{ID: "org-1", Category: "department"}, nil)
	repo.EXPECT().UpdateOrgUnit(mock.Anything, data).Return(nil)
	core := NewOrgUnitCore(Config{Repo: repo})

	err := core.UpdateOrgUnit(ctx, data)

	require.NoError(t, err)
}

func TestUpdateOrgUnit_RejectsCompanyOrgUnit(t *testing.T) {
	ctx := context.Background()
	data := coreentity.OrgUnit{ID: "org-1", TenantID: "tenant-1", Name: "Company"}
	repo := dbmocks.NewOrgUnitRepository(t)
	repo.EXPECT().
		GetOrgUnit(mock.Anything, coreentity.OrgUnit{TenantID: "tenant-1", ID: "org-1"}).
		Return(&coreentity.OrgUnit{ID: "org-1", Category: "company"}, nil)
	core := NewOrgUnitCore(Config{Repo: repo})

	err := core.UpdateOrgUnit(ctx, data)

	require.Error(t, err)
}

func TestUpdateOrgUnit_ReturnsRepositoryError(t *testing.T) {
	ctx := context.Background()
	data := coreentity.OrgUnit{ID: "org-1", TenantID: "tenant-1", Name: "Sales"}
	wantErr := errors.New("repository failed")
	repo := dbmocks.NewOrgUnitRepository(t)
	repo.EXPECT().
		GetOrgUnit(mock.Anything, coreentity.OrgUnit{TenantID: "tenant-1", ID: "org-1"}).
		Return(nil, wantErr)
	core := NewOrgUnitCore(Config{Repo: repo})

	err := core.UpdateOrgUnit(ctx, data)

	require.ErrorIs(t, err, wantErr)
}
