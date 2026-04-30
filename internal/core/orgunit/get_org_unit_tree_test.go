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

func TestGetOrgUnitTree_ReturnsNestedTree(t *testing.T) {
	ctx := context.Background()
	companyID := "company-1"
	branchID := "branch-1"
	items := []coreentity.OrgUnit{
		{ID: companyID, Code: "CO", Name: "Company", Category: "company"},
		{ID: branchID, Code: "BR", Name: "Branch", Category: "branch", ParentID: &companyID},
		{ID: "dept-1", Code: "DP", Name: "Department", Category: "department", ParentID: &branchID},
	}
	repo := dbmocks.NewOrgUnitRepository(t)
	repo.EXPECT().
		GetAllOrgUnitsByCompany(mock.Anything, "tenant-1", "company-1").
		Return(items, nil)
	core := NewOrgUnitCore(Config{Repo: repo})

	got, err := core.GetOrgUnitTree(ctx, "tenant-1", "company-1")

	require.NoError(t, err)
	require.Equal(t, []coreentity.OrgUnitTreeNode{
		{
			ID:       "company-1",
			Code:     "CO",
			Name:     "Company",
			Category: "company",
			Children: []coreentity.OrgUnitTreeNode{
				{
					ID:       "branch-1",
					Code:     "BR",
					Name:     "Branch",
					Category: "branch",
					Children: []coreentity.OrgUnitTreeNode{
						{
							ID:       "dept-1",
							Code:     "DP",
							Name:     "Department",
							Category: "department",
							Children: []coreentity.OrgUnitTreeNode{},
						},
					},
				},
			},
		},
	}, got)
}

func TestGetOrgUnitTree_ReturnsRepositoryError(t *testing.T) {
	ctx := context.Background()
	wantErr := errors.New("repository failed")
	repo := dbmocks.NewOrgUnitRepository(t)
	repo.EXPECT().
		GetAllOrgUnitsByCompany(mock.Anything, "tenant-1", "company-1").
		Return(nil, wantErr)
	core := NewOrgUnitCore(Config{Repo: repo})

	got, err := core.GetOrgUnitTree(ctx, "tenant-1", "company-1")

	require.ErrorIs(t, err, wantErr)
	require.Nil(t, got)
}
