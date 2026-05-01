package coreentity

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

const AuthProviderGoogle = "google"

type GoogleIdentity struct {
	Subject       string
	Email         string
	EmailVerified bool
	DisplayName   string
	PictureURL    *string
	Nonce         string
}

type UserAuthIdentity struct {
	ID              string
	TenantID        string
	UserID          string
	Provider        string
	ProviderSubject string
	Email           string
	EmailVerified   bool
	DisplayName     string
	PictureURL      *string
}

type GoogleRegisterInput struct {
	IDToken            string
	RegistrationToken  string
	Nonce              string
	TenantName         string
	TenantCode         string
	ConfirmTenantSetup bool
	PreferredLanguage  string
}

type GoogleRegisterResult struct {
	AccessToken  string
	RefreshToken string
	TenantCode   string
}

type GoogleLoginInput struct {
	IDToken string
	Nonce   string
}

type GoogleRegisterInitInput struct {
	IDToken string
	Nonce   string
}

type GoogleRegisterInitResult struct {
	RegistrationToken string
	Email             string
	DisplayName       string
	PictureURL        *string
}

type TenantProfile struct {
	ID                  string
	Name                string
	Code                string
	CanChangeTenantCode bool
}
