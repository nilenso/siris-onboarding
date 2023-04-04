package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var environmentVariables = map[string]string{
	"dbHost":                "DB_HOST",
	"dbPort":                "DB_PORT",
	"dbUsername":            "DB_USERNAME",
	"dbPassword":            "DB_PASSWORD",
	"dbName":                "DB_NAME",
	"dbSSLMode":             "DB_SSL_MODE",
	"logLevel":              "LOG_LEVEL",
	"dbMigrationSourcePath": "DB_MIGRATION_SOURCE_PATH",
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
	for identifier, envVar := range environmentVariables {
		value, ok := os.LookupEnv(envVar)
		if !ok {
			return nil, fmt.Errorf("unable to read environment variable: %s", envVar)
		}
		config[identifier] = value
	}
	return &Config{
			Postgres: PostgresConfig{
				Host:     config["dbHost"],
				Port:     config["dbPort"],
				Username: config["dbUsername"],
				Password: config["dbPassword"],
				DBName:   config["dbName"],
				SSLMode:  config["dbSSLMode"],
			},
			LogLevel:              config["logLevel"],
			DBMigrationSourcePath: config["dbMigrationSourcePath"],
		},
		nil
}
