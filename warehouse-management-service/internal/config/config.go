package config

import (
	"encoding/json"
	"os"
)

type Postgres struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"dbName"`
	SSLMode  bool   `json:"sslMode"`
}

type DBMigration struct {
	SourcePath string `json:"path"`
}

type Config struct {
	Postgres    Postgres    `json:"postgres"`
	DBMigration DBMigration `json:"dbMigration"`
}

func FromFile(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
