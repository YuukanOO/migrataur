package migrataur

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// Utility funcs used for testing

func assertEquals(t *testing.T, expected, actual interface{}) {
	if actual != expected {
		t.Errorf("Expected: %s, Got: %s", expected, actual)
	}
}

func assertNotEquals(t *testing.T, expected, actual interface{}) {
	if actual == expected {
		t.Errorf("Should be not equals: %s, And: %s", expected, actual)
	}
}

func assertEqualsWith(t *testing.T, predicate bool, expected, actual interface{}) {
	if !predicate {
		t.Errorf("Expected: %s, Got: %s", expected, actual)
	}
}

func assertNotNil(t *testing.T, actual interface{}) {
	if actual == nil {
		t.Error("Should not be nil!")
	}
}

func assertNil(t *testing.T, actual interface{}) {
	if actual != nil {
		t.Error(actual)
	}
}

func cleanUpMigrationsDir() {
	fullpath, _ := filepath.Abs(DefaultOptions.Directory)

	if err := os.RemoveAll(fullpath); err != nil {
		panic(err)
	}
}

func assertMigrationsEquals(t *testing.T, migrations []*Migration, names ...string) {
	lenActual, lenExpected := len(migrations), len(names)

	assertEquals(t, lenExpected, lenActual)

	for i := 0; i < len(migrations); i++ {
		assertEquals(t, names[i], migrations[i].Name)
	}
}

func TestGetRangeStr(t *testing.T) {
	first, last := getMigrationRange("")

	if first != "" || last != "" {
		t.Error("Start and end should be empty")
	}

	first, last = getMigrationRange("migration01")

	assertEquals(t, "migration01", first)
	assertEquals(t, "", last)

	first, last = getMigrationRange("migration02..migration07")

	assertEquals(t, "migration02", first)
	assertEquals(t, "migration07", last)
}

func TestMigrataurInit(t *testing.T) {
	cleanUpMigrationsDir()

	instance := New(&mockAdapter{}, DefaultOptions)

	migration, err := instance.Init()

	assertNotNil(t, migration)
	assertNil(t, err)
	assertEqualsWith(t, strings.HasSuffix(migration.Name, "initMigrataur.sql"), "initMigrataur.sql", migration.Name)
	assertEquals(t, false, migration.HasBeenApplied())
}

func TestMigrataurNew(t *testing.T) {
	cleanUpMigrationsDir()

	instance := New(&mockAdapter{}, DefaultOptions)

	migration, err := instance.NewMigration("migration01")

	assertNotNil(t, migration)
	assertNil(t, err)
	assertEqualsWith(t, strings.HasSuffix(migration.Name, "migration01.sql"), "migration01.sql", migration.Name)
}

func TestMigrataurMigrateToLatest(t *testing.T) {
	cleanUpMigrationsDir()

	adapter := &mockAdapter{}
	instance := New(adapter, DefaultOptions)

	instance.NewMigration("migration01")
	instance.NewMigration("migration02")
	instance.NewMigration("migration03")
	instance.NewMigration("migration04")

	_, err := instance.MigrateToLatest()

	assertNil(t, err)
	assertEquals(t, 4, len(adapter.appliedMigrations))
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

	_, err := instance.Migrate("migration02..migration04")

	assertNil(t, err)
	assertEquals(t, 3, len(adapter.appliedMigrations))

	_, err = instance.Migrate("migration05")

	assertNil(t, err)
	assertEquals(t, 4, len(adapter.appliedMigrations))

	// Migrations count should not have changed
	_, err = instance.Migrate("migration05")

	assertNil(t, err)
	assertEquals(t, 4, len(adapter.appliedMigrations))
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

	migrations, err := instance.GetAll()

	assertNil(t, err)
	assertEquals(t, 4, len(migrations))

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

	_, err := instance.MigrateToLatest()

	assertNil(t, err)
	assertEquals(t, 5, len(adapter.appliedMigrations))

	_, err = instance.Rollback("migration05..migration03")

	assertNil(t, err)
	assertEquals(t, 2, len(adapter.appliedMigrations))

	_, err = instance.Rollback("migration02")

	assertNil(t, err)
	assertEquals(t, 1, len(adapter.appliedMigrations))

	// Twice should redo the down func
	_, err = instance.Rollback("migration02")

	assertNil(t, err)
	assertEquals(t, 1, len(adapter.appliedMigrations))
}

func TestMigrataurReset(t *testing.T) {
	cleanUpMigrationsDir()

	adapter := &mockAdapter{}

	instance := New(adapter, DefaultOptions)

	instance.NewMigration("migration01")
	instance.NewMigration("migration02")
	instance.NewMigration("migration03")
	instance.NewMigration("migration04")

	_, err := instance.MigrateToLatest()

	assertNil(t, err)
	assertEquals(t, 4, len(adapter.appliedMigrations))

	_, err = instance.Reset()

	assertNil(t, err)
	assertEquals(t, 0, len(adapter.appliedMigrations))
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
