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

func (a *mockAdapter) AddMigration(migration *Migration) error {
	a.appliedMigrations = append(a.appliedMigrations, migration)

	return nil
}

func (a *mockAdapter) RemoveMigration(migration *Migration) error {
	for i, m := range a.appliedMigrations {
		if m.Name == migration.Name {
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

	now := time.Now()

	mig1 := &Migration{Name: "migration01.sql", AppliedAt: &now}
	mig2 := &Migration{Name: "migration02.sql", AppliedAt: &now}

	assert.
		nil(err).
		equals(0, len(migrations)).
		nil(adapter.AddMigration(mig1)).
		nil(adapter.AddMigration(mig2))

	migrations, err = adapter.GetAll()

	assert.
		nil(err).
		equals(2, len(migrations)).
		nil(adapter.RemoveMigration(mig1))

	migrations, err = adapter.GetAll()

	assert.
		nil(err).
		equals(1, len(migrations))
}
