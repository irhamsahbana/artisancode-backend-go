package repository

import (
	postgresTx "codebase-app/internal/framework/secondary/db/postgres/transaction"
	portsRepo "codebase-app/internal/ports/secondary/db"
	"context"

	"github.com/jmoiron/sqlx"
)

type userInvitationRepo struct {
	db *sqlx.DB
}

func (r *userInvitationRepo) executor(ctx context.Context) postgresTx.SQLExecutor {
	return postgresTx.ExecutorFromContext(ctx, r.db)
}

var _ portsRepo.UserInvitationRepository = &userInvitationRepo{}

type UserInvitationRepositoryConfig struct {
	DB *sqlx.DB
}

func NewUserInvitationRepository(cfg UserInvitationRepositoryConfig) portsRepo.UserInvitationRepository {
	return &userInvitationRepo{db: cfg.DB}
}
