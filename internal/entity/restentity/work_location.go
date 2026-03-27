package restentity

import "codebase-app/pkg/types"

type WorkLocation struct {
	ID           string   `json:"id"`
	OrgUnitID    *string  `json:"org_unit_id"`
	OrgUnitName  *string  `json:"org_unit_name"`
	Name         string   `json:"name"`
	Address      *string  `json:"address"`
	Timezone     string   `json:"timezone"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	RadiusMeters *int     `json:"radius_meters"`
}

type GetWorkLocationsReq struct {
	Q         string  `query:"q" validate:"omitempty,min=2"`
	OrgUnitID *string `query:"org_unit_id" validate:"omitempty,uuidv7"`
	types.MetaQuery
}

func (r *GetWorkLocationsReq) SetDefault() {
	r.MetaQuery.SetDefault()
}

type GetWorkLocationsResp struct {
	Items []WorkLocation `json:"items"`
	Meta  types.Meta     `json:"meta"`
}

type GetWorkLocationReq struct {
	ID string `params:"id" validate:"required"`
}

type GetWorkLocationResp struct {
	WorkLocation
}

type CreateWorkLocationReq struct {
	Name         string   `json:"name" validate:"required,min=2"`
	OrgUnitID    *string  `json:"org_unit_id" validate:"omitempty,uuidv7"`
	Address      *string  `json:"address"`
	Timezone     string   `json:"timezone" validate:"required,timezone"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	RadiusMeters *int     `json:"radius_meters" validate:"omitempty,min=1"`
}

type CreateWorkLocationResp struct {
	ID string `json:"id"`
}

type UpdateWorkLocationReq struct {
	ID           string   `params:"id" validate:"required"`
	Name         string   `json:"name" validate:"required,min=2"`
	OrgUnitID    *string  `json:"org_unit_id" validate:"omitempty,uuidv7"`
	Address      *string  `json:"address"`
	Timezone     string   `json:"timezone" validate:"required,timezone"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	RadiusMeters *int     `json:"radius_meters" validate:"omitempty,min=1"`
}

type DeleteWorkLocationReq struct {
	ID string `params:"id" validate:"required"`
}
