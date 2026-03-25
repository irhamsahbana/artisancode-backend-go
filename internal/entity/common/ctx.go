package common

import "context"

type UserContextKey string

const (
	UserContextKeyClaims UserContextKey = "claims"
)

type UserContext struct {
	UserID      string
	UserName    string
	TenantID    string
	TenantName  string
	Roles       []string
	CompanyID   *string
	CompanyName *string
}

func GetUserContext(ctx context.Context) UserContext {
	if uc, ok := ctx.Value(UserContextKeyClaims).(UserContext); ok {
		return uc
	}
	return UserContext{}
}

func (uc UserContext) IsOwner() bool {
	return uc.CompanyID == nil
}

func (uc UserContext) HasRole(role string) bool {
	for _, r := range uc.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func (uc UserContext) CanAccessCompany(companyID string) bool {
	if uc.IsOwner() {
		return true
	}
	return uc.CompanyID != nil && *uc.CompanyID == companyID
}
