package postgres

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Username        string
	Password        string
	Database        string
	Host            string
	Port            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

func New(cfg Config) (*sqlx.DB, error) {
	connectionString := "user=" + cfg.Username +
		" password=" + cfg.Password +
		" host=" + cfg.Host +
		" port=" + cfg.Port +
		" dbname=" + cfg.Database +
		" sslmode=" + cfg.SSLMode +
		" TimeZone=UTC"

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
