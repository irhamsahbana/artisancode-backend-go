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

func TestCreateOrgUnit_CreatesOrgUnitWithoutParent(t *testing.T) {
	ctx := context.Background()
	data := coreentity.OrgUnit{
		TenantID: "tenant-1",
		Code:     "BR-1",
		Name:     "Branch 1",
		Category: "branch",
	}
	want := &coreentity.OrgUnit{ID: "org-1", TenantID: "tenant-1", Code: "BR-1"}
	repo := dbmocks.NewOrgUnitRepository(t)
	repo.EXPECT().ExistsOrgUnitByCode(mock.Anything, "tenant-1", "BR-1").Return(false, nil)
	repo.EXPECT().CreateOrgUnit(mock.Anything, data).Return(want, nil)
	core := NewOrgUnitCore(Config{Repo: repo})

	got, err := core.CreateOrgUnit(ctx, data)

	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestCreateOrgUnit_CreatesOrgUnitWhenParentCategoryAllowsChild(t *testing.T) {
	ctx := context.Background()
	parentID := "parent-1"
	data := coreentity.OrgUnit{
		TenantID: "tenant-1",
		Code:     "DIV-1",
		Name:     "Division 1",
		Category: "division",
		ParentID: &parentID,
	}
	want := &coreentity.OrgUnit{ID: "org-1", TenantID: "tenant-1", Code: "DIV-1"}
	repo := dbmocks.NewOrgUnitRepository(t)
	repo.EXPECT().ExistsOrgUnitByCode(mock.Anything, "tenant-1", "DIV-1").Return(false, nil)
	repo.EXPECT().
		GetOrgUnit(mock.Anything, coreentity.OrgUnit{TenantID: "tenant-1", ID: "parent-1"}).
		Return(&coreentity.OrgUnit{ID: "parent-1", Category: "branch"}, nil)
	repo.EXPECT().CreateOrgUnit(mock.Anything, data).Return(want, nil)
	core := NewOrgUnitCore(Config{Repo: repo})

	got, err := core.CreateOrgUnit(ctx, data)

	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestCreateOrgUnit_RejectsDuplicateCode(t *testing.T) {
	ctx := context.Background()
	data := coreentity.OrgUnit{TenantID: "tenant-1", Code: "BR-1", Category: "branch"}
	repo := dbmocks.NewOrgUnitRepository(t)
	repo.EXPECT().ExistsOrgUnitByCode(mock.Anything, "tenant-1", "BR-1").Return(true, nil)
	core := NewOrgUnitCore(Config{Repo: repo})

	got, err := core.CreateOrgUnit(ctx, data)

	require.Error(t, err)
	require.Nil(t, got)
}

func TestCreateOrgUnit_RejectsInvalidParentCategory(t *testing.T) {
	ctx := context.Background()
	parentID := "parent-1"
	data := coreentity.OrgUnit{
		TenantID: "tenant-1",
		Code:     "DEPT-1",
		Category: "department",
		ParentID: &parentID,
	}
	repo := dbmocks.NewOrgUnitRepository(t)
	repo.EXPECT().ExistsOrgUnitByCode(mock.Anything, "tenant-1", "DEPT-1").Return(false, nil)
	repo.EXPECT().
		GetOrgUnit(mock.Anything, coreentity.OrgUnit{TenantID: "tenant-1", ID: "parent-1"}).
		Return(&coreentity.OrgUnit{ID: "parent-1", Category: "branch"}, nil)
	core := NewOrgUnitCore(Config{Repo: repo})

	got, err := core.CreateOrgUnit(ctx, data)

	require.Error(t, err)
	require.Nil(t, got)
}

func TestCreateOrgUnit_ReturnsRepositoryError(t *testing.T) {
	ctx := context.Background()
	data := coreentity.OrgUnit{TenantID: "tenant-1", Code: "BR-1", Category: "branch"}
	wantErr := errors.New("repository failed")
	repo := dbmocks.NewOrgUnitRepository(t)
	repo.EXPECT().ExistsOrgUnitByCode(mock.Anything, "tenant-1", "BR-1").Return(false, wantErr)
	core := NewOrgUnitCore(Config{Repo: repo})

	got, err := core.CreateOrgUnit(ctx, data)

	require.ErrorIs(t, err, wantErr)
	require.Nil(t, got)
}
