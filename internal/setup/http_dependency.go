package setup

type httpRouteRegistrar interface {
	registerIdentityRoutes()
	registerOrganizationRoutes()
	registerAttendanceRoutes()
	registerInternalRoutes()
	registerUtilityRoutes()
	registerRouteNotFound()
}

var newHTTPRouteRegistrar = func() httpRouteRegistrar {
	deps := newHTTPDependencies()

	return deps
}

func HttpDependencies() {
	registerHTTPRoutes(newHTTPRouteRegistrar())
}

func registerHTTPRoutes(registrar httpRouteRegistrar) {
	registrar.registerIdentityRoutes()
	registrar.registerOrganizationRoutes()
	registrar.registerAttendanceRoutes()
	registrar.registerInternalRoutes()
	registrar.registerUtilityRoutes()
	registrar.registerRouteNotFound()
}
