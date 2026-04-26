package core

import (
	corePorts "codebase-app/internal/ports/core"
	integrationPorts "codebase-app/internal/ports/integration"
	portsRepo "codebase-app/internal/ports/secondary/db"
)

var _ corePorts.UserInvitationCore = &userInvitationCore{}

type userInvitationCore struct {
	repo         portsRepo.UserInvitationRepository
	userRepo     portsRepo.UserRepository
	employeeRepo portsRepo.EmployeeRepository
	tx           portsRepo.Transactor
	bus          integrationPorts.MessagePublisher
}

type Config struct {
	Repo         portsRepo.UserInvitationRepository
	UserRepo     portsRepo.UserRepository
	EmployeeRepo portsRepo.EmployeeRepository
	Tx           portsRepo.Transactor
	Bus          integrationPorts.MessagePublisher
}

func NewUserInvitationCore(cfg Config) *userInvitationCore {
	return &userInvitationCore{
		repo:         cfg.Repo,
		userRepo:     cfg.UserRepo,
		employeeRepo: cfg.EmployeeRepo,
		tx:           cfg.Tx,
		bus:          cfg.Bus,
	}
}
