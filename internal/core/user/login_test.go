package core

import (
	"context"
	"testing"
	"time"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/integration/tokencache"
	dbMocks "codebase-app/internal/ports/secondary/db/mocks"
	"codebase-app/pkg/errmsg"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserCore_Login(t *testing.T) {
	ctx := context.Background()
	hashedPassword, _ := hashPassword("password123")

	tests := []struct {
		name      string
		input     coreentity.User
		setup     func(repo *dbMocks.UserRepository)
		wantError bool
	}{
		{
			name: "success",
			input: coreentity.User{
				Email:      "test@example.com",
				Password:   "password123",
				TenantCode: "tenant-1",
			},
			setup: func(repo *dbMocks.UserRepository) {
				repo.EXPECT().
					FindActiveUserByEmailAndTenant(mock.Anything, "test@example.com", "tenant-1").
					Return(&coreentity.User{
						ID:              "user-1",
						Email:           "test@example.com",
						Password:        hashedPassword,
						TenantID:        "tenant-1",
						TenantName:      "Acme Corp",
						UserName:        "Test User",
						RoleNames:       []string{"admin"},
						CompanyID:       strPtr("company-1"),
						CompanyName:     strPtr("Acme Company"),
						EmailVerifiedAt: timePtr("2026-04-30T09:00:00Z"),
					}, nil)
			},
		},
		{
			name: "invalid email",
			input: coreentity.User{
				Email:      "nonexistent@example.com",
				Password:   "password123",
				TenantCode: "tenant-1",
			},
			setup: func(repo *dbMocks.UserRepository) {
				repo.EXPECT().
					FindActiveUserByEmailAndTenant(mock.Anything, "nonexistent@example.com", "tenant-1").
					Return(nil, errmsg.NewCustomErrors(404).SetMessage("user not found"))
			},
			wantError: true,
		},
		{
			name: "dependency error from repository",
			input: coreentity.User{
				Email:      "test@example.com",
				Password:   "password123",
				TenantCode: "tenant-1",
			},
			setup: func(repo *dbMocks.UserRepository) {
				repo.EXPECT().
					FindActiveUserByEmailAndTenant(mock.Anything, "test@example.com", "tenant-1").
					Return(nil, errmsg.NewCustomErrors(500).SetMessage("database error"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := dbMocks.NewUserRepository(t)

			if tt.setup != nil {
				tt.setup(repo)
			}

			core := NewUserCore(Config{
				Repo:       repo,
				TokenCache: tokencache.NewTokenCache(time.Hour, time.Minute),
			})

			got, err := core.Login(ctx, tt.input)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, got.AccessToken)
			require.NotEmpty(t, got.RefreshToken)
		})
	}
}
