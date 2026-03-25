package restentity

import "codebase-app/pkg/types"

type JobPosition struct {
	ID    string  `json:"id"`
	Name  string `json:"name"`
	Grade *string `json:"grade"`
}

type GetJobPositionsReq struct {
	Q        string `query:"q" validate:"omitempty,min=2"`
	types.MetaQuery
}

func (r *GetJobPositionsReq) SetDefault() {
	r.MetaQuery.SetDefault()
}

type GetJobPositionsResp struct {
	Items []JobPosition `json:"items"`
	Meta  types.Meta    `json:"meta"`
}

type GetJobPositionReq struct {
	ID string `params:"id" validate:"required"`
}

type GetJobPositionResp struct {
	JobPosition
}

type CreateJobPositionReq struct {
	Name  string  `json:"name" validate:"required,min=2"`
	Grade *string `json:"grade"`
}

type CreateJobPositionResp struct {
	ID string `json:"id"`
}

type UpdateJobPositionReq struct {
	ID    string  `params:"id" validate:"required"`
	Name  string `json:"name" validate:"required,min=2"`
	Grade *string `json:"grade"`
}

type DeleteJobPositionReq struct {
	ID string `params:"id" validate:"required"`
}