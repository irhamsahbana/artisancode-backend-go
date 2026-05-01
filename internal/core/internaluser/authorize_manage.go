package core

import (
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/pkg/errmsg"
)

func (c *internalUserCore) authorizeManage(userCtx common.UserContext) error {
	if !userCtx.HasRole(coreentity.InternalUserRoleSuperAdmin) {
		return errmsg.NewCustomErrors(403).SetMessage(errmsg.MessageYouAreNotAuthorizedToManageInternalUsers)
	}

	return nil
}
