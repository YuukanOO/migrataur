package migrataur

// Adapter is the interface needed to access the underlying database. This is where
// you should implements the desired behavior. Built-in adapters are found in the subpackage
// /adapters.
type Adapter interface {
	// GetInitialMigration retrieves the migration needed to create the migration table
	GetInitialMigration(name string) *Migration
	// AddMigration adds the given migration to the adapter history. This is where you
	// should insert the migration in the history and not where you should run the migration.
	AddMigration(migration *Migration) error
	// RemoveMigration removes the given migration from the adapter history. This is
	// where you should remove the migration from the history and not where you should
	// run the migration.
	RemoveMigration(migration *Migration) error
	// Exec the given commands. This is call by Migrataur to apply or rollback a migration
	// with the corresponding code.
	Exec(command string) error
	// GetAll retrieves all migrations for this adapter
	GetAll() ([]*Migration, error)
}
