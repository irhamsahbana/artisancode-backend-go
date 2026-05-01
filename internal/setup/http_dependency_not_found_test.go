package setup

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"codebase-app/pkg/errmsg"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestRegisterRouteNotFound(t *testing.T) {
	deps := httpDependencies{
		app: fiber.New(),
	}

	deps.registerRouteNotFound()

	req := httptest.NewRequest("GET", "/unknown-route?foo=bar", nil)
	resp, err := deps.app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var body map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.Equal(t, false, body["success"])
	require.Equal(t, errmsg.MessageRouteNotFound, body["message"])
	require.Equal(t, map[string]any{}, body["errors"])
}
