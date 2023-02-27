package postgres

import "testing"

func TestRunMigrationError(t *testing.T) {
	err := postgres.RunMigration("bad source url")
	if err == nil {
		t.Error(err)
	}
}
