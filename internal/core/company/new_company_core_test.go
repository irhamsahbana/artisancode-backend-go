package core

import (
	"testing"

	"codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/require"
)

func TestNewCompanyCore(t *testing.T) {
	repo := mocks.NewCompanyRepository(t)

	core := NewCompanyCore(Config{
		Repo: repo,
	})

	require.NotNil(t, core)
	require.Same(t, repo, core.repo)
}
