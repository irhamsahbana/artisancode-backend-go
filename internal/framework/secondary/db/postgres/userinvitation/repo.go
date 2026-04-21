package repository

import (
	portsRepo "codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
)

type userInvitationRepo struct {
	db *sqlx.DB
}

var _ portsRepo.UserInvitationRepository = &userInvitationRepo{}

type UserInvitationRepositoryConfig struct {
	DB *sqlx.DB
}

func NewUserInvitationRepository(cfg UserInvitationRepositoryConfig) portsRepo.UserInvitationRepository {
	return &userInvitationRepo{db: cfg.DB}
}
