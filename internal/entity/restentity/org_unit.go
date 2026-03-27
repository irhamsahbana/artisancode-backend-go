package restentity

import "codebase-app/pkg/types"

type OrgUnit struct {
	ID       string  `json:"id"`
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	ParentID *string `json:"parent_id"`
	Category string  `json:"category"`
}

type GetOrgUnitsReq struct {
	Q         string `query:"q" validate:"omitempty,min=2"`
	Category  string `query:"category" validate:"omitempty,oneof=company branch division department unit"`
	types.MetaQuery
}

func (r *GetOrgUnitsReq) SetDefault() {
	r.MetaQuery.SetDefault()
}

type GetOrgUnitsResp struct {
	Items []OrgUnit  `json:"items"`
	Meta  types.Meta `json:"meta"`
}

type GetOrgUnitReq struct {
	ID string `params:"id" validate:"required"`
}

type GetOrgUnitResp struct {
	OrgUnit
}

type CreateOrgUnitReq struct {
	Code      string  `json:"code" validate:"required,min=1"`
	Name      string  `json:"name" validate:"required,min=2"`
	ParentID  *string `json:"parent_id" validate:"omitempty,uuidv7"`
	Category  string  `json:"category" validate:"required,oneof=company branch division department unit"`
}

type CreateOrgUnitResp struct {
	ID string `json:"id"`
}

type UpdateOrgUnitReq struct {
	ID        string  `params:"id" validate:"required"`
	Code      string  `json:"code" validate:"required,min=1"`
	Name      string  `json:"name" validate:"required,min=2"`
	ParentID  *string `json:"parent_id" validate:"omitempty,uuidv7"`
	Category  string  `json:"category" validate:"required,oneof=company branch division department unit"`
}

type DeleteOrgUnitReq struct {
	ID string `params:"id" validate:"required"`
}