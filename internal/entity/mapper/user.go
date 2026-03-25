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
		UserCtx:  uc,
		Email:    req.Email,
		Password: req.Password,
		TenantID: req.TenantID,
	}
}

func RegisterReqToCore(ctx context.Context, req restentity.RegisterReq) coreentity.User {
	uc := common.GetUserContext(ctx)
	return coreentity.User{
		UserCtx:     uc,
		Name:       req.Name,
		UserName:   req.UserName,
		Email:      req.Email,
		Password:   req.Password,
		TenantName: req.TenantName,
	}
}

func RegisterReqToTenant(ctx context.Context, req restentity.RegisterReq) coreentity.Tenant {
	uc := common.GetUserContext(ctx)
	return coreentity.Tenant{
		UserCtx: uc,
		Name:    req.TenantName,
		Code:    req.TenantCode,
	}
}

func RefreshTokenReqToCore(ctx context.Context, req restentity.RefreshTokenReq) coreentity.User {
	uc := common.GetUserContext(ctx)
	return coreentity.User{
		UserCtx:       uc,
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
