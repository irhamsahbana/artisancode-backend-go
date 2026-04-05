package restentity

type UserRole struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserResource struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	UserName    string     `json:"username"`
	Email       string     `json:"email"`
	CompanyID   *string    `json:"company_id"`
	CompanyName *string    `json:"company_name"`
	Roles       []UserRole `json:"roles"`
}

type GetUsersReq struct {
	Q     string `query:"q" validate:"omitempty,min=2"`
	Limit int    `query:"limit" validate:"omitempty,min=1"`
	Page  int    `query:"page" validate:"omitempty,min=1"`
}

func (r *GetUsersReq) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}
	if r.Limit < 1 {
		r.Limit = 25
	}
}

type GetUsersResp struct {
	Items      []UserResource `json:"items"`
	Pagination PaginationResp `json:"pagination"`
}

type GetUserReq struct {
	ID string `params:"id" validate:"required,uuidv7"`
}

type CreateUserReq struct {
	Name      string   `json:"name" validate:"required,min=3"`
	UserName  string   `json:"username" validate:"required,min=3"`
	Email     string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"required,min=8"`
	RoleIDs   []string `json:"role_ids" validate:"required,min=1,dive,uuidv7"`
	CompanyID *string  `json:"company_id" validate:"omitempty,uuidv7"`
}

type CreateUserResp struct {
	ID string `json:"id"`
}

type UpdateUserReq struct {
	ID        string   `params:"id" validate:"required,uuidv7"`
	Name      string   `json:"name" validate:"required,min=3"`
	UserName  string   `json:"username" validate:"required,min=3"`
	Email     string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"omitempty,min=8"`
	RoleIDs   []string `json:"role_ids" validate:"required,min=1,dive,uuidv7"`
	CompanyID *string  `json:"company_id" validate:"omitempty,uuidv7"`
}

type DeleteUserReq struct {
	ID string `params:"id" validate:"required,uuidv7"`
}
