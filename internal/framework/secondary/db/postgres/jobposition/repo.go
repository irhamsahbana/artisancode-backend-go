package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/repoentity"
	"codebase-app/internal/infrastructure/tracing"
	portsRepo "codebase-app/internal/ports/secondary/db"
	"codebase-app/pkg/errmsg"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type jobPositionRepo struct {
	db *sqlx.DB
}

var _ portsRepo.JobPositionRepository = &jobPositionRepo{}

type JobPositionRepositoryConfig struct {
	DB *sqlx.DB
}

func NewJobPositionRepository(cfg JobPositionRepositoryConfig) portsRepo.JobPositionRepository {
	return &jobPositionRepo{db: cfg.DB}
}

func (r *jobPositionRepo) GetJobPositions(ctx context.Context, filter coreentity.JobPositionListFilter) ([]coreentity.JobPosition, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:jobposition:repo:GetJobPositions")
	defer span.End()

	type dao struct {
		TotalData int `db:"total_data"`
		repoentity.JobPosition
	}

	var (
		data  = make([]dao, 0)
		args  = make([]any, 0, 4)
		items = make([]coreentity.JobPosition, 0)
		total = 0
	)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id, tenant_id, name, grade
		FROM job_positions
		WHERE deleted_at IS NULL AND tenant_id = ?
	`

	args = append(args, filter.TenantID)

	if filter.Q != "" {
		query += ` AND name ILIKE '%' || ? || '%'`
		args = append(args, filter.Q)
	}

	query += ` ORDER BY name ASC LIMIT ? OFFSET ?`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), args...)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query job positions")
		return nil, 0, err
	}

	for _, d := range data {
		total = d.TotalData
		items = append(items, coreentity.JobPosition{
			ID:       d.ID,
			TenantID: d.TenantID,
			Name:     d.Name,
			Grade:    d.Grade,
		})
	}

	return items, total, nil
}

func (r *jobPositionRepo) GetJobPosition(ctx context.Context, filter coreentity.JobPosition) (*coreentity.JobPosition, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:jobposition:repo:GetJobPosition")
	defer span.End()

	var data = new(repoentity.JobPosition)

	query := `
		SELECT id, tenant_id, name, grade
		FROM job_positions
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, data, r.db.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Job position not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get job position")
		return nil, err
	}

	result := &coreentity.JobPosition{
		ID:       data.ID,
		TenantID: data.TenantID,
		Name:     data.Name,
		Grade:    data.Grade,
	}
	return result, nil
}

func (r *jobPositionRepo) CreateJobPosition(ctx context.Context, data coreentity.JobPosition) (*coreentity.JobPosition, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:jobposition:repo:CreateJobPosition")
	defer span.End()

	query := `
		INSERT INTO job_positions (tenant_id, name, grade)
		VALUES (?, ?, ?)
		RETURNING id
	`
	var id string
	err := r.db.GetContext(ctx, &id, r.db.Rebind(query), data.TenantID, data.Name, data.Grade)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create job position")
		return nil, err
	}

	data.ID = id
	return &data, nil
}

func (r *jobPositionRepo) UpdateJobPosition(ctx context.Context, data coreentity.JobPosition) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:jobposition:repo:UpdateJobPosition")
	defer span.End()

	query := `
		UPDATE job_positions
		SET name = ?, grade = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), data.Name, data.Grade, data.ID, data.TenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update job position")
		return err
	}
	return nil
}

func (r *jobPositionRepo) DeleteJobPosition(ctx context.Context, filter coreentity.JobPositionDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:jobposition:repo:DeleteJobPosition")
	defer span.End()

	query := `
		UPDATE job_positions
		SET deleted_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to delete job position")
		return err
	}
	return nil
}
