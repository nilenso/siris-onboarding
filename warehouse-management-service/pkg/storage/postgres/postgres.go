package postgres

import (
	"database/sql"
	"fmt"
	storage "warehouse-management-service/pkg/storage"

	// Loads postgres drivers
	_ "github.com/lib/pq"
)

type postgres struct {
	db *sql.DB
}

func New(host, port, user, password, dbName string, sslMode bool) (storage.Service, error) {
	connect := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
	db, err := sql.Open("postgres", connect)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &postgres{
		db: db,
	}, nil
}

func (p *postgres) Close() error {
	return p.db.Close()
}
