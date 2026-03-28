package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type MeRepository interface {
	GetEmployeeByUserID(ctx context.Context, tenantID, userID string) (*coreentity.Employee, error)
	GetWorkShiftByUserID(ctx context.Context, tenantID, userID string) (*coreentity.WorkShift, error)
}
