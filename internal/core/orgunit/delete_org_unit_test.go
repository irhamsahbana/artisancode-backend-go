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

func TestDeleteOrgUnit_DeletesNonCompanyOrgUnit(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.OrgUnitDeleteFilter{TenantID: "tenant-1", ID: "org-1"}
	repo := dbmocks.NewOrgUnitRepository(t)
	repo.EXPECT().
		GetOrgUnit(mock.Anything, coreentity.OrgUnit{TenantID: "tenant-1", ID: "org-1"}).
		Return(&coreentity.OrgUnit{ID: "org-1", Category: "department"}, nil)
	repo.EXPECT().DeleteOrgUnit(mock.Anything, filter).Return(nil)
	core := NewOrgUnitCore(Config{Repo: repo})

	err := core.DeleteOrgUnit(ctx, filter)

	require.NoError(t, err)
}

func TestDeleteOrgUnit_RejectsCompanyOrgUnit(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.OrgUnitDeleteFilter{TenantID: "tenant-1", ID: "org-1"}
	repo := dbmocks.NewOrgUnitRepository(t)
	repo.EXPECT().
		GetOrgUnit(mock.Anything, coreentity.OrgUnit{TenantID: "tenant-1", ID: "org-1"}).
		Return(&coreentity.OrgUnit{ID: "org-1", Category: "company"}, nil)
	core := NewOrgUnitCore(Config{Repo: repo})

	err := core.DeleteOrgUnit(ctx, filter)

	require.Error(t, err)
}

func TestDeleteOrgUnit_ReturnsRepositoryError(t *testing.T) {
	ctx := context.Background()
	filter := coreentity.OrgUnitDeleteFilter{TenantID: "tenant-1", ID: "org-1"}
	wantErr := errors.New("repository failed")
	repo := dbmocks.NewOrgUnitRepository(t)
	repo.EXPECT().
		GetOrgUnit(mock.Anything, coreentity.OrgUnit{TenantID: "tenant-1", ID: "org-1"}).
		Return(nil, wantErr)
	core := NewOrgUnitCore(Config{Repo: repo})

	err := core.DeleteOrgUnit(ctx, filter)

	require.ErrorIs(t, err, wantErr)
}
