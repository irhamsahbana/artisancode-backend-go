package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	portsRepo "codebase-app/internal/ports/repository"
	"codebase-app/pkg/errmsg"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type workShiftRepo struct {
	db *sqlx.DB
}

var _ portsRepo.WorkShiftRepository = &workShiftRepo{}

type WorkShiftRepositoryConfig struct {
	DB *sqlx.DB
}

func NewWorkShiftRepository(cfg WorkShiftRepositoryConfig) portsRepo.WorkShiftRepository {
	return &workShiftRepo{db: cfg.DB}
}

func (r *workShiftRepo) GetWorkShifts(ctx context.Context, filter coreentity.WorkShiftListFilter) ([]coreentity.WorkShift, int, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.GetWorkShifts")
	defer span.End()

	type dao struct {
		TotalData          int    `db:"total_data"`
		ID                 string `db:"id"`
		TenantID           string `db:"tenant_id"`
		Name               string `db:"name"`
		StartTime          string `db:"start_time"`
		EndTime            string `db:"end_time"`
		GracePeriodMinutes int    `db:"grace_period_minutes"`
	}

	var (
		data  = make([]dao, 0)
		args  = make([]any, 0, 4)
		items = make([]coreentity.WorkShift, 0)
		total = 0
	)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id, tenant_id, name, start_time, end_time, grace_period_minutes
		FROM work_shifts
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
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query work shifts")
		return nil, 0, err
	}

	for _, d := range data {
		total = d.TotalData
		items = append(items, coreentity.WorkShift{
			ID:                 d.ID,
			TenantID:           d.TenantID,
			Name:               d.Name,
			StartTime:          d.StartTime,
			EndTime:            d.EndTime,
			GracePeriodMinutes: d.GracePeriodMinutes,
		})
	}

	return items, total, nil
}

func (r *workShiftRepo) GetWorkShift(ctx context.Context, filter coreentity.WorkShift) (*coreentity.WorkShift, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.GetWorkShift")
	defer span.End()

	var data struct {
		ID                 string `db:"id"`
		TenantID           string `db:"tenant_id"`
		Name               string `db:"name"`
		StartTime          string `db:"start_time"`
		EndTime            string `db:"end_time"`
		GracePeriodMinutes int    `db:"grace_period_minutes"`
	}

	query := `
		SELECT id, tenant_id, name, start_time, end_time, grace_period_minutes
		FROM work_shifts
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &data, r.db.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Shift tidak ditemukan")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get work shift")
		return nil, err
	}

	result := &coreentity.WorkShift{
		ID:                 data.ID,
		TenantID:           data.TenantID,
		Name:               data.Name,
		StartTime:          data.StartTime,
		EndTime:            data.EndTime,
		GracePeriodMinutes: data.GracePeriodMinutes,
	}
	return result, nil
}

func (r *workShiftRepo) CreateWorkShift(ctx context.Context, data coreentity.WorkShift) (*coreentity.WorkShift, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.CreateWorkShift")
	defer span.End()

	query := `
		INSERT INTO work_shifts (tenant_id, name, start_time, end_time, grace_period_minutes)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id
	`
	var id string
	err := r.db.GetContext(ctx, &id, r.db.Rebind(query),
		data.TenantID, data.Name, data.StartTime, data.EndTime, data.GracePeriodMinutes,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create work shift")
		return nil, err
	}

	data.ID = id
	return &data, nil
}

func (r *workShiftRepo) UpdateWorkShift(ctx context.Context, data coreentity.WorkShift) error {
	ctx, span := tracing.StartSpan(ctx, "repo.UpdateWorkShift")
	defer span.End()

	query := `
		UPDATE work_shifts
		SET name = ?, start_time = ?, end_time = ?, grace_period_minutes = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query),
		data.Name, data.StartTime, data.EndTime, data.GracePeriodMinutes, data.ID, data.TenantID,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update work shift")
		return err
	}
	return nil
}

func (r *workShiftRepo) DeleteWorkShift(ctx context.Context, filter coreentity.WorkShiftDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "repo.DeleteWorkShift")
	defer span.End()

	query := `
		UPDATE work_shifts
		SET deleted_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to delete work shift")
		return err
	}
	return nil
}
