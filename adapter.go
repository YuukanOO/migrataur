package migrataur

import (
	"time"
)

// Adapter is the interface needed to access the underlying database
type Adapter interface {
	GetInitialMigration(name string) *Migration
	AddMigration(completeName string, at time.Time) error
	RemoveMigration(completeName string) error
	Exec(command string) error
	GetAll() ([]*Migration, error)
}
