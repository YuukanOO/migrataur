package migrataur

import (
	"fmt"
	"os"
	"path/filepath"
)

// Migrataur represents an instance configurated for a particular use
type Migrataur struct {
	options *Options
}

// New instantiates a new Migrataur instance for the given options
func New(opts *Options) *Migrataur {
	return &Migrataur{options: extendOptions(opts)}
}

func (m *Migrataur) getMigrationFullpath(name string) string {
	return filepath.Join(m.options.Directory,
		fmt.Sprintf("%s_%s%s", m.options.UnicityGenerator(), name, m.options.Extension))
}

// NewMigration creates a new migration in the configured folder and returns the instance of the migration
// attached to the newly created file
func (m *Migrataur) NewMigration(name string) *Migration {

	fullPath := m.getMigrationFullpath(name)

	migration := newMigration(filepath.Base(fullPath))

	file, err := os.Create(fullPath)

	if err != nil {
		panic(err)
	}

	migrationData, err := migration.MarshalText()

	if err != nil {
		panic(err)
	}

	_, err = file.Write(migrationData)

	if err != nil {
		panic(err)
	}

	return migration
}

// GetAll retrieve all migrations for the current instance. It will list applied and pending migrations
func (m *Migrataur) GetAll() []*Migration {
	return []*Migration{}
}
