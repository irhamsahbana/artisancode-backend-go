package core

import (
	"testing"

	"codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/require"
)

func TestNewEmployeeCore(t *testing.T) {
	repo := mocks.NewEmployeeRepository(t)
	userRepo := mocks.NewUserRepository(t)

	core := NewEmployeeCore(Config{
		Repo:     repo,
		UserRepo: userRepo,
	})

	require.NotNil(t, core)
	require.Same(t, repo, core.repo)
	require.Same(t, userRepo, core.userRepo)
}
