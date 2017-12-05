package migrataur

import (
	"time"
)

// Adapter is the interface needed to access the underlying database. This is where
// you should implements the desired behavior. Built-in adapters are found in the subpackage
// /adapters.
type Adapter interface {
	GetInitialMigration(name string) *Migration
	AddMigration(completeName string, at time.Time) error
	RemoveMigration(completeName string) error
	Exec(command string) error
	GetAll() ([]*Migration, error)
}
