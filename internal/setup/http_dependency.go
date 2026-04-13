package setup

import (
	"time"

	attendanceCore "codebase-app/internal/core/attendance"
	companyCore "codebase-app/internal/core/company"
	employeeCore "codebase-app/internal/core/employee"
	exportJobCore "codebase-app/internal/core/export_job"
	jobpositionCore "codebase-app/internal/core/jobposition"
	meCore "codebase-app/internal/core/me"
	orgunitCore "codebase-app/internal/core/orgunit"
	rbacCore "codebase-app/internal/core/rbac"
	storageCore "codebase-app/internal/core/storage"
	userCore "codebase-app/internal/core/user"
	worklocationCore "codebase-app/internal/core/worklocation"
	workshiftCore "codebase-app/internal/core/workshift"
	attendanceHandler "codebase-app/internal/framework/primary/http/attendance"
	companyHandler "codebase-app/internal/framework/primary/http/company"
	employeeHandler "codebase-app/internal/framework/primary/http/employee"
	exportJobHandler "codebase-app/internal/framework/primary/http/export_job"
	jobpositionHandler "codebase-app/internal/framework/primary/http/jobposition"
	meHandler "codebase-app/internal/framework/primary/http/me"
	orgunitHandler "codebase-app/internal/framework/primary/http/orgunit"
	rbacHandler "codebase-app/internal/framework/primary/http/rbac"
	storageHandler "codebase-app/internal/framework/primary/http/storage"
	userHandler "codebase-app/internal/framework/primary/http/user"
	worklocationHandler "codebase-app/internal/framework/primary/http/worklocation"
	workshiftHandler "codebase-app/internal/framework/primary/http/workshift"
	attendanceRepo "codebase-app/internal/framework/secondary/db/postgres/attendance"
	companyRepo "codebase-app/internal/framework/secondary/db/postgres/company"
	employeeRepo "codebase-app/internal/framework/secondary/db/postgres/employee"
	exportJobRepo "codebase-app/internal/framework/secondary/db/postgres/export_job"
	jobpositionRepo "codebase-app/internal/framework/secondary/db/postgres/jobposition"
	meRepo "codebase-app/internal/framework/secondary/db/postgres/me"
	orgunitRepo "codebase-app/internal/framework/secondary/db/postgres/orgunit"
	rbacRepo "codebase-app/internal/framework/secondary/db/postgres/rbac"
	storageRepo "codebase-app/internal/framework/secondary/db/postgres/storage"
	userRepo "codebase-app/internal/framework/secondary/db/postgres/user"
	worklocationRepo "codebase-app/internal/framework/secondary/db/postgres/worklocation"
	workshiftRepo "codebase-app/internal/framework/secondary/db/postgres/workshift"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"codebase-app/internal/infrastructure"
	"codebase-app/internal/integration/tokencache"
	"codebase-app/internal/middleware"
	integrationPorts "codebase-app/internal/ports/secondary/integration"
	"codebase-app/pkg/response"
)

func HttpDependencies(
	app *fiber.App,
	db *sqlx.DB,
	s3 integrationPorts.StorageContract,
	bus integrationPorts.MessagePublisher,
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
	meRepository := meRepo.NewMeRepository(meRepo.MeRepositoryConfig{
		DB: db,
	})
	rbacRepository := rbacRepo.NewRbacRepository(rbacRepo.RbacRepositoryConfig{
		DB: db,
	})

	tokenCache := tokencache.NewTokenCache(time.Hour*24*7, time.Minute*10)
	attendanceRepository := attendanceRepo.NewAttendanceRepository(attendanceRepo.AttendanceRepositoryConfig{
		DB: db,
	})
	exportJobRepository := exportJobRepo.NewExportJobRepository(exportJobRepo.ExportJobRepositoryConfig{
		DB: db,
	})
	storageRepository := storageRepo.NewStorageRepository(storageRepo.StorageRepositoryConfig{
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
		Repo:        attendanceRepository,
		CompanyRepo: companyRepository,
		StorageRepo: storageRepository,
		S3:          s3,
	})
	exportJobCoreInst := exportJobCore.NewExportJobCore(exportJobCore.ExportJobCoreConfig{
		Repo:           exportJobRepository,
		AttendanceRepo: attendanceRepository,
		StorageRepo:    storageRepository,
		S3:             s3,
		Bus:            bus,
	})
	storageCoreInst := storageCore.NewStorageCore(s3, storageRepository)
	meCoreInst := meCore.NewMeCore(meCore.MeCoreConfig{
		Repo: meRepository,
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
	attendanceHandlerInst := attendanceHandler.NewAttendanceHandler(attendanceHandler.AttendanceHandlerConfig{
		Core: attendanceCoreInst,
	})
	attendanceHandlerInst.Register(app.Group("/attendance-logs", middleware.Auth))
	exportJobHandler.NewExportJobHandler(exportJobHandler.ExportJobHandlerConfig{
		Core: exportJobCoreInst,
	}).Register(app.Group("/export-jobs", middleware.Auth))
	attendanceHandlerInst.RegisterSummary(app.Group("/attendance-summary", middleware.Auth))
	attendanceHandlerInst.RegisterPolicy(app.Group("/attendance-policy", middleware.Auth))
	meHandler.NewMeHandler(meHandler.MeHandlerConfig{
		Core: meCoreInst,
	}).Register(app.Group("/me", middleware.Auth))
	storageHandler.NewStorageHandler(storageCoreInst).Register(app.Group("/storage"))
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
