package setup

import (
	webhookCore "codebase-app/internal/core/webhook"
	meHandler "codebase-app/internal/framework/primary/http/me"
	rbacHandler "codebase-app/internal/framework/primary/http/rbac"
	webhookHandler "codebase-app/internal/framework/primary/http/webhook"
	"codebase-app/internal/middleware"
)

func buildUtilityDependencies(deps *httpDependencies, ctx httpBootstrapContext) {
	deps.webhookCore = webhookCore.NewWebhookCore(webhookCore.Config{
		DOKUVerifier: ctx.dokuClient,
	})
}

func (deps httpDependencies) registerUtilityRoutes() {
	meHandler.NewMeHandler(meHandler.Config{
		Core: deps.meCore,
	}).Register(deps.app.Group("/me", middleware.Auth))
	webhookHandler.NewWebhookHandler(webhookHandler.Config{
		Core: deps.webhookCore,
	}).Register(deps.app.Group("/webhooks"))
	deps.registerStorageRoutes()
	rbacHandler.NewRbacHandler(rbacHandler.Config{
		Core: deps.rbacCore,
	}).Register(deps.app.Group("/role-and-permissions", middleware.Auth))
}
