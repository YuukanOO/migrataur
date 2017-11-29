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

func TestMigrationMarshaling(t *testing.T) {
	migration := Migration{
		Name: migrationName,
		Up:   upFixture,
		Down: downFixture,
	}

	data, err := migration.Marshal(DefaultMarshalOptions)

	if err != nil {
		t.Error(err)
	}

	content := string(data)

	if content != expectedSerializedContent {
		assertEquals(t, expectedSerializedContent, content)
	}
}

func TestMigrationUnmarshaling(t *testing.T) {
	migration := Migration{
		Name: migrationName,
		Up:   upFixture,
		Down: downFixture,
	}

	data, _ := migration.Marshal(DefaultMarshalOptions)

	migration.Up = ""
	migration.Down = ""

	if err := migration.Unmarshal(data, DefaultMarshalOptions); err != nil {
		t.Error(err)
	}

	if migration.Up != upFixture {
		assertEquals(t, upFixture, migration.Up)
	}

	if migration.Down != downFixture {
		assertEquals(t, downFixture, migration.Down)
	}
}
