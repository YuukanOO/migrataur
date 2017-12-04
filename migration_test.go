package migrataur

import (
	"fmt"
	"testing"
)

func TestMigrationMarshaling(t *testing.T) {
	migration := Migration{
		Name: "migration01",
		Up:   "create table horses (name varchar(50) primary key);",
		Down: "drop table horses;",
	}
	opts := DefaultMarshalOptions

	data, err := migration.Marshal(opts)

	if err != nil {
		t.Error(err)
	}

	content := string(data)
	expected := fmt.Sprintf(`%s %s
%s
%s
%s


%s
%s
%s
`, opts.Header, migration.Name,
		opts.UpStart, migration.Up, opts.UpEnd,
		opts.DownStart, migration.Down, opts.DownEnd)

	assertEquals(t, expected, content)
}

func TestMigrationUnmarshaling(t *testing.T) {
	up := "create table horses (name varchar(50) primary key);"
	down := "drop table horses;"

	migration := Migration{
		Name: "migration02",
		Up:   up,
		Down: down,
	}

	data, _ := migration.Marshal(DefaultMarshalOptions)

	migration.Up = ""
	migration.Down = ""

	if err := migration.Unmarshal(data, DefaultMarshalOptions); err != nil {
		t.Error(err)
	}

	assertEquals(t, up, migration.Up)
	assertEquals(t, down, migration.Down)
}
