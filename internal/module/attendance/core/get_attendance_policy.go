package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *attendanceCore) GetAttendancePolicy(ctx context.Context, filter coreentity.SelfFilter) (*coreentity.AttendancePolicy, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetAttendancePolicy")
	defer span.End()

	company, err := c.getPolicyCompany(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &coreentity.AttendancePolicy{
		UserCtx:                 filter.UserCtx,
		Timezone:                company.Config.Timezone,
		AttendanceRadiusMeters:  company.Config.AttendanceRadiusMeters,
		AttendanceCheckInStart:  company.Config.AttendanceCheckInStart,
		AttendanceCheckInEnd:    company.Config.AttendanceCheckInEnd,
		AttendanceCheckOutStart: company.Config.AttendanceCheckOutStart,
		AttendanceCheckOutEnd:   company.Config.AttendanceCheckOutEnd,
	}, nil
}

func (c *attendanceCore) getPolicyCompany(ctx context.Context, filter coreentity.SelfFilter) (*coreentity.Company, error) {
	if filter.UserCtx.CompanyID != nil {
		return c.companyRepo.GetCompany(ctx, coreentity.Company{
			TenantID: filter.TenantID,
			ID:       *filter.UserCtx.CompanyID,
		})
	}

	if !filter.UserCtx.HasRole("owner") {
		company, err := c.repo.GetCompanyByUserID(ctx, filter.TenantID, filter.UserID)
		if err != nil {
			return nil, err
		}
		if company != nil {
			return company, nil
		}
	}

	companies, _, err := c.companyRepo.GetCompanies(ctx, coreentity.CompanyListFilter{
		TenantID: filter.TenantID,
		Page:     1,
		Paginate: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(companies) == 0 {
		return nil, errmsg.NewCustomErrors(404).SetMessage("Company policy not found")
	}

	return &companies[0], nil
}
