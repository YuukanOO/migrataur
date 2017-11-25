package mock

import (
	"fmt"
	"time"

	"github.com/YuukanOO/migrataur"
)

// Adapter represents a bare in memory adapter. Used for testing purposes.
type Adapter struct{}

func (a *Adapter) CreateMigrationsTableIfNotExists() error {
	return nil
}

func (a *Adapter) AddMigration(name string, at time.Time) error {
	fmt.Println("Adding migration " + name + " " + at.String())
	return nil
}

func (a *Adapter) RemoveMigration(name string) error {
	fmt.Println("Removing migration " + name)
	return nil
}

func (a *Adapter) Exec(command string) error {
	return nil
}

func (a *Adapter) GetAll() ([]*migrataur.Migration, error) {
	return []*migrataur.Migration{
		migrataur.NewAdapterMigration("1511644918_createUsers.sql", time.Now()),
	}, nil
}
