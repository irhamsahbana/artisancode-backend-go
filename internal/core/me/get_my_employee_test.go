package core

import (
	"context"
	"errors"
	"testing"

	"codebase-app/internal/entity/coreentity"
	dbmocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetMyEmployee_ReturnsEmployeeFromRepository(t *testing.T) {
	ctx := context.Background()
	repo := dbmocks.NewMeRepository(t)
	want := &coreentity.Employee{ID: "employee-1", TenantID: "tenant-1", FullName: "Budi"}
	repo.EXPECT().
		GetEmployeeByUserID(mock.Anything, "tenant-1", "user-1").
		Return(want, nil)
	core := NewMeCore(Config{Repo: repo})

	got, err := core.GetMyEmployee(ctx, coreentity.SelfFilter{
		TenantID: "tenant-1",
		UserID:   "user-1",
	})

	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestGetMyEmployee_ReturnsRepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := dbmocks.NewMeRepository(t)
	wantErr := errors.New("repository failed")
	repo.EXPECT().
		GetEmployeeByUserID(mock.Anything, "tenant-1", "user-1").
		Return(nil, wantErr)
	core := NewMeCore(Config{Repo: repo})

	got, err := core.GetMyEmployee(ctx, coreentity.SelfFilter{
		TenantID: "tenant-1",
		UserID:   "user-1",
	})

	require.ErrorIs(t, err, wantErr)
	require.Nil(t, got)
}
