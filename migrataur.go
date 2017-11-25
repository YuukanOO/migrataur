package migrataur

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// Migrataur represents an instance configurated for a particular use
type Migrataur struct {
	options *Options
	adapter Adapter
}

// New instantiates a new Migrataur instance for the given options
func New(adapter Adapter, opts *Options) *Migrataur {
	return &Migrataur{
		adapter: adapter,
		options: extendOptionsAndSanitize(opts),
	}
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

	defer file.Close()

	data, err := migration.MarshalText()

	if err != nil {
		panic(err)
	}

	_, err = file.Write(data)

	if err != nil {
		panic(err)
	}

	return migration
}

func (m *Migrataur) getAllFromFilesystem() []*Migration {
	migrations := []*Migration{}
	files, err := ioutil.ReadDir(m.options.Directory)

	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		existingMigration := newMigration(f.Name())

		data, err := ioutil.ReadFile(filepath.Join(m.options.Directory, f.Name()))

		if err != nil {
			panic(err)
		}

		if err = existingMigration.UnmarshalText(data); err != nil {
			panic(err)
		}

		migrations = append(migrations, existingMigration)
	}

	return migrations
}

// GetAll retrieve all migrations for the current instance. It will list applied and pending migrations
func (m *Migrataur) GetAll() []*Migration {

	fileSystemMigrations := m.getAllFromFilesystem()
	adapterMigrations, err := m.adapter.GetAll()

	if err != nil {
		panic(err)
	}

	// Constructs the migrations map to easily update them with adapter ones
	migrationsMap := map[string]*Migration{}

	for _, m := range fileSystemMigrations {
		migrationsMap[m.name] = m
	}

	for _, m := range adapterMigrations {
		fsMigration, ok := migrationsMap[m.name]

		if !ok {
			panic(fmt.Sprintf("The migration %s was not found in the migrations directory!", m.name))
		}

		fsMigration.hasBeenAppliedAt(*m.appliedAt)
	}

	return fileSystemMigrations
}

// MigrateToLatest migrates the database to the latest version
func (m *Migrataur) MigrateToLatest() {
	for _, migration := range m.GetAll() {
		if migration.HasBeenApplied() {
			continue
		}

		if err := m.adapter.Exec(migration.upStr); err != nil {
			panic(err)
		}

		if err := m.adapter.AddMigration(migration.name, time.Now()); err != nil {
			panic(err)
		}
	}
}
