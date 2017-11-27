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
		name:    migrationName,
		upStr:   upFixture,
		downStr: downFixture,
	}

	data, err := migration.MarshalText()

	if err != nil {
		t.Error(err)
	}

	content := string(data)

	if content != expectedSerializedContent {
		t.Errorf(`Content should be equals to
%s`, expectedSerializedContent)
	}
}

func TestMigrationCanBeDeserializedFromText(t *testing.T) {
	migration := Migration{
		name: migrationName,
	}

	data, _ := migration.MarshalText()

	if err := migration.UnmarshalText(data); err != nil {
		t.Error(err)
	}
}
