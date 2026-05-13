package core

import (
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/secondary/db"

	"golang.org/x/crypto/bcrypt"
)

var _ corePorts.EmployeeCore = &employeeCore{}

type employeeCore struct {
	repo        portsRepo.EmployeeRepository
	userRepo    portsRepo.UserRepository
	billingCore corePorts.InternalTenantBillingCore
}

type Config struct {
	Repo        portsRepo.EmployeeRepository
	UserRepo    portsRepo.UserRepository
	BillingCore corePorts.InternalTenantBillingCore
}

func NewEmployeeCore(cfg Config) *employeeCore {
	return &employeeCore{
		repo:        cfg.Repo,
		userRepo:    cfg.UserRepo,
		billingCore: cfg.BillingCore,
	}
}

// hashPassword generates a bcrypt hash from a plain password
func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
