package config

import (
	"reflect"
	"testing"
)

func TestFromFile(t *testing.T) {
	config, err := FromFile("./testdata/config.json")
	wantConfig := &Config{
		LogLevel: "debug",
		Postgres: PostgresConfig{
			Host:     "localhost",
			Port:     "5432",
			Username: "user",
			Password: "xxxxxxxx",
			DBName:   "db",
			SSLMode:  "disable",
		},
		DBMigration: DBMigration{SourcePath: "file://warehouse-management-service/db/migrations"},
	}

	if err != nil || !reflect.DeepEqual(config, wantConfig) {
		t.Errorf("want config: %v, err: %v, got config: %v, err: %v", wantConfig, nil, config, err)
	}
}

func TestFromFileError(t *testing.T) {
	_, err := FromFile("bad path")
	if err == nil {
		t.Errorf("got %v, want %v:", err, "error")
	}
}

func TestFromFileParseError(t *testing.T) {
	_, err := FromFile("./testdata/bad_config.json")
	if err == nil {
		t.Errorf("got %v, want %v:", err, "json error")
	}
}
