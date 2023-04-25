package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	EnvKeyDBHost                = "DB_HOST"
	EnvKeyDBPort                = "DB_PORT"
	EnvKeyDBUsername            = "DB_USERNAME"
	EnvKeyDBPassword            = "DB_PASSWORD"
	EnvKeyDBName                = "DB_NAME"
	EnvKeyDBSSlMode             = "DB_SSL_MODE"
	EnvKeyLogLevel              = "LOG_LEVEL"
	EnvKeyDBMigrationSourcePath = "DB_MIGRATION_SOURCE_PATH"
)

var environmentVariables = map[string]struct{}{
	EnvKeyDBHost:                struct{}{},
	EnvKeyDBPort:                struct{}{},
	EnvKeyDBUsername:            struct{}{},
	EnvKeyDBPassword:            struct{}{},
	EnvKeyDBName:                struct{}{},
	EnvKeyDBSSlMode:             struct{}{},
	EnvKeyLogLevel:              struct{}{},
	EnvKeyDBMigrationSourcePath: struct{}{},
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"dbName"`
	SSLMode  string `json:"sslMode"`
}

type Config struct {
	LogLevel              string         `json:"logLevel"`
	Postgres              PostgresConfig `json:"postgres"`
	DBMigrationSourcePath string         `json:"dbMigrationSourcePath"`
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

func FromEnv() (*Config, error) {
	config := make(map[string]string)
	for envKey := range environmentVariables {
		value, ok := os.LookupEnv(envKey)
		if !ok {
			return nil, fmt.Errorf("unable to read environment variable: %s", envKey)
		}
		config[envKey] = value
	}
	return &Config{
			Postgres: PostgresConfig{
				Host:     config[EnvKeyDBHost],
				Port:     config[EnvKeyDBPort],
				Username: config[EnvKeyDBUsername],
				Password: config[EnvKeyDBPassword],
				DBName:   config[EnvKeyDBName],
				SSLMode:  config[EnvKeyDBSSlMode],
			},
			LogLevel:              config[EnvKeyLogLevel],
			DBMigrationSourcePath: config[EnvKeyDBMigrationSourcePath],
		},
		nil
}
