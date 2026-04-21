package core

import (
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/secondary/db"
)

var _ corePorts.UserInvitationCore = &userInvitationCore{}

type userInvitationCore struct {
	repo         portsRepo.UserInvitationRepository
	userRepo     portsRepo.UserRepository
	employeeRepo portsRepo.EmployeeRepository
}

type UserInvitationCoreConfig struct {
	Repo         portsRepo.UserInvitationRepository
	UserRepo     portsRepo.UserRepository
	EmployeeRepo portsRepo.EmployeeRepository
}

func NewUserInvitationCore(cfg UserInvitationCoreConfig) *userInvitationCore {
	return &userInvitationCore{
		repo:         cfg.Repo,
		userRepo:     cfg.UserRepo,
		employeeRepo: cfg.EmployeeRepo,
	}
}
