// Package sql implements a generic adapter for SQL databases. It uses only
// the database/sql standard package.
package sql

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/YuukanOO/migrataur"
)

// DefaultTableName represents the default name of the migrations table
const DefaultTableName = "__migrations"

// DefaultPlaceholder holds the default value for the sql placeholder
// If you're on postgres, you should use ${i} where {i} will be replaced
// by the arg position.
const DefaultPlaceholder = "?"

// Adapter implements the interface defined by migrataur for common SQL databases
type Adapter struct {
	tableName   string
	placeholder string
	db          *sql.DB
}

// WithDB constructs a sql adapter with the given DB handle.
// It will use the default tableName "__migrations" and "?" placeholder.
func WithDB(db *sql.DB) *Adapter {
	return &Adapter{
		db:          db,
		tableName:   DefaultTableName,
		placeholder: DefaultPlaceholder,
	}
}

// WithDBAndOptions constructs a sql adapter with the given DB handle and
// options. table is the the name of the migrations table and placeholder is DB
// dependent. If you're on postgres, you should use ${i} where {i} will be replaced
// by the arg position.
func WithDBAndOptions(db *sql.DB, table, placeholder string) *Adapter {
	adapter := WithDB(db)

	adapter.tableName = table
	adapter.placeholder = placeholder

	return adapter
}

func (a *Adapter) getPlaceholder(idx int) string {
	return strings.Replace(a.placeholder, "{i}", strconv.Itoa(idx), -1)
}

func (a *Adapter) GetInitialMigration() (up, down string) {
	return fmt.Sprintf(`-- Do not edit this migration unless you know what you're doing!
create table %s(
	name varchar(250) primary key,
	applied_at timestamp not null
);`, a.tableName), fmt.Sprintf(`-- Warning also apply to this section ;)
drop table %s;`, a.tableName)
}

func (a *Adapter) MigrationApplied(migration *migrataur.Migration) error {
	_, err := a.db.Exec(fmt.Sprintf("insert into %s values (%s, %s)", a.tableName, a.getPlaceholder(1), a.getPlaceholder(2)), migration.Name, *migration.AppliedAt)

	return err
}

func (a *Adapter) MigrationRollbacked(migration *migrataur.Migration) error {
	if migration.IsInitial() {
		return nil
	}

	_, err := a.db.Exec(fmt.Sprintf("delete from %s where name = %s", a.tableName, a.getPlaceholder(1)), migration.Name)

	return err
}

func (a *Adapter) Exec(command string) error {
	_, err := a.db.Exec(command)

	return err
}

func (a *Adapter) GetAll() ([]*migrataur.Migration, error) {
	// If the database has not been initialized, the migration table doesn't exist yet
	// so fail silently for now
	rows, err := a.db.Query(fmt.Sprintf("select name, applied_at from %s order by name", a.tableName))

	migrations := []*migrataur.Migration{}

	if err != nil {
		return migrations, nil
	}

	defer rows.Close()

	for rows.Next() {
		var migration = &migrataur.Migration{}

		if err = rows.Scan(&migration.Name, &migration.AppliedAt); err != nil {
			panic(err)
		}

		migrations = append(migrations, migration)
	}

	return migrations, nil
}
