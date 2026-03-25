package setup

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"codebase-app/internal/infrastructure"
	"codebase-app/internal/integration/tokencache"
	companyCore "codebase-app/internal/module/company/core"
	companyHandler "codebase-app/internal/module/company/handler"
	companyRepo "codebase-app/internal/module/company/repository"
	jobpositionCore "codebase-app/internal/module/jobposition/core"
	jobpositionHandler "codebase-app/internal/module/jobposition/handler"
	jobpositionRepo "codebase-app/internal/module/jobposition/repository"
	orgunitCore "codebase-app/internal/module/orgunit/core"
	orgunitHandler "codebase-app/internal/module/orgunit/handler"
	orgunitRepo "codebase-app/internal/module/orgunit/repository"
	userCore "codebase-app/internal/module/user/core"
	userHandler "codebase-app/internal/module/user/handler"
	userRepo "codebase-app/internal/module/user/repository"
	worklocationCore "codebase-app/internal/module/worklocation/core"
	worklocationHandler "codebase-app/internal/module/worklocation/handler"
	worklocationRepo "codebase-app/internal/module/worklocation/repository"
	workshiftCore "codebase-app/internal/module/workshift/core"
	workshiftHandler "codebase-app/internal/module/workshift/handler"
	workshiftRepo "codebase-app/internal/module/workshift/repository"
	"codebase-app/internal/ports/integration"
	"codebase-app/pkg/response"
)

func Dependencies(
	app *fiber.App,
	db *sqlx.DB,
	s3 integration.StorageContract,
) {
	userRepository := userRepo.NewUserRepository(userRepo.UserRepositoryConfig{
		DB: db,
	})
	companyRepository := companyRepo.NewCompanyRepository(companyRepo.CompanyRepositoryConfig{
		DB: db,
	})
	orgUnitRepository := orgunitRepo.NewOrgUnitRepository(orgunitRepo.OrgUnitRepositoryConfig{
		DB: db,
	})
	jobPositionRepository := jobpositionRepo.NewJobPositionRepository(jobpositionRepo.JobPositionRepositoryConfig{
		DB: db,
	})
	workLocationRepository := worklocationRepo.NewWorkLocationRepository(worklocationRepo.WorkLocationRepositoryConfig{
		DB: db,
	})
	workShiftRepository := workshiftRepo.NewWorkShiftRepository(workshiftRepo.WorkShiftRepositoryConfig{
		DB: db,
	})

	tokenCache := tokencache.NewTokenCache(time.Hour*24*7, time.Minute*10)

	userCoreInst := userCore.NewUserCore(userCore.UserCoreConfig{
		Repo:       userRepository,
		TokenCache: tokenCache,
	})
	companyCoreInst := companyCore.NewCompanyCore(companyCore.CompanyCoreConfig{
		Repo: companyRepository,
	})
	orgUnitCoreInst := orgunitCore.NewOrgUnitCore(orgunitCore.OrgUnitCoreConfig{
		Repo: orgUnitRepository,
	})
	jobPositionCoreInst := jobpositionCore.NewJobPositionCore(jobpositionCore.JobPositionCoreConfig{
		Repo: jobPositionRepository,
	})
	workLocationCoreInst := worklocationCore.NewWorkLocationCore(worklocationCore.WorkLocationCoreConfig{
		Repo: workLocationRepository,
	})
	workShiftCoreInst := workshiftCore.NewWorkShiftCore(workshiftCore.WorkShiftCoreConfig{
		Repo: workShiftRepository,
	})

	userHandler.NewUserHandler(userHandler.UserHandlerConfig{
		Core: userCoreInst,
	}).Register(app.Group("/users"))
	companyHandler.NewCompanyHandler(companyHandler.CompanyHandlerConfig{
		Core: companyCoreInst,
	}).Register(app.Group("/companies"))
	orgunitHandler.NewOrgUnitHandler(orgunitHandler.OrgUnitHandlerConfig{
		Core: orgUnitCoreInst,
	}).Register(app.Group("/org-units"))
	jobpositionHandler.NewJobPositionHandler(jobpositionHandler.JobPositionHandlerConfig{
		Core: jobPositionCoreInst,
	}).Register(app.Group("/job-positions"))
	worklocationHandler.NewWorkLocationHandler(worklocationHandler.WorkLocationHandlerConfig{
		Core: workLocationCoreInst,
	}).Register(app.Group("/work-locations"))
	workshiftHandler.NewWorkShiftHandler(workshiftHandler.WorkShiftHandlerConfig{
		Core: workShiftCoreInst,
	}).Register(app.Group("/work-shifts"))

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
