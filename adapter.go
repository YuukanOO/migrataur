package migrataur

import (
	"time"
)

// Adapter is the interface needed to access the underlying database
type Adapter interface {
	GetInitialMigration() *Migration
	AddMigration(completeName string, at time.Time) error
	RemoveMigration(completeName string) error
	Exec(command string) error
	GetAll() ([]*Migration, error)
}
