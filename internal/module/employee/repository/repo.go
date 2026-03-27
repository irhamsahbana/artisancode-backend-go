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

type employeeRepo struct {
	db *sqlx.DB
}

var _ portsRepo.EmployeeRepository = &employeeRepo{}

type EmployeeRepositoryConfig struct {
	DB *sqlx.DB
}

func NewEmployeeRepository(cfg EmployeeRepositoryConfig) portsRepo.EmployeeRepository {
	return &employeeRepo{db: cfg.DB}
}

func (r *employeeRepo) GetEmployees(ctx context.Context, filter coreentity.EmployeeListFilter) ([]coreentity.Employee, int, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.GetEmployees")
	defer span.End()

	type dao struct {
		TotalData int `db:"total_data"`
		ID           string  `db:"id"`
		TenantID     string  `db:"tenant_id"`
		EmployeeNo   string  `db:"employee_no"`
		FullName     string  `db:"full_name"`
		OrgUnitID    *string `db:"org_unit_id"`
		JobPositionID *string `db:"job_position_id"`
		LocationID   *string `db:"location_id"`
		ShiftID      *string `db:"shift_id"`
		Status       string  `db:"status"`
		JoinDate     *string `db:"join_date"`
	}

	var (
		data  = make([]dao, 0)
		args  = make([]any, 0, 8)
		items = make([]coreentity.Employee, 0)
		total = 0
	)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id, tenant_id, employee_no, full_name, org_unit_id, job_position_id,
			location_id, shift_id, status, join_date
		FROM employees
		WHERE deleted_at IS NULL AND tenant_id = ?
	`

	args = append(args, filter.TenantID)

	if filter.Status != "" {
		query += ` AND status = ?`
		args = append(args, filter.Status)
	}
	if filter.OrgUnitID != nil {
		query += ` AND org_unit_id = ?`
		args = append(args, filter.OrgUnitID)
	}
	if filter.Q != "" {
		query += ` AND (full_name ILIKE '%' || ? || '%' OR employee_no ILIKE '%' || ? || '%')`
		args = append(args, filter.Q, filter.Q)
	}

	query += ` ORDER BY full_name ASC LIMIT ? OFFSET ?`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), args...)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query employees")
		return nil, 0, err
	}

	for _, d := range data {
		total = d.TotalData
		items = append(items, coreentity.Employee{
			ID:            d.ID,
			TenantID:      d.TenantID,
			EmployeeNo:    d.EmployeeNo,
			FullName:      d.FullName,
			OrgUnitID:     d.OrgUnitID,
			JobPositionID: d.JobPositionID,
			LocationID:    d.LocationID,
			ShiftID:       d.ShiftID,
			Status:        d.Status,
			JoinDate:      d.JoinDate,
		})
	}

	return items, total, nil
}

func (r *employeeRepo) GetEmployee(ctx context.Context, filter coreentity.Employee) (*coreentity.Employee, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.GetEmployee")
	defer span.End()

	var data struct {
		ID            string  `db:"id"`
		TenantID      string  `db:"tenant_id"`
		EmployeeNo    string  `db:"employee_no"`
		FullName      string  `db:"full_name"`
		OrgUnitID     *string `db:"org_unit_id"`
		JobPositionID *string `db:"job_position_id"`
		LocationID   *string `db:"location_id"`
		ShiftID      *string `db:"shift_id"`
		Status        string  `db:"status"`
		JoinDate      *string `db:"join_date"`
	}

	query := `
		SELECT
			id, tenant_id, employee_no, full_name, org_unit_id, job_position_id,
			location_id, shift_id, status, join_date
		FROM employees
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &data, r.db.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Employee not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get employee")
		return nil, err
	}

	result := &coreentity.Employee{
		ID:            data.ID,
		TenantID:      data.TenantID,
		EmployeeNo:    data.EmployeeNo,
		FullName:      data.FullName,
		OrgUnitID:     data.OrgUnitID,
		JobPositionID: data.JobPositionID,
		LocationID:    data.LocationID,
		ShiftID:       data.ShiftID,
		Status:        data.Status,
		JoinDate:      data.JoinDate,
	}
	return result, nil
}

func (r *employeeRepo) CreateEmployee(ctx context.Context, data coreentity.Employee) (*coreentity.Employee, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.CreateEmployee")
	defer span.End()

	query := `
		INSERT INTO employees (
			tenant_id, employee_no, full_name, org_unit_id, job_position_id, location_id, shift_id, status, join_date
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`
	var id string
	err := r.db.GetContext(ctx, &id, r.db.Rebind(query),
		data.TenantID,
		data.EmployeeNo,
		data.FullName,
		data.OrgUnitID,
		data.JobPositionID,
		data.LocationID,
		data.ShiftID,
		data.Status,
		data.JoinDate,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create employee")
		return nil, err
	}

	data.ID = id
	return &data, nil
}

func (r *employeeRepo) UpdateEmployee(ctx context.Context, data coreentity.Employee) error {
	ctx, span := tracing.StartSpan(ctx, "repo.UpdateEmployee")
	defer span.End()

	query := `
		UPDATE employees
		SET
			tenant_id = ?, employee_no = ?, full_name = ?, org_unit_id = ?,
			job_position_id = ?, location_id = ?, shift_id = ?, status = ?, join_date = ?,
			updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query),
		data.TenantID,
		data.EmployeeNo,
		data.FullName,
		data.OrgUnitID,
		data.JobPositionID,
		data.LocationID,
		data.ShiftID,
		data.Status,
		data.JoinDate,
		data.ID,
		data.TenantID,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update employee")
		return err
	}
	return nil
}

func (r *employeeRepo) DeleteEmployee(ctx context.Context, filter coreentity.EmployeeDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "repo.DeleteEmployee")
	defer span.End()

	query := `
		UPDATE employees
		SET deleted_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), filter.ID, filter.TenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to delete employee")
		return err
	}
	return nil
}