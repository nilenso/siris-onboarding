package postgres

import (
	"database/sql"
	"fmt"
	"warehouse-management-service/internal/config"

	// Loads postgres drivers
	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func New(connectionURL string) (*Postgres, error) {
	db, err := sql.Open("postgres", connectionURL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Postgres{
		db: db,
	}, nil
}

func Connection(config config.Postgres) string {
	return fmt.Sprintf("postgresql://%s@%s:%s/%s?sslmode=%s", config.Username, config.Host, config.Port, config.DBName, config.SSLMode)
}

func (p *Postgres) Close() error {
	return p.db.Close()
}
