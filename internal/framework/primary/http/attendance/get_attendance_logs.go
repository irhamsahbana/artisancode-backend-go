package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"
	"codebase-app/pkg/types"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *attendanceHandler) getAttendanceLogs(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:attendance:get_attendance_logs:getAttendanceLogs")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx = c.UserContext()
		req = new(restentity.GetAttendanceLogsReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse query params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate query params")
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	filter := coreentity.AttendanceLogListFilter{
		UserCtx:        common.GetUserContext(ctx),
		TenantID:       common.GetUserContext(ctx).TenantID,
		EmployeeID:     req.EmployeeID,
		Q:              req.Q,
		Type:           mapper.AttendanceTypePtrFromString(req.Type),
		Source:         mapper.AttendanceSourcePtrFromString(req.Source),
		Status:         mapper.AttendanceStatusPtrFromString(req.Status),
		SelfieStatus:   stringPtrFromValue(req.SelfieStatus),
		OrgUnitID:      req.OrgUnitID,
		BranchID:       req.BranchID,
		WorkLocationID: req.WorkLocationID,
		ExceptionType:  stringPtrFromValue(req.ExceptionType),
		AttendanceDay:  req.AttendanceDay,
		DateFrom:       req.DateFrom,
		DateTo:         req.DateTo,
		Page:           req.Page,
		Paginate:       req.Limit,
	}

	items, total, err := h.core.GetAttendanceLogs(ctx, filter)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get attendance logs")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	restItems := make([]restentity.AttendanceLog, 0, len(items))
	for _, item := range items {
		restItems = append(restItems, mapper.AttendanceLogFromCoreToRest(item))
	}

	meta := types.Meta{
		Page:      req.Page,
		Paginate:  req.Limit,
		TotalData: total,
	}
	meta.CountTotalPage(req.Page, req.Limit, total)

	return c.Status(fiber.StatusOK).JSON(response.Success(restentity.GetAttendanceLogsResp{
		Items:      restItems,
		Pagination: restentity.NewPaginationResp(meta),
	}, ""))
}

func stringPtrFromValue(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}
