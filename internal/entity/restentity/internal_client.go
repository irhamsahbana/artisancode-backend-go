package restentity

import "codebase-app/pkg/types"

type InternalClientResource struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	OwnerNames  string `json:"owner_names"`
	OwnerEmails string `json:"owner_emails"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type GetInternalClientsReq struct {
	Q string `query:"q" validate:"omitempty,min=2"`
	types.MetaQuery
}

func (r *GetInternalClientsReq) SetDefault() {
	r.MetaQuery.SetDefault()
}

type GetInternalClientsResp struct {
	Items []InternalClientResource `json:"items"`
	Meta  types.Meta               `json:"meta"`
}

type GetInternalClientOwnerPermissionsReq struct {
	ID string `params:"id" validate:"required"`
}

type GetInternalClientOwnerPermissionsResp struct {
	ClientID           string       `json:"client_id"`
	Available          []Permission `json:"available_permissions"`
	OwnerPermissionIDs []string     `json:"owner_permission_ids"`
}

type UpdateInternalClientOwnerPermissionsReq struct {
	ID            string   `params:"id" validate:"required"`
	PermissionIDs []string `json:"permission_ids"`
}
