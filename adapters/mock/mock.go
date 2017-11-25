package mock

import "github.com/YuukanOO/migrataur"
import "time"

// Adapter represents a bare in memory adapter. Used for testing purposes.
type Adapter struct{}

func (a *Adapter) CreateMigrationsTableIfNotExists() error {
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
