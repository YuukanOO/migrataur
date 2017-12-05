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

func (a *mockAdapter) GetInitialMigration(name string) *Migration {
	return &Migration{Name: name}
}

func (a *mockAdapter) AddMigration(name string, at time.Time) error {
	a.appliedMigrations = append(a.appliedMigrations, &Migration{
		Name:      name,
		AppliedAt: &at,
	})

	return nil
}

func (a *mockAdapter) RemoveMigration(name string) error {
	for i, m := range a.appliedMigrations {
		if m.Name == name {
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
	assert := assert(t)

	assert.equals(0, len(adapter.appliedMigrations))

	migrations, err := adapter.GetAll()

	assert.
		nil(err).
		equals(0, len(migrations)).
		nil(adapter.AddMigration("migration01.sql", time.Now())).
		nil(adapter.AddMigration("migration02.sql", time.Now()))

	migrations, err = adapter.GetAll()

	assert.
		nil(err).
		equals(2, len(migrations)).
		nil(adapter.RemoveMigration("migration01.sql"))

	migrations, err = adapter.GetAll()

	assert.
		nil(err).
		equals(1, len(migrations))
}
