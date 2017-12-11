package migrataur

// Adapter is the interface needed to access the underlying database. This is where
// you should implements the desired behavior. Built-in adapters are found in the subpackage
// /adapters.
type Adapter interface {
	// GetInitialMigration retrieves the migration up and down code and is used to populate
	// the migrations history table.
	GetInitialMigration() (up, down string)
	// MigrationApplied is called when the migration has been successfully applied by the
	// adapter. This is where you should insert the migration in the history.
	MigrationApplied(migration *Migration) error
	// MigrationRollbacked is called when the migration has been successfully rolled back.
	// This is where you should remove the migration from the history.
	MigrationRollbacked(migration *Migration) error
	// Exec the given commands. This is call by Migrataur to apply or rollback a migration
	// with the corresponding code.
	Exec(command string) error
	// GetAll retrieves all migrations for this adapter
	GetAll() ([]*Migration, error)
}
