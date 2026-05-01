package restentity

type LoginReq struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
	TenantCode string `json:"tenant_code" validate:"required,min=2,max=5"`
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

type GoogleLoginReq struct {
	IDToken string `json:"id_token" validate:"required"`
	Nonce   string `json:"nonce" validate:"omitempty"`
}

func (r *GoogleLoginReq) Log() map[string]interface{} {
	return map[string]interface{}{}
}

type GoogleRegisterInitReq struct {
	IDToken string `json:"id_token" validate:"required"`
	Nonce   string `json:"nonce" validate:"omitempty"`
}

func (r *GoogleRegisterInitReq) Log() map[string]interface{} {
	return map[string]interface{}{
		"id_token_present": r.IDToken != "",
		"nonce_present":    r.Nonce != "",
	}
}

type GoogleRegisterInitResp struct {
	RegistrationToken string  `json:"registration_token"`
	Email             string  `json:"email"`
	DisplayName       string  `json:"display_name"`
	PictureURL        *string `json:"picture_url,omitempty"`
}

type GoogleRegisterReq struct {
	IDToken            string `json:"id_token" validate:"omitempty"`
	RegistrationToken  string `json:"registration_token" validate:"omitempty"`
	Nonce              string `json:"nonce" validate:"omitempty"`
	TenantName         string `json:"tenant_name" validate:"required"`
	TenantCode         string `json:"tenant_code" validate:"required"`
	ConfirmTenantSetup bool   `json:"confirm_tenant_setup"`
	Language           string `json:"language" validate:"omitempty,oneof=id en"`
}

func (r *GoogleRegisterReq) Log() map[string]interface{} {
	return map[string]interface{}{
		"tenant_name":          r.TenantName,
		"tenant_code":          r.TenantCode,
		"confirm_tenant_setup": r.ConfirmTenantSetup,
		"language":             r.Language,
		"id_token_present":     r.IDToken != "",
		"registration_token":   r.RegistrationToken != "",
		"nonce_present":        r.Nonce != "",
	}
}

type GoogleRegisterResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TenantCode   string `json:"tenant_code"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LogoutReq struct {
	RefreshToken string `json:"refresh_token"`
}

type RegisterReq struct {
	Name       string `json:"name" validate:"required,min=3"`
	UserName   string `json:"username" validate:"required,min=3"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
	TenantCode string `json:"tenant_code" validate:"required,min=2,max=5,uppercase,alphanum"`
	TenantName string `json:"tenant_name" validate:"required"`
	Language   string `json:"language" validate:"omitempty,oneof=id en"`
}

type RegisterResp struct {
	Email                string `json:"email"`
	VerificationRequired bool   `json:"verification_required"`
}

type TenantProfileResp struct {
	TenantID            string `json:"tenant_id"`
	TenantName          string `json:"tenant_name"`
	TenantCode          string `json:"tenant_code"`
	CanChangeTenantCode bool   `json:"can_change_tenant_code"`
}

type VerifyEmailReq struct {
	Token string `json:"token" validate:"required"`
}

type ResendVerificationEmailReq struct {
	Email      string `json:"email" validate:"required,email"`
	TenantCode string `json:"tenant_code" validate:"required,min=2,max=5,uppercase,alphanum"`
}

type ForgotPasswordReq struct {
	Email      string `json:"email" validate:"required,email"`
	TenantCode string `json:"tenant_code" validate:"required,min=2,max=5,uppercase,alphanum"`
}

type ResetPasswordReq struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}
