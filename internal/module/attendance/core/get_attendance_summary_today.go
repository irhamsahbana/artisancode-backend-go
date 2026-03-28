package core

import (
	"context"
	"time"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *attendanceCore) GetAttendanceSummaryToday(ctx context.Context, filter coreentity.SelfFilter) (*coreentity.AttendanceSummary, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetAttendanceSummaryToday")
	defer span.End()

	employee, err := c.repo.GetEmployeeByUserID(ctx, filter.TenantID, filter.UserID)
	if err != nil {
		return nil, err
	}

	policy, err := c.GetAttendancePolicy(ctx, filter)
	if err != nil {
		return nil, err
	}

	today := time.Now().Format("2006-01-02")
	if policy.Timezone != "" {
		location, loadErr := time.LoadLocation(policy.Timezone)
		if loadErr == nil {
			today = time.Now().In(location).Format("2006-01-02")
		}
	}

	logs, _, err := c.repo.GetAttendanceLogs(ctx, coreentity.AttendanceLogListFilter{
		UserCtx:       filter.UserCtx,
		TenantID:      filter.TenantID,
		EmployeeID:    &employee.ID,
		AttendanceDay: &today,
		Page:          1,
		Paginate:      10,
	})
	if err != nil {
		return nil, err
	}

	summary := coreentity.AttendanceSummary{
		UserCtx:        filter.UserCtx,
		AttendanceDate: today,
		CheckedIn:      false,
		CheckedOut:     false,
		CanCheckIn:     true,
		CanCheckOut:    false,
	}

	for _, item := range logs {
		if summary.LastLogType == nil {
			summary.LastLogType = stringPtr(item.Type)
			summary.LastLoggedAt = stringPtr(item.LoggedAt)
		}

		if item.Type == "check_in" && summary.CheckInLogID == nil {
			summary.CheckedIn = true
			summary.CheckInLogID = stringPtr(item.ID)
		}

		if item.Type == "check_out" && summary.CheckOutLogID == nil {
			summary.CheckedOut = true
			summary.CheckOutLogID = stringPtr(item.ID)
		}
	}

	summary.CanCheckIn = !summary.CheckedIn
	summary.CanCheckOut = summary.CheckedIn && !summary.CheckedOut

	return &summary, nil
}

func stringPtr(value string) *string {
	result := value
	return &result
}
