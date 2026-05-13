package mapper

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
)

func LoginReqToCore(_ context.Context, req restentity.LoginReq) coreentity.User {
	return coreentity.User{
		Email:      req.Email,
		Password:   req.Password,
		TenantCode: req.TenantCode,
	}
}

func RegisterReqToCore(_ context.Context, req restentity.RegisterReq) coreentity.User {
	return coreentity.User{
		Name:       req.Name,
		UserName:   req.UserName,
		Email:      req.Email,
		Password:   req.Password,
		TenantName: req.TenantName,
	}
}

func RegisterReqToTenant(_ context.Context, req restentity.RegisterReq) coreentity.Tenant {
	preferredLanguage := req.Language
	if preferredLanguage == "" {
		preferredLanguage = "id"
	}
	return coreentity.Tenant{
		Name:              req.TenantName,
		Code:              req.TenantCode,
		PreferredLanguage: preferredLanguage,
	}
}

func RefreshTokenReqToCore(_ context.Context, req restentity.RefreshTokenReq) coreentity.User {
	return coreentity.User{
		RefreshToken: req.RefreshToken,
	}
}

func LogoutReqToCore(_ context.Context, req restentity.LogoutReq) coreentity.User {
	return coreentity.User{
		RefreshToken: req.RefreshToken,
	}
}

func VerifyEmailReqToCore(_ context.Context, req restentity.VerifyEmailReq) coreentity.UserActionToken {
	return coreentity.UserActionToken{
		Token:   req.Token,
		Purpose: coreentity.UserActionTokenPurposeEmailVerification,
	}
}

func ResendVerificationEmailReqToCore(_ context.Context, req restentity.ResendVerificationEmailReq) coreentity.User {
	return coreentity.User{
		Email:      req.Email,
		TenantCode: req.TenantCode,
	}
}

func ForgotPasswordReqToCore(_ context.Context, req restentity.ForgotPasswordReq) coreentity.User {
	return coreentity.User{
		Email:      req.Email,
		TenantCode: req.TenantCode,
	}
}

func ResetPasswordReqToCore(_ context.Context, req restentity.ResetPasswordReq) (coreentity.UserActionToken, coreentity.User) {
	return coreentity.UserActionToken{
			Token:   req.Token,
			Purpose: coreentity.UserActionTokenPurposePasswordReset,
		}, coreentity.User{
			Password: req.Password,
		}
}

func AuthTokensToLoginResp(tokens coreentity.AuthTokens) restentity.LoginResp {
	return restentity.LoginResp{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}

func GoogleRegisterReqToCore(_ context.Context, req restentity.GoogleRegisterReq) coreentity.GoogleRegisterInput {
	preferredLanguage := req.Language
	if preferredLanguage == "" {
		preferredLanguage = "id"
	}
	return coreentity.GoogleRegisterInput{
		IDToken:            req.IDToken,
		RegistrationToken:  req.RegistrationToken,
		Nonce:              req.Nonce,
		TenantName:         req.TenantName,
		TenantCode:         req.TenantCode,
		ConfirmTenantSetup: req.ConfirmTenantSetup,
		PreferredLanguage:  preferredLanguage,
	}
}

func GoogleLoginReqToCore(_ context.Context, req restentity.GoogleLoginReq) coreentity.GoogleLoginInput {
	return coreentity.GoogleLoginInput{
		IDToken: req.IDToken,
		Nonce:   req.Nonce,
	}
}

func GoogleRegisterInitReqToCore(
	_ context.Context,
	req restentity.GoogleRegisterInitReq,
) coreentity.GoogleRegisterInitInput {
	return coreentity.GoogleRegisterInitInput{
		IDToken: req.IDToken,
		Nonce:   req.Nonce,
	}
}

func GoogleRegisterResultToResp(result coreentity.GoogleRegisterResult) restentity.GoogleRegisterResp {
	return restentity.GoogleRegisterResp{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		TenantCode:   result.TenantCode,
	}
}

func GoogleRegisterInitResultToResp(
	result coreentity.GoogleRegisterInitResult,
) restentity.GoogleRegisterInitResp {
	return restentity.GoogleRegisterInitResp{
		RegistrationToken: result.RegistrationToken,
		Email:             result.Email,
		DisplayName:       result.DisplayName,
		PictureURL:        result.PictureURL,
	}
}

func TenantProfileToResp(profile coreentity.TenantProfile) restentity.TenantProfileResp {
	return restentity.TenantProfileResp{
		TenantID:            profile.ID,
		TenantName:          profile.Name,
		TenantCode:          profile.Code,
		CanChangeTenantCode: profile.CanChangeTenantCode,
	}
}

func RegisterResultToRegisterResp(result coreentity.RegisterResult) restentity.RegisterResp {
	return restentity.RegisterResp{
		Email:                result.Email,
		VerificationRequired: result.VerificationRequired,
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
