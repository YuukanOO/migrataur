package migrataur

// Adapter is the interface needed to access the underlying database
type Adapter interface {
	CreateMigrationsTableIfNotExists() error
	Exec(command string) error
	GetAll() ([]*Migration, error)
}
