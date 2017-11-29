package migrataur

import (
	"fmt"
	"testing"
)

const (
	migrationName = "TestMigration.sql"
	upFixture     = `
create table my_resources (
	id int primary key,
	name varchar(50),
);
`
	downFixture = `
drop table my_resources;
`
)

var (
	expectedSerializedContent = fmt.Sprintf(`-- migration %s
-- +migrataur up
%s
-- -migrataur up


-- +migrataur down
%s
-- -migrataur down
`, migrationName, upFixture, downFixture)
)

func TestMigrationCanBeSerializedToText(t *testing.T) {
	migration := Migration{
		Name: migrationName,
		Up:   upFixture,
		Down: downFixture,
	}

	data, err := migration.MarshalText()

	if err != nil {
		t.Error(err)
	}

	content := string(data)

	if content != expectedSerializedContent {
		t.Errorf(`Content should be equals to
%s, was %s`, expectedSerializedContent, content)
	}
}

func TestMigrationCanBeDeserializedFromText(t *testing.T) {
	migration := Migration{
		Name: migrationName,
		Up:   upFixture,
		Down: downFixture,
	}

	data, _ := migration.MarshalText()

	migration.Up = ""
	migration.Down = ""

	if err := migration.UnmarshalText(data); err != nil {
		t.Error(err)
	}

	if migration.Up != upFixture {
		t.Errorf("Up migration should be equal to %s, was %s", upFixture, migration.Up)
	}

	if migration.Down != downFixture {
		t.Errorf("Down migration should be equal to %s, was %s", downFixture, migration.Down)
	}
}
