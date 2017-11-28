package sql

import (
	"database/sql"
	"time"

	"github.com/YuukanOO/migrataur"
)

// Adapter implements the interface defined by migrataur for common SQL databases
type Adapter struct {
	db *sql.DB
}

func NewAdapter(db *sql.DB) *Adapter {
	return &Adapter{
		db: db,
	}
}

func (a *Adapter) CreateMigrationsTableIfNotExists() error {
	a.db.Exec(`
create table __migrations(
	name varchar(250) not null,
	applied_at date not null
);
`)

	return nil
}

func (a *Adapter) AddMigration(completeName string, at time.Time) error {
	_, err := a.db.Exec("insert into __migrations values (?, ?)")

	return err
}

func (a *Adapter) RemoveMigration(completeName string) error {
	_, err := a.db.Exec("delete from __migration where name = ?")

	return err
}

func (a *Adapter) Exec(command string) error {
	_, err := a.db.Exec(command)

	return err
}

func (a *Adapter) GetAll() ([]*migrataur.Migration, error) {
	return nil, nil
}
