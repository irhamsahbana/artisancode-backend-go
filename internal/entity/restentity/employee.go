package restentity

import "codebase-app/pkg/types"

type Employee struct {
	ID            string  `json:"id"`
	EmployeeNo    string  `json:"employee_no"`
	FullName      string  `json:"full_name"`
	OrgUnitID     *string `json:"org_unit_id"`
	JobPositionID *string `json:"job_position_id"`
	LocationID    *string `json:"location_id"`
	ShiftID       *string `json:"shift_id"`
	Status        string  `json:"status"`
	JoinDate      *string `json:"join_date"`
}

type GetEmployeesReq struct {
	Q         string  `query:"q" validate:"omitempty,min=2"`
	Status    string  `query:"status" validate:"omitempty,oneof=active inactive"`
	OrgUnitID *string `query:"org_unit_id"`
	types.MetaQuery
}

func (r *GetEmployeesReq) SetDefault() {
	r.MetaQuery.SetDefault()
}

type GetEmployeesResp struct {
	Items []Employee `json:"items"`
	Meta  types.Meta `json:"meta"`
}

type GetEmployeeReq struct {
	ID string `params:"id" validate:"required"`
}

type GetEmployeeResp struct {
	Employee
}

type CreateEmployeeReq struct {
	EmployeeNo    string  `json:"employee_no" validate:"required,min=3,max=50"`
	FullName      string  `json:"full_name" validate:"required,min=3"`
	OrgUnitID     *string `json:"org_unit_id" validate:"omitempty,uuidv7"`
	JobPositionID *string `json:"job_position_id" validate:"omitempty,uuidv7"`
	LocationID    *string `json:"location_id" validate:"omitempty,uuidv7"`
	ShiftID       *string `json:"shift_id" validate:"omitempty,uuidv7"`
	Status        string  `json:"status" validate:"required,oneof=active inactive"`
	JoinDate      *string `json:"join_date" validate:"omitempty,datetime=2006-01-02"`
}

type CreateEmployeeResp struct {
	ID string `json:"id"`
}

type UpdateEmployeeReq struct {
	ID            string  `params:"id" validate:"required"`
	EmployeeNo    string  `json:"employee_no" validate:"required,min=3,max=50"`
	FullName      string  `json:"full_name" validate:"required,min=3"`
	OrgUnitID     *string `json:"org_unit_id" validate:"omitempty,uuidv7"`
	JobPositionID *string `json:"job_position_id" validate:"omitempty,uuidv7"`
	LocationID    *string `json:"location_id" validate:"omitempty,uuidv7"`
	ShiftID       *string `json:"shift_id" validate:"omitempty,uuidv7"`
	Status        string  `json:"status" validate:"required,oneof=active inactive"`
	JoinDate      *string `json:"join_date" validate:"omitempty,datetime=2006-01-02"`
}

type DeleteEmployeeReq struct {
	ID string `params:"id" validate:"required"`
}
