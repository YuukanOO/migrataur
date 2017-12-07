package migrataur

import (
	"fmt"
	"testing"
	"time"
)

func TestMigrationMarshaling(t *testing.T) {
	assert := assert(t)

	migration := Migration{
		Name: "migration01",
		up:   "create table horses (name varchar(50) primary key);",
		down: "drop table horses;",
	}
	opts := DefaultMarshalOptions

	data, err := migration.marshal(opts)

	assert.nil(err)

	content := string(data)
	expected := fmt.Sprintf(`%s
%s
%s


%s
%s
%s`, opts.UpStart, migration.up, opts.UpEnd,
		opts.DownStart, migration.down, opts.DownEnd)

	assert.equals(expected, content)
}

func TestMigrationUnmarshaling(t *testing.T) {
	up := "create table horses (name varchar(50) primary key);"
	down := "drop table horses;"

	migration := Migration{
		Name: "migration02",
		up:   up,
		down: down,
	}

	data, _ := migration.marshal(DefaultMarshalOptions)

	migration.up = ""
	migration.down = ""

	if err := migration.unmarshal(data, DefaultMarshalOptions); err != nil {
		t.Error(err)
	}

	assert(t).
		equals(up, migration.up).
		equals(down, migration.down)
}

func TestMigrationToString(t *testing.T) {
	now := time.Now()
	appliedMigration := &Migration{Name: "migration01.sql", AppliedAt: &now}
	notAppliedMigration := &Migration{Name: "migration02.sql"}

	assert(t).
		true(appliedMigration.HasBeenApplied()).
		equals(fmt.Sprintf("[✓]\t%s", appliedMigration.Name), appliedMigration.String()).
		false(notAppliedMigration.HasBeenApplied()).
		equals(fmt.Sprintf("[ ]\t%s", notAppliedMigration.Name), notAppliedMigration.String())
}
