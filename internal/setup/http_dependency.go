package setup

import (
	"time"

	"codebase-app/internal/adapter"
	attendanceCore "codebase-app/internal/core/attendance"
	companyCore "codebase-app/internal/core/company"
	employeeCore "codebase-app/internal/core/employee"
	exportJobCore "codebase-app/internal/core/export_job"
	internalClientCore "codebase-app/internal/core/internalclient"
	internalInvoiceCore "codebase-app/internal/core/internalinvoice"
	internalOrderCore "codebase-app/internal/core/internalorder"
	internalProductCore "codebase-app/internal/core/internalproduct"
	internalQuotationCore "codebase-app/internal/core/internalquotation"
	internalUserCore "codebase-app/internal/core/internaluser"
	jobpositionCore "codebase-app/internal/core/jobposition"
	meCore "codebase-app/internal/core/me"
	orgunitCore "codebase-app/internal/core/orgunit"
	rbacCore "codebase-app/internal/core/rbac"
	storageCore "codebase-app/internal/core/storage"
	userCore "codebase-app/internal/core/user"
	userInvitationCore "codebase-app/internal/core/userinvitation"
	webhookCore "codebase-app/internal/core/webhook"
	worklocationCore "codebase-app/internal/core/worklocation"
	workshiftCore "codebase-app/internal/core/workshift"
	attendanceHandler "codebase-app/internal/framework/primary/http/attendance"
	companyHandler "codebase-app/internal/framework/primary/http/company"
	employeeHandler "codebase-app/internal/framework/primary/http/employee"
	exportJobHandler "codebase-app/internal/framework/primary/http/export_job"
	internalClientHandler "codebase-app/internal/framework/primary/http/internalclient"
	internalInvoiceHandler "codebase-app/internal/framework/primary/http/internalinvoice"
	internalOrderHandler "codebase-app/internal/framework/primary/http/internalorder"
	internalProductHandler "codebase-app/internal/framework/primary/http/internalproduct"
	internalQuotationHandler "codebase-app/internal/framework/primary/http/internalquotation"
	internalUserHandler "codebase-app/internal/framework/primary/http/internaluser"
	jobpositionHandler "codebase-app/internal/framework/primary/http/jobposition"
	meHandler "codebase-app/internal/framework/primary/http/me"
	orgunitHandler "codebase-app/internal/framework/primary/http/orgunit"
	rbacHandler "codebase-app/internal/framework/primary/http/rbac"
	storageHandler "codebase-app/internal/framework/primary/http/storage"
	userHandler "codebase-app/internal/framework/primary/http/user"
	userInvitationHandler "codebase-app/internal/framework/primary/http/userinvitation"
	webhookHandler "codebase-app/internal/framework/primary/http/webhook"
	worklocationHandler "codebase-app/internal/framework/primary/http/worklocation"
	workshiftHandler "codebase-app/internal/framework/primary/http/workshift"
	attendanceRepo "codebase-app/internal/framework/secondary/db/postgres/attendance"
	companyRepo "codebase-app/internal/framework/secondary/db/postgres/company"
	employeeRepo "codebase-app/internal/framework/secondary/db/postgres/employee"
	exportJobRepo "codebase-app/internal/framework/secondary/db/postgres/export_job"
	internalClientRepo "codebase-app/internal/framework/secondary/db/postgres/internalclient"
	internalInvoiceRepo "codebase-app/internal/framework/secondary/db/postgres/internalinvoice"
	internalOrderRepo "codebase-app/internal/framework/secondary/db/postgres/internalorder"
	internalProductRepo "codebase-app/internal/framework/secondary/db/postgres/internalproduct"
	internalQuotationRepo "codebase-app/internal/framework/secondary/db/postgres/internalquotation"
	internalUserRepo "codebase-app/internal/framework/secondary/db/postgres/internaluser"
	jobpositionRepo "codebase-app/internal/framework/secondary/db/postgres/jobposition"
	meRepo "codebase-app/internal/framework/secondary/db/postgres/me"
	orgunitRepo "codebase-app/internal/framework/secondary/db/postgres/orgunit"
	rbacRepo "codebase-app/internal/framework/secondary/db/postgres/rbac"
	storageRepo "codebase-app/internal/framework/secondary/db/postgres/storage"
	postgresTx "codebase-app/internal/framework/secondary/db/postgres/transaction"
	userRepo "codebase-app/internal/framework/secondary/db/postgres/user"
	userInvitationRepo "codebase-app/internal/framework/secondary/db/postgres/userinvitation"
	worklocationRepo "codebase-app/internal/framework/secondary/db/postgres/worklocation"
	workshiftRepo "codebase-app/internal/framework/secondary/db/postgres/workshift"

	dokuIntegration "codebase-app/internal/integration/doku"
	storage "codebase-app/internal/integration/storage"

	"github.com/gofiber/fiber/v2"

	"codebase-app/internal/infrastructure"
	"codebase-app/internal/integration/ratelimit"
	"codebase-app/internal/integration/tokencache"
	"codebase-app/internal/middleware"
	"codebase-app/pkg/response"
)

func HttpDependencies() {
	var (
		app = adapter.Adapters.RestServer
		db  = adapter.Adapters.Postgres
		bus = adapter.Adapters.MessagePublisher
		s3  = storage.NewStorageIntegration(adapter.Adapters.Storage)
		tx  = postgresTx.NewTransactor(db)
	)

	userRepository := userRepo.NewUserRepository(userRepo.Config{
		DB: db,
	})
	companyRepository := companyRepo.NewCompanyRepository(companyRepo.Config{
		DB: db,
	})
	orgUnitRepository := orgunitRepo.NewOrgUnitRepository(orgunitRepo.Config{
		DB: db,
	})
	jobPositionRepository := jobpositionRepo.NewJobPositionRepository(jobpositionRepo.Config{
		DB: db,
	})
	workLocationRepository := worklocationRepo.NewWorkLocationRepository(worklocationRepo.Config{
		DB: db,
	})
	workShiftRepository := workshiftRepo.NewWorkShiftRepository(workshiftRepo.Config{
		DB: db,
	})
	meRepository := meRepo.NewMeRepository(meRepo.Config{
		DB: db,
	})
	rbacRepository := rbacRepo.NewRbacRepository(rbacRepo.Config{
		DB: db,
	})

	tokenCache := tokencache.NewTokenCache(time.Hour*24*7, time.Minute*10)
	authRateLimiter := ratelimit.NewCacheLimiter()
	attendanceRepository := attendanceRepo.NewAttendanceRepository(attendanceRepo.Config{
		DB: db,
	})
	exportJobRepository := exportJobRepo.NewExportJobRepository(exportJobRepo.Config{
		DB: db,
	})
	storageRepository := storageRepo.NewStorageRepository(storageRepo.Config{
		DB: db,
	})
	userInvitationRepository := userInvitationRepo.NewUserInvitationRepository(userInvitationRepo.Config{
		DB: db,
	})
	internalProductRepository := internalProductRepo.NewInternalProductRepository(internalProductRepo.Config{
		DB: db,
	})
	internalClientRepository := internalClientRepo.NewInternalClientRepository(internalClientRepo.Config{
		DB: db,
	})
	internalUserRepository := internalUserRepo.NewInternalUserRepository(internalUserRepo.Config{
		DB: db,
	})
	internalQuotationRepository := internalQuotationRepo.NewInternalQuotationRepository(internalQuotationRepo.Config{
		DB: db,
	})
	internalOrderRepository := internalOrderRepo.NewInternalOrderRepository(internalOrderRepo.Config{
		DB: db,
	})
	internalInvoiceRepository := internalInvoiceRepo.NewInternalInvoiceRepository(internalInvoiceRepo.Config{
		DB: db,
	})

	userCoreInst := userCore.NewUserCore(userCore.Config{
		Repo:       userRepository,
		Tx:         tx,
		TokenCache: tokenCache,
		Bus:        bus,
	})
	companyCoreInst := companyCore.NewCompanyCore(companyCore.Config{
		Repo: companyRepository,
	})
	orgUnitCoreInst := orgunitCore.NewOrgUnitCore(orgunitCore.Config{
		Repo: orgUnitRepository,
	})
	jobPositionCoreInst := jobpositionCore.NewJobPositionCore(jobpositionCore.Config{
		Repo: jobPositionRepository,
	})
	workLocationCoreInst := worklocationCore.NewWorkLocationCore(worklocationCore.Config{
		Repo:        workLocationRepository,
		OrgUnitRepo: orgUnitRepository,
	})
	workShiftCoreInst := workshiftCore.NewWorkShiftCore(workshiftCore.Config{
		Repo: workShiftRepository,
	})
	employeeRepository := employeeRepo.NewEmployeeRepository(employeeRepo.Config{
		DB: db,
	})
	rbacCoreInst := rbacCore.NewRbacCore(rbacCore.Config{
		Repo: rbacRepository,
	})

	employeeCoreInst := employeeCore.NewEmployeeCore(employeeCore.Config{
		Repo:     employeeRepository,
		UserRepo: userRepository,
	})
	userInvitationCoreInst := userInvitationCore.NewUserInvitationCore(userInvitationCore.Config{
		Repo:         userInvitationRepository,
		UserRepo:     userRepository,
		EmployeeRepo: employeeRepository,
		Tx:           tx,
		Bus:          bus,
	})
	dokuClient := dokuIntegration.NewClientFromEnv()
	attendanceCoreInst := attendanceCore.NewAttendanceCore(attendanceCore.Config{
		Repo:        attendanceRepository,
		StorageRepo: storageRepository,
		Tx:          tx,
		S3:          s3,
	})
	exportJobCoreInst := exportJobCore.NewExportJobCore(exportJobCore.Config{
		Repo:           exportJobRepository,
		AttendanceRepo: attendanceRepository,
		StorageRepo:    storageRepository,
		Tx:             tx,
		S3:             s3,
		Bus:            bus,
	})
	internalProductCoreInst := internalProductCore.NewInternalProductCore(internalProductCore.Config{
		Repo: internalProductRepository,
	})
	internalClientCoreInst := internalClientCore.NewInternalClientCore(internalClientCore.Config{
		Repo: internalClientRepository,
	})
	internalQuotationCoreInst := internalQuotationCore.NewInternalQuotationCore(internalQuotationCore.Config{
		QuotationRepo: internalQuotationRepository,
		OrderRepo:     internalOrderRepository,
		InvoiceRepo:   internalInvoiceRepository,
		Tx:            tx,
	})
	internalOrderCoreInst := internalOrderCore.NewInternalOrderCore(internalOrderCore.Config{
		OrderRepo:   internalOrderRepository,
		InvoiceRepo: internalInvoiceRepository,
		Tx:          tx,
	})
	internalInvoiceCoreInst := internalInvoiceCore.NewInternalInvoiceCore(internalInvoiceCore.Config{
		InvoiceRepo:   internalInvoiceRepository,
		OrderRepo:     internalOrderRepository,
		QuotationRepo: internalQuotationRepository,
		Tx:            tx,
		DOKU:          dokuClient,
	})
	internalUserCoreInst := internalUserCore.NewInternalUserCore(internalUserCore.Config{
		Repo:       internalUserRepository,
		TokenCache: tokenCache,
	})
	storageCoreInst := storageCore.NewStorageCore(s3, storageRepository)
	meCoreInst := meCore.NewMeCore(meCore.Config{
		Repo: meRepository,
	})
	webhookCoreInst := webhookCore.NewWebhookCore(webhookCore.Config{
		DOKUVerifier: dokuClient,
	})
	userHandler.NewUserHandler(userHandler.Config{
		Core:        userCoreInst,
		RateLimiter: authRateLimiter,
	}).Register(app.Group("/users"))
	userInvitationHandler.NewUserInvitationHandler(userInvitationHandler.Config{
		Core: userInvitationCoreInst,
	}).Register(app.Group("/user-invitations"))
	companyHandler.NewCompanyHandler(companyHandler.Config{
		Core: companyCoreInst,
	}).Register(app.Group("/companies", middleware.Auth))
	orgunitHandler.NewOrgUnitHandler(orgunitHandler.Config{
		Core: orgUnitCoreInst,
	}).Register(app.Group("/org-units", middleware.Auth))
	jobpositionHandler.NewJobPositionHandler(jobpositionHandler.Config{
		Core: jobPositionCoreInst,
	}).Register(app.Group("/job-positions", middleware.Auth))
	worklocationHandler.NewWorkLocationHandler(worklocationHandler.Config{
		Core: workLocationCoreInst,
	}).Register(app.Group("/work-locations", middleware.Auth))
	workshiftHandler.NewWorkShiftHandler(workshiftHandler.Config{
		Core: workShiftCoreInst,
	}).Register(app.Group("/work-shifts", middleware.Auth))
	employeeHandler.NewEmployeeHandler(employeeHandler.Config{
		Core: employeeCoreInst,
	}).Register(app.Group("/employees", middleware.Auth))
	attendanceHandlerInst := attendanceHandler.NewAttendanceHandler(attendanceHandler.Config{
		Core: attendanceCoreInst,
	})
	attendanceHandlerInst.Register(app.Group("/attendance-logs", middleware.Auth))
	exportJobHandler.NewExportJobHandler(exportJobHandler.Config{
		Core: exportJobCoreInst,
	}).Register(app.Group("/export-jobs", middleware.Auth))
	internalProductHandler.NewInternalProductHandler(internalProductHandler.Config{
		Core: internalProductCoreInst,
	}).Register(app.Group("/internal-products", middleware.InternalAuth))
	internalClientHandler.NewInternalClientHandler(internalClientHandler.Config{
		Core: internalClientCoreInst,
	}).Register(app.Group("/internal-clients", middleware.InternalAuth))
	internalCommerceGroup := app.Group("/internal-commerce", middleware.InternalAuth)
	internalQuotationHandler.NewInternalQuotationHandler(internalQuotationHandler.Config{
		Core: internalQuotationCoreInst,
	}).Register(internalCommerceGroup)
	internalOrderHandler.NewInternalOrderHandler(internalOrderHandler.Config{
		Core: internalOrderCoreInst,
	}).Register(internalCommerceGroup)
	internalInvoiceHandler.NewInternalInvoiceHandler(internalInvoiceHandler.Config{
		Core: internalInvoiceCoreInst,
	}).Register(internalCommerceGroup)
	internalUserHandler.NewInternalUserHandler(internalUserHandler.Config{
		Core: internalUserCoreInst,
	}).Register(app.Group("/internal-users"))
	attendanceHandlerInst.RegisterSummary(app.Group("/attendance-summary", middleware.Auth))
	attendanceHandlerInst.RegisterPolicy(app.Group("/attendance-policy", middleware.Auth))
	meHandler.NewMeHandler(meHandler.Config{
		Core: meCoreInst,
	}).Register(app.Group("/me", middleware.Auth))
	webhookHandler.NewWebhookHandler(webhookHandler.Config{
		Core: webhookCoreInst,
	}).Register(app.Group("/webhooks"))
	storageHandler.NewStorageHandler(storageCoreInst).Register(app.Group("/storage"))
	rbacHandler.NewRbacHandler(rbacHandler.Config{
		Core: rbacCoreInst,
	}).Register(app.Group("/role-and-permissions", middleware.Auth))

	app.Use(func(c *fiber.Ctx) error {
		var (
			method = c.Method()
			path   = c.Path()
			query  = c.Context().QueryArgs().String()
			ua     = c.Get("User-Agent")
			ip     = c.IP()
		)

		infrastructure.AccessLogger.Debug().
			Str("method", method).
			Str("path", path).
			Str("query", query).
			Str("ua", ua).
			Str("ip", ip).
			Msg("route not found")
		return c.Status(fiber.StatusNotFound).JSON(response.Error("Route not found"))
	})
}
