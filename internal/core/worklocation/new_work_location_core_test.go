package core

import (
	"testing"

	"codebase-app/internal/ports/secondary/db/mocks"

	"github.com/stretchr/testify/require"
)

func TestNewWorkLocationCore(t *testing.T) {
	repo := mocks.NewWorkLocationRepository(t)
	orgUnitRepo := mocks.NewOrgUnitRepository(t)

	core := NewWorkLocationCore(Config{
		Repo:        repo,
		OrgUnitRepo: orgUnitRepo,
	})

	require.NotNil(t, core)
	require.Same(t, repo, core.repo)
	require.Same(t, orgUnitRepo, core.orgUnitRepo)
}
