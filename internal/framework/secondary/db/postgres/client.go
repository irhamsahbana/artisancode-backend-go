package postgres

import (
	"fmt"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Username        string
	Password        string
	Database        string
	Host            string
	Port            string
	SSLMode         string
	ChannelBinding  string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

func New(cfg Config) (*sqlx.DB, error) {
	sslMode, err := normalizeConnectionSetting("sslmode", cfg.SSLMode)
	if err != nil {
		return nil, err
	}

	channelBinding, err := normalizeConnectionSetting("channel_binding", cfg.ChannelBinding)
	if err != nil {
		return nil, err
	}

	connectionString := "user=" + cfg.Username +
		" password=" + cfg.Password +
		" host=" + cfg.Host +
		" port=" + cfg.Port +
		" dbname=" + cfg.Database +
		" TimeZone=UTC"

	if sslMode != "" {
		connectionString += " sslmode=" + sslMode

		if sslMode == "require" {
			sslNegotiation := "direct"
			connectionString += " sslnegotiation=" + sslNegotiation
		}
	}

	if channelBinding != "" {
		connectionString += " channel_binding=" + channelBinding
	}

	db, err := sqlx.Connect("pgx", connectionString)
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

func normalizeConnectionSetting(name, value string) (string, error) {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return "require", nil
	}

	if value != "require" && value != "disable" {
		return "", fmt.Errorf("invalid postgres %s: %s", name, value)
	}

	return value, nil
}
