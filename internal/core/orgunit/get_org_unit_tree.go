package core

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"context"
)

func (c *orgUnitCore) GetOrgUnitTree(ctx context.Context, tenantID string, companyID string) ([]coreentity.OrgUnitTreeNode, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:orgunit:get_org_unit_tree:GetOrgUnitTree")
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
	// Build lookup maps: node data and parent->children relationships
	type nodeData struct {
		ID       string
		Code     string
		Name     string
		Category string
	}
	nodeMap := make(map[string]nodeData)
	parentMap := make(map[string][]string) // parentID -> childIDs

	for _, item := range items {
		nodeMap[item.ID] = nodeData{
			ID:       item.ID,
			Code:     item.Code,
			Name:     item.Name,
			Category: item.Category,
		}
		if item.ParentID != nil {
			parentMap[*item.ParentID] = append(parentMap[*item.ParentID], item.ID)
		}
	}

	// Recursively build children for a given parent ID
	var buildChildren func(parentID string) []coreentity.OrgUnitTreeNode
	buildChildren = func(parentID string) []coreentity.OrgUnitTreeNode {
		childIDs := parentMap[parentID]
		if len(childIDs) == 0 {
			return []coreentity.OrgUnitTreeNode{}
		}
		children := make([]coreentity.OrgUnitTreeNode, 0, len(childIDs))
		for _, childID := range childIDs {
			nd := nodeMap[childID]
			children = append(children, coreentity.OrgUnitTreeNode{
				ID:       nd.ID,
				Code:     nd.Code,
				Name:     nd.Name,
				Category: nd.Category,
				Children: buildChildren(childID),
			})
		}
		return children
	}

	// Find root nodes and build their full subtrees
	roots := make([]coreentity.OrgUnitTreeNode, 0)
	for _, item := range items {
		if item.ParentID == nil {
			nd := nodeMap[item.ID]
			roots = append(roots, coreentity.OrgUnitTreeNode{
				ID:       nd.ID,
				Code:     nd.Code,
				Name:     nd.Name,
				Category: nd.Category,
				Children: buildChildren(item.ID),
			})
		}
	}

	return roots
}
