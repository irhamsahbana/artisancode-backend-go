package restentity

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	TenantID string `json:"tenant_id" validate:"required,uuidv7"`
}

func (r *LoginReq) Log() map[string]interface{} {
	return map[string]interface{}{
		"email": r.Email,
	}
}

type LoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RegisterReq struct {
	Name       string `json:"name" validate:"required,min=3"`
	UserName   string `json:"username" validate:"required,min=3"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
	TenantCode string `json:"tenant_code" validate:"required,min=3,max=64"`
	TenantName string `json:"tenant_name" validate:"required"`
}

type RegisterResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
