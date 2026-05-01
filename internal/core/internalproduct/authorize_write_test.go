package core

import (
	"testing"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/pkg/errmsg"

	"github.com/stretchr/testify/require"
)

func TestAuthorizeWrite(t *testing.T) {
	tests := []struct {
		name        string
		userCtx     common.UserContext
		wantError   bool
		wantCode    int
		wantMessage string
	}{
		{
			name: "allows super admin",
			userCtx: common.UserContext{
				Roles: []string{coreentity.InternalUserRoleSuperAdmin},
			},
		},
		{
			name: "allows operator",
			userCtx: common.UserContext{
				Roles: []string{coreentity.InternalUserRoleOperator},
			},
		},
		{
			name: "rejects unauthorized role",
			userCtx: common.UserContext{
				Roles: []string{"viewer"},
			},
			wantError:   true,
			wantCode:    403,
			wantMessage: errmsg.MessageYouAreNotAuthorizedToManageInternalProducts,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core := NewInternalProductCore(Config{})

			err := core.authorizeWrite(tt.userCtx)

			if !tt.wantError {
				require.NoError(t, err)
				return
			}

			requireHelperError(t, err, tt.wantCode, tt.wantMessage)
		})
	}
}

func requireHelperError(t *testing.T, err error, wantCode int, wantMessage string) {
	t.Helper()

	require.Error(t, err)

	var customErr *errmsg.CustomError
	require.ErrorAs(t, err, &customErr)
	require.Equal(t, wantCode, customErr.Code)
	require.Equal(t, wantMessage, customErr.Msg)
}
