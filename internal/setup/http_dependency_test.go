package setup

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type httpRouteRegistrarSpy struct {
	calls []string
}

func (s *httpRouteRegistrarSpy) registerIdentityRoutes() {
	s.calls = append(s.calls, "identity")
}

func (s *httpRouteRegistrarSpy) registerOrganizationRoutes() {
	s.calls = append(s.calls, "organization")
}

func (s *httpRouteRegistrarSpy) registerAttendanceRoutes() {
	s.calls = append(s.calls, "attendance")
}

func (s *httpRouteRegistrarSpy) registerInternalRoutes() {
	s.calls = append(s.calls, "internal")
}

func (s *httpRouteRegistrarSpy) registerUtilityRoutes() {
	s.calls = append(s.calls, "utility")
}

func (s *httpRouteRegistrarSpy) registerRouteNotFound() {
	s.calls = append(s.calls, "not_found")
}

func TestRegisterHTTPRoutes(t *testing.T) {
	registrar := &httpRouteRegistrarSpy{}

	registerHTTPRoutes(registrar)

	require.Equal(t, []string{
		"identity",
		"organization",
		"attendance",
		"internal",
		"utility",
		"not_found",
	}, registrar.calls)
}

func TestHttpDependencies(t *testing.T) {
	originalFactory := newHTTPRouteRegistrar
	t.Cleanup(func() {
		newHTTPRouteRegistrar = originalFactory
	})

	registrar := &httpRouteRegistrarSpy{}
	factoryCalls := 0

	newHTTPRouteRegistrar = func() httpRouteRegistrar {
		factoryCalls++

		return registrar
	}

	HttpDependencies()

	require.Equal(t, 1, factoryCalls)
	require.Equal(t, []string{
		"identity",
		"organization",
		"attendance",
		"internal",
		"utility",
		"not_found",
	}, registrar.calls)
}
