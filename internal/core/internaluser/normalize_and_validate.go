package core

import (
	"strings"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/pkg/errmsg"

	"golang.org/x/crypto/bcrypt"
)

func (c *internalUserCore) normalizeAndValidate(data *coreentity.InternalUser, passwordRequired bool) error {
	data.FullName = strings.TrimSpace(data.FullName)
	data.Email = strings.ToLower(strings.TrimSpace(data.Email))
	data.RoleCode = strings.TrimSpace(data.RoleCode)
	data.Status = strings.TrimSpace(data.Status)

	if data.RoleCode != coreentity.InternalUserRoleSuperAdmin && data.RoleCode != coreentity.InternalUserRoleOperator {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageInternalUserRoleIsInvalid)
	}
	if data.Status != coreentity.InternalUserStatusInvited &&
		data.Status != coreentity.InternalUserStatusActive &&
		data.Status != coreentity.InternalUserStatusInactive {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageInternalUserStatusIsInvalid)
	}
	if passwordRequired && strings.TrimSpace(data.Password) == "" {
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessagePasswordIsRequired)
	}

	return nil
}

func hashInternalPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}
