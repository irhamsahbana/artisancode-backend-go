package mapper

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
)

func LoginReqToCore(ctx context.Context, req restentity.LoginReq) coreentity.User {
	uc := common.GetUserContext(ctx)
	return coreentity.User{
		UserCtx:    uc,
		Email:      req.Email,
		Password:   req.Password,
		TenantCode: req.TenantCode,
	}
}

func RegisterReqToCore(ctx context.Context, req restentity.RegisterReq) coreentity.User {
	uc := common.GetUserContext(ctx)
	return coreentity.User{
		UserCtx:    uc,
		Name:       req.Name,
		UserName:   req.UserName,
		Email:      req.Email,
		Password:   req.Password,
		TenantName: req.TenantName,
	}
}

func RegisterReqToTenant(ctx context.Context, req restentity.RegisterReq) coreentity.Tenant {
	uc := common.GetUserContext(ctx)
	preferredLanguage := req.Language
	if preferredLanguage == "" {
		preferredLanguage = "id"
	}
	return coreentity.Tenant{
		UserCtx:           uc,
		Name:              req.TenantName,
		Code:              req.TenantCode,
		PreferredLanguage: preferredLanguage,
	}
}

func RefreshTokenReqToCore(ctx context.Context, req restentity.RefreshTokenReq) coreentity.User {
	uc := common.GetUserContext(ctx)
	return coreentity.User{
		UserCtx:      uc,
		RefreshToken: req.RefreshToken,
	}
}

func AuthTokensToLoginResp(tokens coreentity.AuthTokens) restentity.LoginResp {
	return restentity.LoginResp{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}

func AuthTokensToRegisterResp(tokens coreentity.AuthTokens) restentity.RegisterResp {
	return restentity.RegisterResp{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}

func UserFromCoreToRest(item coreentity.User) restentity.UserResource {
	roles := make([]restentity.UserRole, 0, len(item.RoleIDs))
	for index, roleID := range item.RoleIDs {
		roleName := ""
		if index < len(item.RoleNames) {
			roleName = item.RoleNames[index]
		}
		roles = append(roles, restentity.UserRole{
			ID:   roleID,
			Name: roleName,
		})
	}

	return restentity.UserResource{
		ID:          item.ID,
		Name:        item.Name,
		UserName:    item.UserName,
		Email:       item.Email,
		CompanyID:   item.CompanyID,
		CompanyName: item.CompanyName,
		Roles:       roles,
	}
}

func UserFromRestCreateToCore(ctx context.Context, req restentity.CreateUserReq) coreentity.User {
	uc := common.GetUserContext(ctx)
	return coreentity.User{
		UserCtx:   uc,
		Name:      req.Name,
		UserName:  req.UserName,
		Email:     req.Email,
		Password:  req.Password,
		TenantID:  uc.TenantID,
		RoleIDs:   req.RoleIDs,
		CompanyID: req.CompanyID,
	}
}

func UserFromRestUpdateToCore(ctx context.Context, req restentity.UpdateUserReq) coreentity.User {
	uc := common.GetUserContext(ctx)
	return coreentity.User{
		UserCtx:   uc,
		ID:        req.ID,
		Name:      req.Name,
		UserName:  req.UserName,
		Email:     req.Email,
		Password:  req.Password,
		TenantID:  uc.TenantID,
		RoleIDs:   req.RoleIDs,
		CompanyID: req.CompanyID,
	}
}
