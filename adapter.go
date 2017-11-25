package migrataur

import (
	"time"
)

// Adapter is the interface needed to access the underlying database
type Adapter interface {
	CreateMigrationsTableIfNotExists() error
	AddMigration(name string, at time.Time) error
	RemoveMigration(name string) error
	Exec(command string) error
	GetAll() ([]*Migration, error)
}
