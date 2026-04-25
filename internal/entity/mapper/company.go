package mapper

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
)

func CompanyFromCoreToRest(item coreentity.Company) restentity.Company {
	return restentity.Company{
		ID:   item.ID,
		Code: item.Code,
		Name: item.Name,
		Config: restentity.CompanyConfig{
			Logo:               item.Config.Logo,
			PreferredLanguage:  item.Config.PreferredLanguage,
			SupportedLanguages: item.Config.SupportedLanguages,
			Timezone:           item.Config.Timezone,
			DateFormat:         item.Config.DateFormat,
			TimeFormat:         item.Config.TimeFormat,
		},
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func CompanyFromRestCreateToCore(ctx context.Context, req restentity.CreateCompanyReq) coreentity.Company {
	uc := common.GetUserContext(ctx)
	return coreentity.Company{
		UserCtx:  uc,
		TenantID: uc.TenantID,
		Code:     req.Code,
		Name:     req.Name,
		Config: coreentity.CompanyConfig{
			Logo:               req.Config.Logo,
			PreferredLanguage:  req.Config.PreferredLanguage,
			SupportedLanguages: req.Config.SupportedLanguages,
			Timezone:           req.Config.Timezone,
			DateFormat:         req.Config.DateFormat,
			TimeFormat:         req.Config.TimeFormat,
		},
	}
}

func CompanyFromRestUpdateToCore(ctx context.Context, req restentity.UpdateCompanyReq) coreentity.Company {
	uc := common.GetUserContext(ctx)
	return coreentity.Company{
		UserCtx:  uc,
		TenantID: uc.TenantID,
		ID:       req.ID,
		Name:     req.Name,
	}
}
