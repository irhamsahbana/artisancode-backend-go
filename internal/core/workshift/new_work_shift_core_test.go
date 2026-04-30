package core

import (
	"testing"

	"codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/require"
)

func TestNewWorkShiftCore(t *testing.T) {
	repo := mocks.NewWorkShiftRepository(t)

	core := NewWorkShiftCore(Config{
		Repo: repo,
	})

	require.NotNil(t, core)
	require.Same(t, repo, core.repo)
}
