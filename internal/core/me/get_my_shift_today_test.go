package core

import (
	"context"
	"errors"
	"testing"
	"time"

	"codebase-app/internal/entity/coreentity"
	dbmocks "codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetMyShiftToday_ReturnsShiftWithDateInShiftTimezone(t *testing.T) {
	ctx := context.Background()
	repo := dbmocks.NewMeRepository(t)
	repo.EXPECT().
		GetWorkShiftByUserID(mock.Anything, "tenant-1", "user-1").
		Return(&coreentity.WorkShift{
			ID:        "shift-1",
			Name:      "Morning",
			StartTime: "08:00",
			EndTime:   "17:00",
			Timezone:  "Asia/Makassar",
		}, nil)
	core := NewMeCore(Config{Repo: repo})

	got, err := core.GetMyShiftToday(ctx, coreentity.SelfFilter{
		TenantID: "tenant-1",
		UserID:   "user-1",
	})

	require.NoError(t, err)
	require.Equal(t, "shift-1", got.ShiftID)
	require.Equal(t, "Morning", got.ShiftName)
	require.Equal(t, "08:00", got.StartTime)
	require.Equal(t, "17:00", got.EndTime)
	require.Equal(t, "Asia/Makassar", got.Timezone)
	require.Equal(t, time.Now().In(mustLoadLocation(t, "Asia/Makassar")).Format("2006-01-02"), got.AttendanceDate)
}

func TestGetMyShiftToday_ReturnsNilWhenRepositoryHasNoShift(t *testing.T) {
	ctx := context.Background()
	repo := dbmocks.NewMeRepository(t)
	repo.EXPECT().
		GetWorkShiftByUserID(mock.Anything, "tenant-1", "user-1").
		Return(nil, nil)
	core := NewMeCore(Config{Repo: repo})

	got, err := core.GetMyShiftToday(ctx, coreentity.SelfFilter{
		TenantID: "tenant-1",
		UserID:   "user-1",
	})

	require.NoError(t, err)
	require.Nil(t, got)
}

func TestGetMyShiftToday_ReturnsRepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := dbmocks.NewMeRepository(t)
	wantErr := errors.New("repository failed")
	repo.EXPECT().
		GetWorkShiftByUserID(mock.Anything, "tenant-1", "user-1").
		Return(nil, wantErr)
	core := NewMeCore(Config{Repo: repo})

	got, err := core.GetMyShiftToday(ctx, coreentity.SelfFilter{
		TenantID: "tenant-1",
		UserID:   "user-1",
	})

	require.ErrorIs(t, err, wantErr)
	require.Nil(t, got)
}

func mustLoadLocation(t *testing.T, name string) *time.Location {
	t.Helper()

	location, err := time.LoadLocation(name)
	require.NoError(t, err)
	return location
}
