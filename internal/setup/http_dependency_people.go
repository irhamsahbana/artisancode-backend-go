package setup

import (
	companyCore "codebase-app/internal/core/company"
	employeeCore "codebase-app/internal/core/employee"
	jobpositionCore "codebase-app/internal/core/jobposition"
	meCore "codebase-app/internal/core/me"
	orgunitCore "codebase-app/internal/core/orgunit"
	rbacCore "codebase-app/internal/core/rbac"
	userCore "codebase-app/internal/core/user"
	userInvitationCore "codebase-app/internal/core/userinvitation"
	worklocationCore "codebase-app/internal/core/worklocation"
	workshiftCore "codebase-app/internal/core/workshift"
	companyHandler "codebase-app/internal/framework/primary/http/company"
	employeeHandler "codebase-app/internal/framework/primary/http/employee"
	jobpositionHandler "codebase-app/internal/framework/primary/http/jobposition"
	orgunitHandler "codebase-app/internal/framework/primary/http/orgunit"
	userHandler "codebase-app/internal/framework/primary/http/user"
	userInvitationHandler "codebase-app/internal/framework/primary/http/userinvitation"
	worklocationHandler "codebase-app/internal/framework/primary/http/worklocation"
	workshiftHandler "codebase-app/internal/framework/primary/http/workshift"
	companyRepo "codebase-app/internal/framework/secondary/db/postgres/company"
	employeeRepo "codebase-app/internal/framework/secondary/db/postgres/employee"
	jobpositionRepo "codebase-app/internal/framework/secondary/db/postgres/jobposition"
	meRepo "codebase-app/internal/framework/secondary/db/postgres/me"
	orgunitRepo "codebase-app/internal/framework/secondary/db/postgres/orgunit"
	rbacRepo "codebase-app/internal/framework/secondary/db/postgres/rbac"
	userRepo "codebase-app/internal/framework/secondary/db/postgres/user"
	userInvitationRepo "codebase-app/internal/framework/secondary/db/postgres/userinvitation"
	worklocationRepo "codebase-app/internal/framework/secondary/db/postgres/worklocation"
	workshiftRepo "codebase-app/internal/framework/secondary/db/postgres/workshift"
	"codebase-app/internal/middleware"
)

func buildPeopleDependencies(deps *httpDependencies, ctx *httpBootstrapContext) {
	userRepository := userRepo.NewUserRepository(userRepo.Config{
		DB: ctx.db,
	})
	companyRepository := companyRepo.NewCompanyRepository(companyRepo.Config{
		DB: ctx.db,
	})
	orgUnitRepository := orgunitRepo.NewOrgUnitRepository(orgunitRepo.Config{
		DB: ctx.db,
	})
	jobPositionRepository := jobpositionRepo.NewJobPositionRepository(jobpositionRepo.Config{
		DB: ctx.db,
	})
	workLocationRepository := worklocationRepo.NewWorkLocationRepository(worklocationRepo.Config{
		DB: ctx.db,
	})
	workShiftRepository := workshiftRepo.NewWorkShiftRepository(workshiftRepo.Config{
		DB: ctx.db,
	})
	employeeRepository := employeeRepo.NewEmployeeRepository(employeeRepo.Config{
		DB: ctx.db,
	})
	userInvitationRepository := userInvitationRepo.NewUserInvitationRepository(
		userInvitationRepo.Config{
			DB: ctx.db,
		},
	)
	meRepository := meRepo.NewMeRepository(meRepo.Config{
		DB: ctx.db,
	})
	rbacRepository := rbacRepo.NewRbacRepository(rbacRepo.Config{
		DB: ctx.db,
	})

	deps.userCore = userCore.NewUserCore(userCore.Config{
		Repo:                 userRepository,
		Tx:                   ctx.tx,
		TokenCache:           ctx.tokenCache,
		Bus:                  ctx.bus,
		GoogleTokenValidator: ctx.googleTokenValidator,
	})
	deps.companyCore = companyCore.NewCompanyCore(companyCore.Config{
		Repo: companyRepository,
	})
	deps.orgUnitCore = orgunitCore.NewOrgUnitCore(orgunitCore.Config{
		Repo: orgUnitRepository,
	})
	deps.jobPositionCore = jobpositionCore.NewJobPositionCore(jobpositionCore.Config{
		Repo: jobPositionRepository,
	})
	deps.workLocationCore = worklocationCore.NewWorkLocationCore(
		worklocationCore.Config{
			Repo:        workLocationRepository,
			OrgUnitRepo: orgUnitRepository,
		},
	)
	deps.workShiftCore = workshiftCore.NewWorkShiftCore(workshiftCore.Config{
		Repo: workShiftRepository,
	})
	deps.employeeCore = employeeCore.NewEmployeeCore(employeeCore.Config{
		Repo:     employeeRepository,
		UserRepo: userRepository,
	})
	deps.userInvitationCore = userInvitationCore.NewUserInvitationCore(
		userInvitationCore.Config{
			Repo:         userInvitationRepository,
			UserRepo:     userRepository,
			EmployeeRepo: employeeRepository,
			Tx:           ctx.tx,
			Bus:          ctx.bus,
		},
	)
	deps.meCore = meCore.NewMeCore(meCore.Config{
		Repo: meRepository,
	})
	deps.rbacCore = rbacCore.NewRbacCore(rbacCore.Config{
		Repo: rbacRepository,
	})
}

func (deps httpDependencies) registerIdentityRoutes() {
	userHandlerInst := userHandler.NewUserHandler(userHandler.Config{
		Core:        deps.userCore,
		RateLimiter: deps.authRateLimiter,
	})
	userHandlerInst.Register(deps.app.Group("/users"))
	userHandlerInst.RegisterTenant(deps.app.Group("/tenant", middleware.Auth))

	userInvitationHandler.NewUserInvitationHandler(userInvitationHandler.Config{
		Core: deps.userInvitationCore,
	}).Register(deps.app.Group("/user-invitations"))
}

func (deps httpDependencies) registerOrganizationRoutes() {
	companyHandler.NewCompanyHandler(companyHandler.Config{
		Core: deps.companyCore,
	}).Register(deps.app.Group("/companies", middleware.Auth))
	orgunitHandler.NewOrgUnitHandler(orgunitHandler.Config{
		Core: deps.orgUnitCore,
	}).Register(deps.app.Group("/org-units", middleware.Auth))
	jobpositionHandler.NewJobPositionHandler(jobpositionHandler.Config{
		Core: deps.jobPositionCore,
	}).Register(deps.app.Group("/job-positions", middleware.Auth))
	worklocationHandler.NewWorkLocationHandler(worklocationHandler.Config{
		Core: deps.workLocationCore,
	}).Register(deps.app.Group("/work-locations", middleware.Auth))
	workshiftHandler.NewWorkShiftHandler(workshiftHandler.Config{
		Core: deps.workShiftCore,
	}).Register(deps.app.Group("/work-shifts", middleware.Auth))
	employeeHandler.NewEmployeeHandler(employeeHandler.Config{
		Core: deps.employeeCore,
	}).Register(deps.app.Group("/employees", middleware.Auth))
}
