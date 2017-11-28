package sql

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/YuukanOO/migrataur"
)

const (
	// DefaultTableName represents the name of the migrations table
	DefaultTableName = "__migrations"
	// DefaultPlaceholder holds the default value for the sql placeholder
	DefaultPlaceholder = "?"
)

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

func (a *Adapter) CreateMigrationsTableIfNotExists() error {
	a.db.Exec(fmt.Sprintf(`
create table %s(
	name varchar(250) primary key,
	applied_at timestamp not null
);
`, a.tableName))

	return nil
}

func (a *Adapter) AddMigration(completeName string, at time.Time) error {
	_, err := a.db.Exec(fmt.Sprintf("insert into %s values (%s, %s)", a.tableName, a.getPlaceholder(1), a.getPlaceholder(2)), completeName, at)

	return err
}

func (a *Adapter) RemoveMigration(completeName string) error {
	_, err := a.db.Exec(fmt.Sprintf("delete from %s where name = %s", a.tableName, a.getPlaceholder(1)), completeName)

	return err
}

func (a *Adapter) Exec(command string) error {
	_, err := a.db.Exec(command)

	return err
}

func (a *Adapter) GetAll() ([]*migrataur.Migration, error) {
	rows, err := a.db.Query(fmt.Sprintf("select name, applied_at from %s order by name", a.tableName))

	if err != nil {
		return nil, err
	}

	migrations := []*migrataur.Migration{}

	defer rows.Close()
	for rows.Next() {
		var (
			name      string
			appliedAt time.Time
		)

		if err = rows.Scan(&name, &appliedAt); err != nil {
			panic(err)
		}

		migrations = append(migrations, migrataur.NewMigration(name, appliedAt))
	}

	return migrations, nil
}
