package http

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestAppBuildRegistersMetricsRoute(t *testing.T) {
	originalSync := syncHTTPAdapters
	originalRegister := registerHTTPDependencies
	t.Cleanup(func() {
		syncHTTPAdapters = originalSync
		registerHTTPDependencies = originalRegister
	})

	syncCalled := false
	registerCalled := false
	syncHTTPAdapters = func(app *fiber.App) {
		syncCalled = true
	}
	registerHTTPDependencies = func() {
		registerCalled = true
	}

	app := NewApp(AppConfig{})
	built, err := app.build("test-app", "v1", "test")

	require.NoError(t, err)
	require.True(t, syncCalled)
	require.True(t, registerCalled)

	resp, err := built.Test(httptest.NewRequest("GET", "/metrics", nil))
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(body), "artisancode_backend_build_info")
}
