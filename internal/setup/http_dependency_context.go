package setup

import (
	"time"

	"codebase-app/internal/adapter"
	postgresTx "codebase-app/internal/framework/secondary/db/postgres/transaction"
	dokuIntegration "codebase-app/internal/integration/doku"
	"codebase-app/internal/integration/googleidentity"
	"codebase-app/internal/integration/ratelimit"
	storageIntegration "codebase-app/internal/integration/storage"
	"codebase-app/internal/integration/tokencache"
	corePorts "codebase-app/internal/ports/core"
	integrationPorts "codebase-app/internal/ports/integration"
	repositoryPorts "codebase-app/internal/ports/secondary/db"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type httpBootstrapContext struct {
	app                  *fiber.App
	db                   *sqlx.DB
	bus                  integrationPorts.MessagePublisher
	tx                   repositoryPorts.Transactor
	s3                   integrationPorts.StorageContract
	tokenCache           tokencache.TokenCacheContract
	authRateLimiter      ratelimit.AttemptLimiter
	googleTokenValidator integrationPorts.GoogleIDTokenValidator
	dokuClient           integrationPorts.DokuClient
}

type httpDependencies struct {
	app             *fiber.App
	authRateLimiter ratelimit.AttemptLimiter

	userCore           corePorts.UserCore
	companyCore        corePorts.CompanyCore
	orgUnitCore        corePorts.OrgUnitCore
	jobPositionCore    corePorts.JobPositionCore
	workLocationCore   corePorts.WorkLocationCore
	workShiftCore      corePorts.WorkShiftCore
	employeeCore       corePorts.EmployeeCore
	userInvitationCore corePorts.UserInvitationCore
	meCore             corePorts.MeCore
	rbacCore           corePorts.RbacCore

	attendanceCore corePorts.AttendanceCore
	exportJobCore  corePorts.ExportJobCore
	storageCore    corePorts.StorageCore

	internalProductCore       corePorts.InternalProductCore
	internalCurrencyCore      corePorts.InternalCurrencyCore
	internalClientCore        corePorts.InternalClientCore
	internalQuotationCore     corePorts.InternalQuotationCore
	internalOrderCore         corePorts.InternalOrderCore
	internalInvoiceCore       corePorts.InternalInvoiceCore
	internalTenantBillingCore corePorts.InternalTenantBillingCore
	internalUserCore          corePorts.InternalUserCore

	webhookCore corePorts.WebhookCore
}

func newHTTPDependencies() httpDependencies {
	ctx := newHTTPBootstrapContext()
	deps := httpDependencies{
		app:             ctx.app,
		authRateLimiter: ctx.authRateLimiter,
	}

	buildPeopleDependencies(&deps, ctx)
	buildAttendanceDependencies(&deps, ctx)
	buildInternalDependencies(&deps, ctx)
	buildUtilityDependencies(&deps, ctx)

	return deps
}

func newHTTPBootstrapContext() httpBootstrapContext {
	db := adapter.Adapters.Postgres

	return httpBootstrapContext{
		app:                  adapter.Adapters.RestServer,
		db:                   db,
		bus:                  adapter.Adapters.MessagePublisher,
		tx:                   postgresTx.NewTransactor(db),
		s3:                   storageIntegration.NewStorageIntegration(adapter.Adapters.Storage),
		tokenCache:           tokencache.NewTokenCache(time.Hour*24*7, time.Minute*10),
		authRateLimiter:      ratelimit.NewCacheLimiter(),
		googleTokenValidator: googleidentity.NewValidatorFromEnv(),
		dokuClient:           dokuIntegration.NewClientFromEnv(),
	}
}
