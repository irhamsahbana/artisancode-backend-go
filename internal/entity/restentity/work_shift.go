package restentity

import "codebase-app/pkg/types"

type WorkShift struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Timezone           string `json:"timezone"`
	StartTime          string `json:"start_time"`
	EndTime            string `json:"end_time"`
	GracePeriodMinutes int    `json:"grace_period_minutes"`
}

type GetWorkShiftsReq struct {
	Q string `query:"q" validate:"omitempty,min=2"`
	types.MetaQuery
}

func (r *GetWorkShiftsReq) SetDefault() {
	r.MetaQuery.SetDefault()
}

type GetWorkShiftsResp struct {
	Items []WorkShift `json:"items"`
	Meta  types.Meta  `json:"meta"`
}

type GetWorkShiftReq struct {
	ID string `params:"id" validate:"required"`
}

type GetWorkShiftResp struct {
	WorkShift
}

type CreateWorkShiftReq struct {
	Name               string `json:"name" validate:"required,min=2"`
	Timezone           string `json:"timezone" validate:"required,timezone"`
	StartTime          string `json:"start_time" validate:"required,datetime=15:04"`
	EndTime            string `json:"end_time" validate:"required,datetime=15:04"`
	GracePeriodMinutes int    `json:"grace_period_minutes" validate:"required,min=0"`
}

type CreateWorkShiftResp struct {
	ID string `json:"id"`
}

type UpdateWorkShiftReq struct {
	ID                 string `params:"id" validate:"required"`
	Name               string `json:"name" validate:"required,min=2"`
	Timezone           string `json:"timezone" validate:"required,timezone"`
	StartTime          string `json:"start_time" validate:"required,datetime=15:04"`
	EndTime            string `json:"end_time" validate:"required,datetime=15:04"`
	GracePeriodMinutes int    `json:"grace_period_minutes" validate:"required,min=0"`
}

type DeleteWorkShiftReq struct {
	ID string `params:"id" validate:"required"`
}
