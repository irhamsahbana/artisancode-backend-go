package setup

import (
	internalClientCore "codebase-app/internal/core/internalclient"
	internalInvoiceCore "codebase-app/internal/core/internalinvoice"
	internalOrderCore "codebase-app/internal/core/internalorder"
	internalProductCore "codebase-app/internal/core/internalproduct"
	internalQuotationCore "codebase-app/internal/core/internalquotation"
	internalUserCore "codebase-app/internal/core/internaluser"
	internalClientHandler "codebase-app/internal/framework/primary/http/internalclient"
	internalInvoiceHandler "codebase-app/internal/framework/primary/http/internalinvoice"
	internalOrderHandler "codebase-app/internal/framework/primary/http/internalorder"
	internalProductHandler "codebase-app/internal/framework/primary/http/internalproduct"
	internalQuotationHandler "codebase-app/internal/framework/primary/http/internalquotation"
	internalUserHandler "codebase-app/internal/framework/primary/http/internaluser"
	internalClientRepo "codebase-app/internal/framework/secondary/db/postgres/internalclient"
	internalInvoiceRepo "codebase-app/internal/framework/secondary/db/postgres/internalinvoice"
	internalOrderRepo "codebase-app/internal/framework/secondary/db/postgres/internalorder"
	internalProductRepo "codebase-app/internal/framework/secondary/db/postgres/internalproduct"
	internalQuotationRepo "codebase-app/internal/framework/secondary/db/postgres/internalquotation"
	internalUserRepo "codebase-app/internal/framework/secondary/db/postgres/internaluser"
	"codebase-app/internal/middleware"
)

func buildInternalDependencies(deps *httpDependencies, ctx httpBootstrapContext) {
	internalProductRepository := internalProductRepo.NewInternalProductRepository(
		internalProductRepo.Config{
			DB: ctx.db,
		},
	)
	internalClientRepository := internalClientRepo.NewInternalClientRepository(
		internalClientRepo.Config{
			DB: ctx.db,
		},
	)
	internalUserRepository := internalUserRepo.NewInternalUserRepository(
		internalUserRepo.Config{
			DB: ctx.db,
		},
	)
	internalQuotationRepository := internalQuotationRepo.NewInternalQuotationRepository(
		internalQuotationRepo.Config{
			DB: ctx.db,
		},
	)
	internalOrderRepository := internalOrderRepo.NewInternalOrderRepository(
		internalOrderRepo.Config{
			DB: ctx.db,
		},
	)
	internalInvoiceRepository := internalInvoiceRepo.NewInternalInvoiceRepository(
		internalInvoiceRepo.Config{
			DB: ctx.db,
		},
	)

	deps.internalProductCore = internalProductCore.NewInternalProductCore(
		internalProductCore.Config{
			Repo: internalProductRepository,
		},
	)
	deps.internalClientCore = internalClientCore.NewInternalClientCore(
		internalClientCore.Config{
			Repo: internalClientRepository,
		},
	)
	deps.internalQuotationCore = internalQuotationCore.NewInternalQuotationCore(
		internalQuotationCore.Config{
			QuotationRepo: internalQuotationRepository,
			OrderRepo:     internalOrderRepository,
			InvoiceRepo:   internalInvoiceRepository,
			Tx:            ctx.tx,
		},
	)
	deps.internalOrderCore = internalOrderCore.NewInternalOrderCore(
		internalOrderCore.Config{
			OrderRepo:   internalOrderRepository,
			InvoiceRepo: internalInvoiceRepository,
			Tx:          ctx.tx,
		},
	)
	deps.internalInvoiceCore = internalInvoiceCore.NewInternalInvoiceCore(
		internalInvoiceCore.Config{
			InvoiceRepo:   internalInvoiceRepository,
			OrderRepo:     internalOrderRepository,
			QuotationRepo: internalQuotationRepository,
			Tx:            ctx.tx,
			DOKU:          ctx.dokuClient,
		},
	)
	deps.internalUserCore = internalUserCore.NewInternalUserCore(
		internalUserCore.Config{
			Repo:       internalUserRepository,
			TokenCache: ctx.tokenCache,
		},
	)
}

func (deps httpDependencies) registerInternalRoutes() {
	internalProductHandler.NewInternalProductHandler(
		internalProductHandler.Config{
			Core: deps.internalProductCore,
		},
	).Register(deps.app.Group("/internal-products", middleware.InternalAuth))
	internalClientHandler.NewInternalClientHandler(
		internalClientHandler.Config{
			Core: deps.internalClientCore,
		},
	).Register(deps.app.Group("/internal-clients", middleware.InternalAuth))

	internalCommerceGroup := deps.app.Group(
		"/internal-commerce",
		middleware.InternalAuth,
	)
	internalQuotationHandler.NewInternalQuotationHandler(
		internalQuotationHandler.Config{
			Core: deps.internalQuotationCore,
		},
	).Register(internalCommerceGroup)
	internalOrderHandler.NewInternalOrderHandler(internalOrderHandler.Config{
		Core: deps.internalOrderCore,
	}).Register(internalCommerceGroup)
	internalInvoiceHandler.NewInternalInvoiceHandler(
		internalInvoiceHandler.Config{
			Core: deps.internalInvoiceCore,
		},
	).Register(internalCommerceGroup)

	internalUserHandler.NewInternalUserHandler(internalUserHandler.Config{
		Core: deps.internalUserCore,
	}).Register(deps.app.Group("/internal-users"))
}
