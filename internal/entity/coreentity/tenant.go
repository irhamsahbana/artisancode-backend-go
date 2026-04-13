package coreentity

import "codebase-app/internal/entity/common"

type Tenant struct {
	UserCtx common.UserContext

	ID                string
	Name              string
	Code              string
	PreferredLanguage string
}
