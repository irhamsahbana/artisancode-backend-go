package core

import (
	"context"
	"testing"

	"codebase-app/internal/entity/coreentity"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"
	"codebase-app/pkg/errmsg"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserCore_UpdateUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		input     coreentity.User
		setup     func(tx *dbMocks.Transactor)
		wantError bool
	}{
		{
			name: "success",
			input: coreentity.User{
				TenantID: "tenant-1",
				ID:       "user-1",
				UserName: "Updated User",
			},
			setup: func(tx *dbMocks.Transactor) {
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
					Return(nil)
			},
		},
		{
			name: "dependency error from transaction",
			input: coreentity.User{
				TenantID: "tenant-1",
				ID:       "user-1",
				UserName: "Updated User",
			},
			setup: func(tx *dbMocks.Transactor) {
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
					Return(errmsg.NewCustomErrors(500).SetMessage("transaction error"))
			},
			wantError: true,
		},
		{
			name: "success with password change",
			input: coreentity.User{
				TenantID: "tenant-1",
				ID:       "user-1",
				UserName: "Updated User",
				Password: "newpassword123",
			},
			setup: func(tx *dbMocks.Transactor) {
				tx.EXPECT().
					WithinTransaction(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
					Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := dbMocks.NewUserRepository(t)
			tx := dbMocks.NewTransactor(t)

			if tt.setup != nil {
				tt.setup(tx)
			}

			core := NewUserCore(Config{
				Repo: repo,
				Tx:   tx,
			})

			err := core.UpdateUser(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
