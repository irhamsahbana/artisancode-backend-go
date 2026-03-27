package restentity

import "codebase-app/pkg/types"

type PaginationResp struct {
	Total    int `json:"total"`
	Page     int `json:"page"`
	PerPage  int `json:"per_page"`
	LastPage int `json:"last_page"`
}

func NewPaginationResp(meta types.Meta) PaginationResp {
	return PaginationResp{
		Total:    meta.TotalData,
		Page:     meta.Page,
		PerPage:  meta.Paginate,
		LastPage: meta.TotalPage,
	}
}

type Role struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
}

type Permission struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// --- Roles request/response ---

type GetRolesReq struct {
	Q     string `query:"q" validate:"omitempty,min=2"`
	Limit int    `query:"limit" validate:"omitempty,min=1"`
	Page  int    `query:"page" validate:"omitempty,min=1"`
}

func (r *GetRolesReq) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}
	if r.Limit < 1 {
		r.Limit = 25
	}
}

type GetRolesResp struct {
	Items      []RoleWithPermissions `json:"items"`
	Pagination PaginationResp        `json:"pagination"`
}

type RoleWithPermissions struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
}

type GetRoleReq struct {
	ID string `params:"id" validate:"required"`
}

type CreateRoleReq struct {
	Name        string   `json:"name" validate:"required,min=2"`
	Permissions []string `json:"permissions"`
}

type CreateRoleResp struct {
	ID string `json:"id"`
}

type UpdateRoleReq struct {
	ID            string   `params:"id" validate:"required"`
	Name          string   `json:"name" validate:"required,min=2"`
	PermissionIDs []string `json:"permission_ids"`
}

type DeleteRoleReq struct {
	ID string `params:"id" validate:"required"`
}

// --- Permissions request/response ---

type GetPermissionsReq struct {
	Q     string `query:"q" validate:"omitempty,min=2"`
	Limit int    `query:"limit" validate:"omitempty,min=1"`
	Page  int    `query:"page" validate:"omitempty,min=1"`
}

func (r *GetPermissionsReq) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}
	if r.Limit < 1 {
		r.Limit = 25
	}
}

type GetPermissionsResp struct {
	Items      []Permission   `json:"items"`
	Pagination PaginationResp `json:"pagination"`
}
