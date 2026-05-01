package mapper

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
)

func InternalUserLoginReqToCore(ctx context.Context, req restentity.InternalUserLoginReq) coreentity.InternalUser {
	uc := common.GetUserContext(ctx)
	return coreentity.InternalUser{
		UserCtx:  uc,
		Email:    req.Email,
		Password: req.Password,
	}
}

func InternalUserRefreshTokenReqToCore(ctx context.Context, req restentity.InternalUserRefreshTokenReq) coreentity.InternalUser {
	uc := common.GetUserContext(ctx)
	return coreentity.InternalUser{
		UserCtx:      uc,
		RefreshToken: req.RefreshToken,
	}
}

func InternalUserLogoutReqToCore(ctx context.Context, req restentity.InternalUserLogoutReq) coreentity.InternalUser {
	uc := common.GetUserContext(ctx)
	return coreentity.InternalUser{
		UserCtx:      uc,
		RefreshToken: req.RefreshToken,
	}
}

func AuthTokensToInternalUserLoginResp(tokens coreentity.AuthTokens) restentity.InternalUserLoginResp {
	return restentity.InternalUserLoginResp{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}

func InternalUserFromCoreToRest(item coreentity.InternalUser) restentity.InternalUserResource {
	return restentity.InternalUserResource{
		ID:          item.ID,
		FullName:    item.FullName,
		Email:       item.Email,
		RoleCode:    item.RoleCode,
		Status:      item.Status,
		LastLoginAt: item.LastLoginAt,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func InternalUserFromRestCreateToCore(ctx context.Context, req restentity.CreateInternalUserReq) coreentity.InternalUser {
	uc := common.GetUserContext(ctx)
	return coreentity.InternalUser{
		UserCtx:  uc,
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
		RoleCode: req.RoleCode,
		Status:   req.Status,
	}
}

func InternalUserFromRestUpdateToCore(ctx context.Context, req restentity.UpdateInternalUserReq) coreentity.InternalUser {
	uc := common.GetUserContext(ctx)
	return coreentity.InternalUser{
		UserCtx:  uc,
		ID:       req.ID,
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
		RoleCode: req.RoleCode,
		Status:   req.Status,
	}
}
