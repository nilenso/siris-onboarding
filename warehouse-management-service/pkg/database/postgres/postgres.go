package postgres

import (
	"database/sql"
	"fmt"
	"warehouse-management-service/internal/config"

	// Loads postgres drivers
	_ "github.com/lib/pq"
)

type Postgres struct {
	config.PostgresConfig
}

func New(config config.PostgresConfig) *Postgres {
	return &Postgres{PostgresConfig: config}
}

func (p *Postgres) Open() (*sql.DB, error) {
	connectionURL := buildConnectionURL(p.PostgresConfig)
	db, err := sql.Open("postgres", connectionURL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func buildConnectionURL(config config.PostgresConfig) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", config.Username, config.Password, config.Host, config.Port, config.DBName, config.SSLMode)
}
