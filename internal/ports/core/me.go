package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type MeCore interface {
	GetMe(ctx context.Context, filter coreentity.SelfFilter) (*coreentity.Me, error)
	GetMyEmployee(ctx context.Context, filter coreentity.SelfFilter) (*coreentity.Employee, error)
	GetMyShiftToday(ctx context.Context, filter coreentity.SelfFilter) (*coreentity.WorkShiftToday, error)
}
