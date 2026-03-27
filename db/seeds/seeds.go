package seeds

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type Seed struct {
	db *sqlx.DB
}

func newSeed(db *sqlx.DB) Seed {
	return Seed{
		db: db,
	}
}

func Execute(db *sqlx.DB, table string, total int) {
	_ = total
	seed := newSeed(db)
	seed.run(table)
}

func (s *Seed) run(table string) {
	switch table {
	case "tenant", "company":
		s.tenantSeed()
	case "owner_user":
		s.ownerUserSeed()
	case "rbac":
		s.rbacSeed()
	case "org_unit":
		s.orgUnitSeed()
	case "job_position":
		s.jobPositionSeed()
	case "work_location":
		s.workLocationSeed()
	case "work_shift":
		s.workShiftSeed()
	case "template":
		s.rbacSeed()
		s.orgUnitSeed()
	case "all":
		s.tenantSeed()
		s.rbacSeed()
		s.ownerUserSeed()
	case "all_with_template":
		s.rbacSeed()
		s.orgUnitSeed()
		s.tenantSeed()
		s.ownerUserSeed()
		s.jobPositionSeed()
		s.workLocationSeed()
		s.workShiftSeed()
	default:
		log.Warn().Str("table", table).Msg("No seed to run")
	}
}
