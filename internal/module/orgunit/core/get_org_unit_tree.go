package core

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"context"
)

func (c *orgUnitCore) GetOrgUnitTree(ctx context.Context, tenantID string, companyID string) ([]coreentity.OrgUnitTreeNode, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetOrgUnitTree")
	defer span.End()

	items, err := c.repo.GetAllOrgUnitsByCompany(ctx, tenantID, companyID)
	if err != nil {
		return nil, err
	}

	// Build tree structure
	return buildOrgUnitTree(items), nil
}

// buildOrgUnitTree converts flat list of org units into recursive tree structure
func buildOrgUnitTree(items []coreentity.OrgUnit) []coreentity.OrgUnitTreeNode {
	nodeMap := make(map[string]*coreentity.OrgUnitTreeNode)
	roots := make([]coreentity.OrgUnitTreeNode, 0)

	// First pass: create all nodes
	for i := range items {
		item := items[i]
		nodeMap[item.ID] = &coreentity.OrgUnitTreeNode{
			ID:       item.ID,
			Name:     item.Name,
			Category: item.Category,
			Children: make([]coreentity.OrgUnitTreeNode, 0),
		}
	}

	// Second pass: build hierarchy
	for i := range items {
		item := items[i]
		node := nodeMap[item.ID]

		if item.ParentID == nil {
			roots = append(roots, *node)
		} else {
			if parent, exists := nodeMap[*item.ParentID]; exists {
				parent.Children = append(parent.Children, *node)
			}
		}
	}

	return roots
}
