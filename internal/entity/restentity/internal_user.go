package restentity

import "codebase-app/pkg/types"

type InternalUserLoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type InternalUserLoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type InternalUserRefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type InternalUserLogoutReq struct {
	RefreshToken string `json:"refresh_token"`
}

type InternalUserResource struct {
	ID          string  `json:"id"`
	FullName    string  `json:"full_name"`
	Email       string  `json:"email"`
	RoleCode    string  `json:"role_code"`
	Status      string  `json:"status"`
	LastLoginAt *string `json:"last_login_at"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type GetInternalUsersReq struct {
	Q        string `query:"q" validate:"omitempty"`
	RoleCode string `query:"role_code" validate:"omitempty,oneof=super_admin operator finance operations"`
	Status   string `query:"status" validate:"omitempty,oneof=invited active inactive"`
	types.MetaQuery
}

func (r *GetInternalUsersReq) SetDefault() {
	r.MetaQuery.SetDefault()
}

type GetInternalUsersResp struct {
	Items []InternalUserResource `json:"items"`
	Meta  types.Meta             `json:"meta"`
}

type GetInternalUserReq struct {
	ID string `params:"id" validate:"required"`
}

type GetInternalUserResp struct {
	InternalUserResource
}

type CreateInternalUserReq struct {
	FullName string `json:"full_name" validate:"required,min=3,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	RoleCode string `json:"role_code" validate:"required,oneof=super_admin operator finance operations"`
	Status   string `json:"status" validate:"required,oneof=invited active inactive"`
}

type CreateInternalUserResp struct {
	ID string `json:"id"`
}

type UpdateInternalUserReq struct {
	ID       string `params:"id" validate:"required"`
	FullName string `json:"full_name" validate:"required,min=3,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"omitempty,min=8"`
	RoleCode string `json:"role_code" validate:"required,oneof=super_admin operator finance operations"`
	Status   string `json:"status" validate:"required,oneof=invited active inactive"`
}

type DeleteInternalUserReq struct {
	ID string `params:"id" validate:"required"`
}
