package core

import (
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/secondary/db"
	integrationPorts "codebase-app/internal/ports/secondary/integration"
)

var _ corePorts.UserInvitationCore = &userInvitationCore{}

type userInvitationCore struct {
	repo         portsRepo.UserInvitationRepository
	userRepo     portsRepo.UserRepository
	employeeRepo portsRepo.EmployeeRepository
	bus          integrationPorts.MessagePublisher
}

type UserInvitationCoreConfig struct {
	Repo         portsRepo.UserInvitationRepository
	UserRepo     portsRepo.UserRepository
	EmployeeRepo portsRepo.EmployeeRepository
	Bus          integrationPorts.MessagePublisher
}

func NewUserInvitationCore(cfg UserInvitationCoreConfig) *userInvitationCore {
	return &userInvitationCore{
		repo:         cfg.Repo,
		userRepo:     cfg.UserRepo,
		employeeRepo: cfg.EmployeeRepo,
		bus:          cfg.Bus,
	}
}
