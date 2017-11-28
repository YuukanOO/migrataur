package migrataur

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetRangeStr(t *testing.T) {
	start, end := getMigrationRange("")

	if start != "" || end != "" {
		t.Error("Start and end should be empty")
	}

	start, end = getMigrationRange("migration01")

	if start != "migration01" {
		t.Error("Start should be equals to migration01")
	}

	if end != "" {
		t.Error("End should be empty")
	}

	start, end = getMigrationRange("migration02..migration07")

	if start != "migration02" {
		t.Error("Start should be equals to migration02")
	}

	if end != "migration07" {
		t.Error("End should be equals to migration07")
	}
}

func TestMigrataur(t *testing.T) {
	// Ok that's a big ugly test and it writes files to disk, I should try to
	// refactor it someday

	fullpath, _ := filepath.Abs("./migrations")

	if err := os.RemoveAll(fullpath); err != nil {
		t.Error(err)
	}

	inst := New(newMockAdapter(), &Options{
		Extension: ".sql",
	})

	migration := inst.NewMigration("migration01")

	if migration == nil {
		t.Error("Migration should be set")
	}

	if !strings.HasSuffix(migration.name, "migration01.sql") {
		t.Error("New migration should have name and extension")
	}

	inst.NewMigration("migration02")
	inst.NewMigration("migration03")
	inst.NewMigration("migration04")
	inst.NewMigration("migration05")

	migrations := inst.GetAll()

	if len(migrations) != 5 {
		t.Error("We should have 5 migrations now")
	}

	for _, m := range migrations {
		if m.HasBeenApplied() {
			t.Error("All migrations should be pending")
		}
	}

	inst.Migrate("migration02..migration05")

	migrations = inst.GetAll()

	for _, m := range migrations {
		if m.name != "migration01" && !m.HasBeenApplied() {
			t.Error("Migrations should be applied")
		}
	}

	inst.Migrate("migration01")

	migrations = inst.GetAll()

	for _, m := range migrations {
		if !m.HasBeenApplied() {
			t.Error("Migrations should be applied")
		}
	}
}
