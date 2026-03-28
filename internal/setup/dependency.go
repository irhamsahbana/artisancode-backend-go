package setup

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"codebase-app/internal/infrastructure"
	"codebase-app/internal/integration/tokencache"
	"codebase-app/internal/middleware"
	attendanceCore "codebase-app/internal/module/attendance/core"
	attendanceHandler "codebase-app/internal/module/attendance/handler"
	attendanceRepo "codebase-app/internal/module/attendance/repository"
	companyCore "codebase-app/internal/module/company/core"
	companyHandler "codebase-app/internal/module/company/handler"
	companyRepo "codebase-app/internal/module/company/repository"
	employeeCore "codebase-app/internal/module/employee/core"
	employeeHandler "codebase-app/internal/module/employee/handler"
	employeeRepo "codebase-app/internal/module/employee/repository"
	jobpositionCore "codebase-app/internal/module/jobposition/core"
	jobpositionHandler "codebase-app/internal/module/jobposition/handler"
	jobpositionRepo "codebase-app/internal/module/jobposition/repository"
	orgunitCore "codebase-app/internal/module/orgunit/core"
	orgunitHandler "codebase-app/internal/module/orgunit/handler"
	orgunitRepo "codebase-app/internal/module/orgunit/repository"
	rbacCore "codebase-app/internal/module/rbac/core"
	rbacHandler "codebase-app/internal/module/rbac/handler"
	rbacRepo "codebase-app/internal/module/rbac/repository"
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
	rbacRepository := rbacRepo.NewRbacRepository(rbacRepo.RbacRepositoryConfig{
		DB: db,
	})

	tokenCache := tokencache.NewTokenCache(time.Hour*24*7, time.Minute*10)
	attendanceRepository := attendanceRepo.NewAttendanceRepository(attendanceRepo.AttendanceRepositoryConfig{
		DB: db,
	})

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
		Repo:        workLocationRepository,
		OrgUnitRepo: orgUnitRepository,
	})
	workShiftCoreInst := workshiftCore.NewWorkShiftCore(workshiftCore.WorkShiftCoreConfig{
		Repo: workShiftRepository,
	})
	employeeRepository := employeeRepo.NewEmployeeRepository(employeeRepo.EmployeeRepositoryConfig{
		DB: db,
	})
	rbacCoreInst := rbacCore.NewRbacCore(rbacCore.RbacCoreConfig{
		Repo: rbacRepository,
	})

	employeeCoreInst := employeeCore.NewEmployeeCore(employeeCore.EmployeeCoreConfig{
		Repo:     employeeRepository,
		UserRepo: userRepository,
	})
	attendanceCoreInst := attendanceCore.NewAttendanceCore(attendanceCore.AttendanceCoreConfig{
		Repo: attendanceRepository,
	})
	userHandler.NewUserHandler(userHandler.UserHandlerConfig{
		Core: userCoreInst,
	}).Register(app.Group("/users"))
	companyHandler.NewCompanyHandler(companyHandler.CompanyHandlerConfig{
		Core: companyCoreInst,
	}).Register(app.Group("/companies", middleware.Auth))
	orgunitHandler.NewOrgUnitHandler(orgunitHandler.OrgUnitHandlerConfig{
		Core: orgUnitCoreInst,
	}).Register(app.Group("/org-units", middleware.Auth))
	jobpositionHandler.NewJobPositionHandler(jobpositionHandler.JobPositionHandlerConfig{
		Core: jobPositionCoreInst,
	}).Register(app.Group("/job-positions", middleware.Auth))
	worklocationHandler.NewWorkLocationHandler(worklocationHandler.WorkLocationHandlerConfig{
		Core: workLocationCoreInst,
	}).Register(app.Group("/work-locations", middleware.Auth))
	workshiftHandler.NewWorkShiftHandler(workshiftHandler.WorkShiftHandlerConfig{
		Core: workShiftCoreInst,
	}).Register(app.Group("/work-shifts", middleware.Auth))
	employeeHandler.NewEmployeeHandler(employeeHandler.EmployeeHandlerConfig{
		Core: employeeCoreInst,
	}).Register(app.Group("/employees", middleware.Auth))
	attendanceHandler.NewAttendanceHandler(attendanceHandler.AttendanceHandlerConfig{
		Core: attendanceCoreInst,
	}).Register(app.Group("/attendance-logs", middleware.Auth))
	rbacHandler.NewRbacHandler(rbacHandler.RbacHandlerConfig{
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
