package migrataur

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Utility funcs used for testing

func cleanUpMigrationsDir() {
	fullpath, _ := filepath.Abs(DefaultOptions.Directory)

	if err := os.RemoveAll(fullpath); err != nil {
		panic(err)
	}
}

func TestGetRangeStr(t *testing.T) {
	assert := assert(t)

	first, last := getMigrationRange("")

	assert.
		equals("", first).
		equals("", last)

	first, last = getMigrationRange("migration01")

	assert.
		equals("migration01", first).
		equals("", last)

	first, last = getMigrationRange("migration02..migration07")

	assert.
		equals("migration02", first).
		equals("migration07", last)
}

func TestMigrataurInit(t *testing.T) {
	cleanUpMigrationsDir()

	instance := New(&mockAdapter{}, DefaultOptions)

	migration, err := instance.Init()

	assert(t).
		notNil(migration).
		nil(err).
		contains(DefaultOptions.InitialMigrationName, migration.Name).
		false(migration.HasBeenApplied()).
		equals(mockInitialUp, migration.up).
		equals(mockInitialDown, migration.down)
}

func TestMigrataurNew(t *testing.T) {
	cleanUpMigrationsDir()

	instance := New(&mockAdapter{}, DefaultOptions)

	migration, err := instance.NewMigration("migration01")

	assert(t).
		notNil(migration).
		nil(err).
		contains("migration01.sql", migration.Name)
}

func TestMigrataurMigrateToLatest(t *testing.T) {
	cleanUpMigrationsDir()

	instance := New(&mockAdapter{}, DefaultOptions)

	instance.NewMigration("migration01")
	instance.NewMigration("migration02")
	instance.NewMigration("migration03")
	instance.NewMigration("migration04")

	applied, err := instance.MigrateToLatest()

	assert(t).
		nil(err).
		equals(4, len(applied)).
		true(applied[0].IsInitial()).
		migrationsEquals(applied, "migration01", "migration02", "migration03", "migration04")
}

func TestMigrataurMigrate(t *testing.T) {
	cleanUpMigrationsDir()

	assert := assert(t)

	instance := New(&mockAdapter{}, DefaultOptions)

	instance.NewMigration("migration01")
	instance.NewMigration("migration02")
	instance.NewMigration("migration03")
	instance.NewMigration("migration04")
	instance.NewMigration("migration05")
	instance.NewMigration("migration06")

	applied, err := instance.Migrate("migration02..migration04")

	assert.
		nil(err).
		equals(3, len(applied)).
		migrationsEquals(applied, "migration02", "migration03", "migration04")

	applied, err = instance.Migrate("migration05")

	assert.
		nil(err).
		equals(1, len(applied)).
		migrationsEquals(applied, "migration05")

	applied, err = instance.Migrate("migration05")

	assert.
		nil(err).
		equals(0, len(applied))
}

func TestMigrataurGetAll(t *testing.T) {
	cleanUpMigrationsDir()

	assert := assert(t)
	instance := New(&mockAdapter{}, DefaultOptions)

	migrations, err := instance.GetAll()

	assert.
		nil(err).
		equals(0, len(migrations))

	instance.NewMigration("migration01")
	instance.NewMigration("migration02")
	instance.NewMigration("migration03")
	instance.NewMigration("migration04")

	instance.Migrate("migration01..migration02")
	instance.Migrate("migration04")

	migrations, err = instance.GetAll()

	assert.
		nil(err).
		equals(4, len(migrations))

	for _, m := range migrations {
		if strings.Contains(m.Name, "migration03") {
			if m.HasBeenApplied() {
				t.Errorf("%s should not have been applied", m.Name)
			}
		} else if !m.HasBeenApplied() {
			t.Errorf("%s should have been applied", m.Name)
		}
	}
}

func TestMigrataurRollback(t *testing.T) {
	cleanUpMigrationsDir()

	assert := assert(t)

	instance := New(&mockAdapter{}, DefaultOptions)

	instance.NewMigration("migration01")
	instance.NewMigration("migration02")
	instance.NewMigration("migration03")
	instance.NewMigration("migration04")
	instance.NewMigration("migration05")

	applied, err := instance.MigrateToLatest()

	assert.
		nil(err).
		equals(5, len(applied)).
		migrationsEquals(applied, "migration01", "migration02", "migration03", "migration04", "migration05")

	applied, err = instance.Rollback("migration05..migration03")

	assert.
		nil(err).
		equals(3, len(applied)).
		migrationsEquals(applied, "migration05", "migration04", "migration03")

	applied, err = instance.Rollback("migration02")

	assert.
		nil(err).
		equals(1, len(applied)).
		migrationsEquals(applied, "migration02")

	applied, err = instance.Rollback("migration02")

	assert.
		nil(err).
		equals(0, len(applied))
}

func TestMigrataurReset(t *testing.T) {
	cleanUpMigrationsDir()

	assert := assert(t)

	instance := New(&mockAdapter{}, DefaultOptions)

	instance.NewMigration("migration01")
	instance.NewMigration("migration02")
	instance.NewMigration("migration03")
	instance.NewMigration("migration04")

	applied, err := instance.MigrateToLatest()

	assert.
		nil(err).
		equals(4, len(applied)).
		migrationsEquals(applied, "migration01", "migration02", "migration03", "migration04")

	applied, err = instance.Reset()

	assert.
		nil(err).
		equals(4, len(applied)).
		true(applied[3].IsInitial()).
		migrationsEquals(applied, "migration04", "migration03", "migration02", "migration01")
}

func TestMigrationsSorting(t *testing.T) {
	assert := assert(t)

	migrations := []*Migration{
		&Migration{Name: "migration03"},
		&Migration{Name: "migration04"},
		&Migration{Name: "migration02"},
		&Migration{Name: "migration01"},
	}

	sortMigrations(migrations, dirUp)

	assert.migrationsEquals(migrations, "migration01", "migration02", "migration03", "migration04")

	sortMigrations(migrations, dirDown)

	assert.migrationsEquals(migrations, "migration04", "migration03", "migration02", "migration01")
}

func TestWithLoggerNil(t *testing.T) {
	instance := New(&mockAdapter{}, Options{
		Logger: nil,
	})

	instance.GetAll()

	assert(t).nil(instance.Options.Logger)
}
