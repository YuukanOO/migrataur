package migrataur

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
		fmt.Sprintf("%s_%s%s", m.options.SequenceGenerator(), name, m.options.Extension))
}

// NewMigration creates a new migration in the configured folder and returns the instance of the migration
// attached to the newly created file
func (m *Migrataur) NewMigration(name string) *Migration {

	fullPath := m.getMigrationFullpath(name)
	migration := newMigration(filepath.Base(fullPath))

	file, err := os.Create(fullPath)

	if err != nil {
		m.options.Logger.Panic(err)
	}

	defer file.Close()

	data, err := migration.MarshalText()

	if err != nil {
		m.options.Logger.Panic(err)
	}

	_, err = file.Write(data)

	if err != nil {
		m.options.Logger.Panic(err)
	}

	return migration
}

func (m *Migrataur) getAllFromFilesystem() []*Migration {
	migrations := []*Migration{}
	files, err := ioutil.ReadDir(m.options.Directory)

	if err != nil {
		m.options.Logger.Panic(err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		existingMigration := newMigration(f.Name())

		data, err := ioutil.ReadFile(filepath.Join(m.options.Directory, f.Name()))

		if err != nil {
			m.options.Logger.Panic(err)
		}

		if err = existingMigration.UnmarshalText(data); err != nil {
			m.options.Logger.Panic(err)
		}

		migrations = append(migrations, existingMigration)
	}

	return migrations
}

func sortMigrations(migrations []*Migration) {
	sort.Sort(ByName(migrations))
}

// GetAll retrieve all migrations for the current instance. It will list applied and pending migrations
func (m *Migrataur) GetAll() []*Migration {

	fileSystemMigrations := m.getAllFromFilesystem()
	adapterMigrations, err := m.adapter.GetAll()

	if err != nil {
		m.options.Logger.Panic(err)
	}

	// Constructs the migrations map to easily update them with adapter ones
	migrationsMap := map[string]*Migration{}

	for _, m := range fileSystemMigrations {
		migrationsMap[m.name] = m
	}

	for _, mig := range adapterMigrations {
		fsMigration, ok := migrationsMap[mig.name]

		if !ok {
			m.options.Logger.Panicf("The migration %s was not found in the migrations directory!", mig.name)
		}

		fsMigration.hasBeenAppliedAt(*mig.appliedAt)
	}

	sortMigrations(fileSystemMigrations)

	return fileSystemMigrations
}

func getMigrationRange(rangeStr string) (start, end string) {
	if rangeStr == "" {
		return "", ""
	}

	splitted := strings.Split(rangeStr, "..")

	if len(splitted) == 1 {
		return splitted[0], ""
	}

	return splitted[0], splitted[1]
}

// Migrate migrates the database.
// rangeOrName can be the exact migration name or a range such as <migration>..<another migration name>
func (m *Migrataur) Migrate(rangeOrName string) {
	start, end := getMigrationRange(rangeOrName)

	startApplied := false

	for _, migration := range m.GetAll() {
		if !startApplied {
			if strings.Contains(migration.name, start) {
				m.applyMigration(migration)

				startApplied = true

				// Break early if no end migration has been set or if the end is the same
				if end == "" || start == end {
					break
				}
			}
		} else {
			m.applyMigration(migration)

			// If we reach the end, break
			if strings.Contains(migration.name, end) {
				break
			}
		}
	}
}

func (m *Migrataur) applyMigration(migration *Migration) {
	if migration.HasBeenApplied() {
		return
	}

	if err := m.adapter.Exec(migration.upStr); err != nil {
		m.options.Logger.Panicf("✗\t%s", err)
	}

	if err := m.adapter.AddMigration(migration.name, time.Now()); err != nil {
		m.options.Logger.Panicf("✗\t%s", err)
	}

	m.options.Logger.Printf("✓\t%s", migration.name)
}

// MigrateToLatest migrates the database to the latest version
func (m *Migrataur) MigrateToLatest() {
	for _, migration := range m.GetAll() {
		m.applyMigration(migration)
	}
}
