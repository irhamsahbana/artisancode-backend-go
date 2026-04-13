package repository

import (
	"codebase-app/internal/ports/secondary/db"

	"github.com/jmoiron/sqlx"
)

var _ repository.UserRepository = &userRepo{}

type userRepo struct {
	db *sqlx.DB
}

type UserRepositoryConfig struct {
	DB *sqlx.DB
}

func NewUserRepository(cfg UserRepositoryConfig) repository.UserRepository {
	return &userRepo{
		db: cfg.DB,
	}
}
