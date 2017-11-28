package migrataur

import (
	"testing"
	"time"
)

// mockAdapter implements an in memory migrataur adapter used for testing
type mockAdapter struct {
	appliedMigrations []*Migration
}

func newMockAdapter() *mockAdapter {
	return &mockAdapter{}
}

func (a *mockAdapter) CreateMigrationsTableIfNotExists() error {
	return nil
}

func (a *mockAdapter) AddMigration(name string, at time.Time) error {
	a.appliedMigrations = append(a.appliedMigrations, NewAdapterMigration(name, at))

	return nil
}

func (a *mockAdapter) RemoveMigration(name string) error {
	for i, m := range a.appliedMigrations {
		if m.name == name {
			a.appliedMigrations = append(a.appliedMigrations[:i], a.appliedMigrations[i+1:]...)
			break
		}
	}

	return nil
}

func (a *mockAdapter) Exec(command string) error {
	return nil
}

func (a *mockAdapter) GetAll() ([]*Migration, error) {
	return a.appliedMigrations, nil
}

func TestMockAdapter(t *testing.T) {
	adapter := newMockAdapter()

	if len(adapter.appliedMigrations) != 0 {
		t.Error("Adapter should contains no migration")
	}

	if migs, _ := adapter.GetAll(); len(migs) != 0 {
		t.Error("GetAll should returns an empty array")
	}

	if err := adapter.AddMigration("migration01.sql", time.Now()); err != nil {
		t.Error(err)
	}

	if err := adapter.AddMigration("migration02.sql", time.Now()); err != nil {
		t.Error(err)
	}

	migs, _ := adapter.GetAll()

	if len(migs) != 2 {
		t.Error("Should contains 2 migrations")
	}

	if err := adapter.RemoveMigration("migration01.sql"); err != nil {
		t.Error(err)
	}

	migs, _ = adapter.GetAll()

	if len(migs) != 1 {
		t.Error("Should contains 1 migration now")
	}
}
