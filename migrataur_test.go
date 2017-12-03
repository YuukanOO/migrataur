package migrataur

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// Utility funcs used for testing

func shouldHaveBeenEquals(t *testing.T, expected, actual interface{}) {
	t.Errorf("Expected: %s, Got: %s", expected, actual)
}

func cleanUpMigrationsDir() {
	fullpath, _ := filepath.Abs(DefaultOptions.Directory)

	if err := os.RemoveAll(fullpath); err != nil {
		panic(err)
	}
}

func TestGetRangeStr(t *testing.T) {
	first, last := getMigrationRange("")

	if first != "" || last != "" {
		t.Error("Start and end should be empty")
	}

	first, last = getMigrationRange("migration01")

	if first != "migration01" {
		shouldHaveBeenEquals(t, "migration01", first)
	}

	if last != "" {
		shouldHaveBeenEquals(t, "", last)
	}

	first, last = getMigrationRange("migration02..migration07")

	if first != "migration02" {
		shouldHaveBeenEquals(t, "migration02", first)
	}

	if last != "migration07" {
		shouldHaveBeenEquals(t, "migration07", last)
	}
}

func TestMigrataurInit(t *testing.T) {
	cleanUpMigrationsDir()

	instance := New(&mockAdapter{}, DefaultOptions)

	migration, err := instance.Init()

	if err != nil {
		t.Error(err)
	}

	if migration == nil {
		t.Error("Initial migration should not be nil")
	}

	if !strings.HasSuffix(migration.Name, "initMigrataur.sql") {
		shouldHaveBeenEquals(t, "initMigrataur.sql", migration.Name)
	}

	if migration.HasBeenApplied() {
		shouldHaveBeenEquals(t, false, migration.HasBeenApplied())
	}
}

func TestMigrataurNew(t *testing.T) {
	cleanUpMigrationsDir()

	instance := New(&mockAdapter{}, DefaultOptions)

	migration := instance.NewMigration("migration01")

	if !strings.HasSuffix(migration.Name, "migration01.sql") {
		t.Error("Migration name should contains migration01.sql")
	}
}

func TestMigrataurMigrateToLatest(t *testing.T) {
	cleanUpMigrationsDir()

	adapter := &mockAdapter{}
	instance := New(adapter, DefaultOptions)

	instance.NewMigration("migration01")
	instance.NewMigration("migration02")
	instance.NewMigration("migration03")
	instance.NewMigration("migration04")

	instance.MigrateToLatest()

	if len(adapter.appliedMigrations) != 4 {
		shouldHaveBeenEquals(t, 4, len(adapter.appliedMigrations))
	}
}

func TestMigrataurMigrate(t *testing.T) {
	cleanUpMigrationsDir()

	adapter := &mockAdapter{}
	instance := New(adapter, DefaultOptions)

	instance.NewMigration("migration01")
	instance.NewMigration("migration02")
	instance.NewMigration("migration03")
	instance.NewMigration("migration04")
	instance.NewMigration("migration05")
	instance.NewMigration("migration06")

	instance.Migrate("migration02..migration04")

	if len(adapter.appliedMigrations) != 3 {
		shouldHaveBeenEquals(t, 3, len(adapter.appliedMigrations))
	}

	instance.Migrate("migration05")

	if len(adapter.appliedMigrations) != 4 {
		shouldHaveBeenEquals(t, 4, len(adapter.appliedMigrations))
	}

	// Migrations count should not have changed
	instance.Migrate("migration05")

	if len(adapter.appliedMigrations) != 4 {
		shouldHaveBeenEquals(t, 4, len(adapter.appliedMigrations))
	}
}

func TestMigrataurGetAll(t *testing.T) {
	cleanUpMigrationsDir()

	instance := New(&mockAdapter{}, DefaultOptions)

	instance.NewMigration("migration01")
	instance.NewMigration("migration02")
	instance.NewMigration("migration03")
	instance.NewMigration("migration04")

	instance.Migrate("migration01..migration02")
	instance.Migrate("migration04")

	migrations := instance.GetAll()

	if len(migrations) != 4 {
		shouldHaveBeenEquals(t, 4, len(migrations))
	}

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

	adapter := &mockAdapter{}

	instance := New(adapter, DefaultOptions)

	instance.NewMigration("migration01")
	instance.NewMigration("migration02")
	instance.NewMigration("migration03")
	instance.NewMigration("migration04")
	instance.NewMigration("migration05")

	instance.MigrateToLatest()

	if len(adapter.appliedMigrations) != 5 {
		shouldHaveBeenEquals(t, 5, len(adapter.appliedMigrations))
	}

	instance.Rollback("migration05..migration03")

	if len(adapter.appliedMigrations) != 2 {
		shouldHaveBeenEquals(t, 2, len(adapter.appliedMigrations))
	}

	instance.Rollback("migration02")

	if len(adapter.appliedMigrations) != 1 {
		shouldHaveBeenEquals(t, 1, len(adapter.appliedMigrations))
	}

	// Twice should redo the down func
	instance.Rollback("migration02")

	if len(adapter.appliedMigrations) != 1 {
		shouldHaveBeenEquals(t, 1, len(adapter.appliedMigrations))
	}
}

func TestMigrataurReset(t *testing.T) {
	cleanUpMigrationsDir()

	adapter := &mockAdapter{}

	instance := New(adapter, DefaultOptions)

	instance.NewMigration("migration01")
	instance.NewMigration("migration02")
	instance.NewMigration("migration03")
	instance.NewMigration("migration04")

	instance.MigrateToLatest()

	if len(adapter.appliedMigrations) != 4 {
		shouldHaveBeenEquals(t, 4, len(adapter.appliedMigrations))
	}

	instance.Reset()

	if len(adapter.appliedMigrations) != 0 {
		shouldHaveBeenEquals(t, 0, len(adapter.appliedMigrations))
	}
}

func TestMigrationsSorting(t *testing.T) {
	migrations := []*Migration{
		&Migration{Name: "migration03"},
		&Migration{Name: "migration04"},
		&Migration{Name: "migration02"},
		&Migration{Name: "migration01"},
	}

	expected := []string{
		"migration01",
		"migration02",
		"migration03",
		"migration04",
	}

	sortMigrations(migrations, dirUp)

	for i := 0; i < len(expected); i++ {
		if migrations[i].Name != expected[i] {
			t.Errorf("Expecting %s, got %s when sorted", expected[i], migrations[i].Name)
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(expected)))

	sortMigrations(migrations, dirDown)

	for i := 0; i < len(expected); i++ {
		if migrations[i].Name != expected[i] {
			t.Errorf("Expecting %s, got %s when sorted", expected[i], migrations[i].Name)
		}
	}
}
