package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

// allowedChildCategories defines which categories can be created under a given parent category.
var allowedChildCategories = map[string][]string{
	"company":    {"branch"},
	"branch":     {"division"},
	"division":   {"department", "division"},
	"department": {"unit", "department"},
	"unit":       {"unit"},
}

func (c *orgUnitCore) CreateOrgUnit(ctx context.Context, data coreentity.OrgUnit) (*coreentity.OrgUnit, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:orgunit:create_org_unit:CreateOrgUnit")
	defer span.End()

	// 1. Check code uniqueness within tenant
	exists, err := c.repo.ExistsOrgUnitByCode(ctx, data.TenantID, data.Code)
	if err != nil {
		return nil, err
	}
	if exists {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{
			"tenant_id": data.TenantID,
			"code":      data.Code,
		}).Msg("Org unit code already exists")
		return nil, errmsg.NewCustomErrors(400).SetMessage("Organization unit code already exists")
	}

	// 2. Validate parent exists and check category hierarchy
	if data.ParentID != nil {
		parent, err := c.repo.GetOrgUnit(ctx, coreentity.OrgUnit{
			TenantID: data.TenantID,
			ID:       *data.ParentID,
		})
		if err != nil {
			return nil, err
		}

		// 3. Validate category is allowed under parent
		if !isCategoryAllowed(parent.Category, data.Category) {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{
				"parent_category": parent.Category,
				"child_category":  data.Category,
			}).Msg("Invalid category hierarchy")
			return nil, errmsg.NewCustomErrors(400).SetMessage("Invalid category for the selected parent")
		}
	}

	return c.repo.CreateOrgUnit(ctx, data)
}

// isCategoryAllowed checks if a child category is allowed under a parent category.
func isCategoryAllowed(parentCategory, childCategory string) bool {
	allowed, ok := allowedChildCategories[parentCategory]
	if !ok {
		return false
	}
	for _, cat := range allowed {
		if cat == childCategory {
			return true
		}
	}
	return false
}
