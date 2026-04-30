package core

import (
	"context"
	"os"
	"testing"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/config"
	dbmocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	config.Envs = &config.Config{}
	config.Envs.App.Name = "test"
	os.Exit(m.Run())
}

func testStringPtr(value string) *string {
	result := value
	return &result
}

func intPtr(value int) *int {
	result := value
	return &result
}

func floatPtr(value float64) *float64 {
	result := value
	return &result
}

func attendanceUserContext(roles ...string) common.UserContext {
	return common.UserContext{
		UserID:   "user-1",
		UserName: "Budi",
		TenantID: "tenant-1",
		Roles:    roles,
	}
}

func attendanceSelfFilter(roles ...string) coreentity.SelfFilter {
	userCtx := attendanceUserContext(roles...)
	return coreentity.SelfFilter{
		UserCtx:  userCtx,
		TenantID: userCtx.TenantID,
		UserID:   userCtx.UserID,
	}
}

func attendanceEmployee() *coreentity.Employee {
	return &coreentity.Employee{
		ID:         "employee-1",
		TenantID:   "tenant-1",
		EmployeeNo: "EMP-001",
		FullName:   "Budi",
		ShiftID:    testStringPtr("shift-1"),
		LocationID: testStringPtr("location-1"),
	}
}

func attendanceShift() *coreentity.WorkShift {
	return &coreentity.WorkShift{
		ID:                 "shift-1",
		TenantID:           "tenant-1",
		Name:               "Morning",
		Timezone:           "UTC",
		StartTime:          "08:00",
		EndTime:            "17:00",
		GracePeriodMinutes: 15,
	}
}

func runTransaction(t *testing.T, tx *dbmocks.Transactor) {
	t.Helper()

	tx.EXPECT().
		WithinTransaction(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
}

func requireDateNearNow(t *testing.T, got string, location *time.Location) {
	t.Helper()

	before := time.Now().In(location).Format("2006-01-02")
	after := time.Now().In(location).Format("2006-01-02")
	require.Contains(t, []string{before, after}, got)
}
