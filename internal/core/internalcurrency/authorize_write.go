package core

import (
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/pkg/errmsg"
)

func (c *internalCurrencyCore) authorizeWrite(userCtx common.UserContext) error {
	if !userCtx.HasRole(coreentity.InternalUserRoleSuperAdmin) &&
		!userCtx.HasRole(coreentity.InternalUserRoleOperator) {
		return errmsg.NewCustomErrors(403).SetMessage(errmsg.MessageYouAreNotAuthorizedToManageInternalProducts)
	}

	return nil
}
