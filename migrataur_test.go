package migrataur

import (
	"strings"
	"testing"
)

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

func TestGetAllMigrationsForRange(t *testing.T) {

	mockFSAdapter.hasFiles(
		mockFileInfo{name: "migration01.sql"},
		mockFileInfo{name: "migration02.sql"},
		mockFileInfo{name: "migration03.sql"},
		mockFileInfo{name: "migration04.sql"},
		mockFileInfo{name: "migration05.sql"},
	)

	assert := assert(t)
	instance := New(&mockAdapter{}, DefaultOptions)

	migrations, err := instance.getAllMigrationsForRange("", "", dirUp)

	assert.
		nil(err).
		equals(0, len(migrations))

	_, err = instance.getAllMigrationsForRange("doesnotexists", "", dirUp)

	assert.
		notNil(err)

	_, err = instance.getAllMigrationsForRange("migration01", "doesnotexists", dirUp)

	assert.
		notNil(err)

	migrations, err = instance.getAllMigrationsForRange("migration01", "", dirUp)

	assert.
		nil(err).
		equals(1, len(migrations)).
		applied(migrations, "migration01")

	migrations, err = instance.getAllMigrationsForRange("migration03", "migration05", dirUp)

	assert.
		nil(err).
		equals(3, len(migrations)).
		applied(migrations, "migration03", "migration04", "migration05")

	migrations, err = instance.getAllMigrationsForRange("migration05", "", dirDown)

	assert.
		nil(err).
		equals(1, len(migrations)).
		applied(migrations, "migration05")

	migrations, err = instance.getAllMigrationsForRange("migration05", "migration02", dirDown)

	assert.
		nil(err).
		equals(4, len(migrations)).
		applied(migrations, "migration05", "migration04", "migration03", "migration02")
}

func TestMigrataurInit(t *testing.T) {
	mockFSAdapter.empty()

	instance := New(&mockAdapter{}, DefaultOptions)
	migration, err := instance.Init()

	assert(t).
		notNil(migration).
		nil(err).
		contains(DefaultOptions.InitialMigrationName, migration.Name).
		false(migration.HasBeenApplied()).
		equals(mockInitialUp, migration.up).
		equals(mockInitialDown, migration.down).
		exists(migration.Name)
}

func TestMigrataurNew(t *testing.T) {
	mockFSAdapter.empty()

	instance := New(&mockAdapter{}, DefaultOptions)
	migration, err := instance.New("migration01")

	assert(t).
		notNil(migration).
		nil(err).
		contains("migration01.sql", migration.Name).
		exists(migration.Name)
}

func TestMigrataurRemove(t *testing.T) {
	mockFSAdapter.hasFiles(
		mockFileInfo{name: "migration01.sql"},
		mockFileInfo{name: "migration02.sql"},
		mockFileInfo{name: "migration03.sql"},
		mockFileInfo{name: "migration04.sql"},
		mockFileInfo{name: "migration05.sql"},
		mockFileInfo{name: "migration06.sql"},
	)

	assert := assert(t)
	instance := New(&mockAdapter{}, DefaultOptions)

	_, err := instance.MigrateToLatest()

	assert.
		nil(err).
		exists("migration03.sql")

	migrations, err := instance.Remove("")

	assert.
		nil(err).
		equals(0, len(migrations))

	migrations, err = instance.Remove("migration03")

	assert.
		nil(err).
		equals(1, len(migrations)).
		applied(migrations, "migration03").
		notExists("migration03.sql")

	migrations, err = instance.Remove("migration02..migration01")

	assert.
		nil(err).
		equals(2, len(migrations)).
		applied(migrations, "migration02", "migration01")

	_, err = instance.Remove("migration04..migration06")

	assert.notNil(err)
}

func TestMigrataurMigrateToLatest(t *testing.T) {
	mockFSAdapter.hasFiles(
		mockFileInfo{name: "migration01.sql"},
		mockFileInfo{name: "migration02.sql"},
		mockFileInfo{name: "migration03.sql"},
		mockFileInfo{name: "migration04.sql"},
	)

	instance := New(&mockAdapter{}, DefaultOptions)
	applied, err := instance.MigrateToLatest()

	assert(t).
		nil(err).
		equals(4, len(applied)).
		true(applied[0].IsInitial()).
		applied(applied, "migration01", "migration02", "migration03", "migration04")

	applied, err = instance.MigrateToLatest()

	assert(t).
		nil(err).
		equals(0, len(applied))
}

func TestMigrataurMigrate(t *testing.T) {
	mockFSAdapter.hasFiles(
		mockFileInfo{name: "migration01.sql"},
		mockFileInfo{name: "migration02.sql"},
		mockFileInfo{name: "migration03.sql"},
		mockFileInfo{name: "migration04.sql"},
		mockFileInfo{name: "migration05.sql"},
		mockFileInfo{name: "migration06.sql"},
	)

	assert := assert(t)

	instance := New(&mockAdapter{}, DefaultOptions)
	applied, err := instance.Migrate("migration02..migration04")

	assert.
		nil(err).
		equals(3, len(applied)).
		applied(applied, "migration02", "migration03", "migration04")

	applied, err = instance.Migrate("migration05")

	assert.
		nil(err).
		equals(1, len(applied)).
		applied(applied, "migration05")

	applied, err = instance.Migrate("migration05")

	assert.
		nil(err).
		equals(0, len(applied))

	_, err = instance.Migrate("doesnotexists")

	assert.notNil(err)
}

func TestMigrataurGetAll(t *testing.T) {
	mockFSAdapter.hasFiles()

	assert := assert(t)
	instance := New(&mockAdapter{}, DefaultOptions)

	migrations, err := instance.GetAll()

	assert.
		nil(err).
		equals(0, len(migrations))

	mockFSAdapter.hasFiles(
		mockFileInfo{name: "migration01.sql"},
		mockFileInfo{name: "migration02.sql"},
		mockFileInfo{name: "migration03.sql"},
		mockFileInfo{name: "migration04.sql"},
	)

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
	mockFSAdapter.hasFiles(
		mockFileInfo{name: "migration01.sql"},
		mockFileInfo{name: "migration02.sql"},
		mockFileInfo{name: "migration03.sql"},
		mockFileInfo{name: "migration04.sql"},
		mockFileInfo{name: "migration05.sql"},
	)

	assert := assert(t)
	instance := New(&mockAdapter{}, DefaultOptions)
	applied, err := instance.MigrateToLatest()

	assert.
		nil(err).
		equals(5, len(applied)).
		applied(applied, "migration01", "migration02", "migration03", "migration04", "migration05")

	applied, err = instance.Rollback("migration05..migration03")

	assert.
		nil(err).
		equals(3, len(applied)).
		applied(applied, "migration05", "migration04", "migration03")

	applied, err = instance.Rollback("migration02")

	assert.
		nil(err).
		equals(1, len(applied)).
		applied(applied, "migration02")

	applied, err = instance.Rollback("migration02")

	assert.
		nil(err).
		equals(0, len(applied))

	_, err = instance.Rollback("doesnotexists")

	assert.notNil(err)
}

func TestMigrataurReset(t *testing.T) {
	mockFSAdapter.hasFiles(
		mockFileInfo{name: "migration01.sql"},
		mockFileInfo{name: "migration02.sql"},
		mockFileInfo{name: "migration03.sql"},
		mockFileInfo{name: "migration04.sql"},
	)

	assert := assert(t)
	instance := New(&mockAdapter{}, DefaultOptions)
	applied, err := instance.MigrateToLatest()

	assert.
		nil(err).
		equals(4, len(applied)).
		applied(applied, "migration01", "migration02", "migration03", "migration04")

	applied, err = instance.Reset()

	assert.
		nil(err).
		equals(4, len(applied)).
		true(applied[3].IsInitial()).
		applied(applied, "migration04", "migration03", "migration02", "migration01")

	applied, err = instance.Reset()

	assert.
		nil(err).
		equals(0, len(applied))
}

func TestMigrationsSorting(t *testing.T) {
	assert := assert(t)

	migrations := []*Migration{
		{Name: "migration03"},
		{Name: "migration04"},
		{Name: "migration02"},
		{Name: "migration01"},
	}

	sortMigrations(migrations, dirUp)

	assert.applied(migrations, "migration01", "migration02", "migration03", "migration04")

	sortMigrations(migrations, dirDown)

	assert.applied(migrations, "migration04", "migration03", "migration02", "migration01")
}

func TestWithLoggerNilShouldNotPanic(t *testing.T) {
	instance := New(&mockAdapter{}, Options{
		Logger: nil,
	})

	instance.GetAll()

	assert(t).nil(instance.options.Logger)
}
