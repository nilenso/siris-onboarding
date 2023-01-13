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

func New(config config.Postgres) (*Postgres, error) {
	connect := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%v",
		config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLMode)
	db, err := sql.Open("postgres", connect)
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

func (p *Postgres) Close() error {
	return p.db.Close()
}
