package postgres

import (
	"testing"
	"warehouse-management-service/internal/config"
)

func TestOpenError(t *testing.T) {
	pg := New(config.PostgresConfig{})
	_, err := pg.Open()
	if err == nil {
		t.Errorf("got %v, want %v", err, "error")
	}
}
