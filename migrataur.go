// Package migrataur is a simple migration tool for the Go language. It's written as
// a library that needs an adapter to work. It has been build with simplicity in mind.
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

// Those ones are used to define migration's directions
type dir int

const (
	// dirUp when applying migrations
	dirUp dir = iota
	// dirDown when rolling them back
	dirDown = iota
)

// Migrataur represents an instance configurated for a particular use.
// This is the main object you will use.
type Migrataur struct {
	Options Options
	adapter Adapter
}

// New instantiates a new Migrataur instance for the given options
func New(adapter Adapter, opts Options) *Migrataur {
	return &Migrataur{
		adapter: adapter,
		Options: opts.ExtendWith(DefaultOptions),
	}
}

// Init writes the initial migration provided by the adapter to create the needed
// migrations table, you should call it at the start of your project.
func (m *Migrataur) Init() (*Migration, error) {
	m.Printf("Initializing migrataur")

	fullPath := m.getMigrationFullpath(m.Options.InitialMigrationName)
	up, down := m.adapter.GetInitialMigration()

	initialMigration := &Migration{
		Name: filepath.Base(fullPath),
		up:   up,
		down: down,
	}

	if err := initialMigration.writeTo(fullPath, m.Options.MarshalOptions); err != nil {
		return nil, err
	}

	m.Printf("\t%s created!", initialMigration.Name)

	return initialMigration, nil
}

// NewMigration creates a new migration in the configured folder and returns the instance of the migration
// attached to the newly created file
func (m *Migrataur) NewMigration(name string) (*Migration, error) {
	m.Printf("Creating %s", name)

	fullPath := m.getMigrationFullpath(name)
	migration := &Migration{Name: filepath.Base(fullPath)}

	if err := migration.writeTo(fullPath, m.Options.MarshalOptions); err != nil {
		return nil, err
	}

	m.Printf("\t%s created!", migration.Name)

	return migration, nil
}

// GetAll retrieve all migrations for the current instance. It will list applied and pending migrations
func (m *Migrataur) GetAll() ([]*Migration, error) {
	m.Printf("Fetching migrations in %s", m.Options.Directory)

	return m.getAllMigrations(dirUp)
}

// Migrate migrates the database and returns an array of effectively applied migrations (it will
// not contains those that were already applied.
// rangeOrName can be the exact migration name or a range such as <migration>..<another migration name>
func (m *Migrataur) Migrate(rangeOrName string) ([]*Migration, error) {
	m.Printf("Applying %s", rangeOrName)

	return m.applyRange(rangeOrName, dirUp)
}

// MigrateToLatest migrates the database to the latest version
func (m *Migrataur) MigrateToLatest() ([]*Migration, error) {
	m.Printf("Applying all pending migrations")

	return m.applyAll(dirUp)
}

// Rollback inverts migrations and return an array of effectively rollbacked migrations
// (it will not contains those that were not applied).
func (m *Migrataur) Rollback(rangeOrName string) ([]*Migration, error) {
	m.Printf("Rollbacking %s", rangeOrName)

	return m.applyRange(rangeOrName, dirDown)
}

// Reset resets the database to its initial state
func (m *Migrataur) Reset() ([]*Migration, error) {
	m.Printf("Resetting database")

	return m.applyAll(dirDown)
}

// Printf logs a message using the provided Logger if any
func (m *Migrataur) Printf(format string, args ...interface{}) {
	if m.Options.Logger != nil {
		m.Options.Logger.Printf(format, args...)
	}
}

// Fatalf logs a fatal message using the provided Logger if any
func (m *Migrataur) Fatalf(format string, args ...interface{}) {
	if m.Options.Logger != nil {
		m.Options.Logger.Fatalf(format, args...)
	}
}

func (m *Migrataur) applyAll(direction dir) ([]*Migration, error) {
	migrations, err := m.getAllMigrations(direction)

	if err != nil {
		return nil, err
	}

	return m.apply(migrations, direction)
}

func (m *Migrataur) applyRange(rangeOrName string, direction dir) ([]*Migration, error) {
	start, end := getMigrationRange(rangeOrName)
	migrations, err := m.getAllMigrationsForRange(start, end, direction)

	if err != nil {
		return nil, err
	}

	return m.apply(migrations, direction)
}

// getAllFromFilesystem reads all migrations in the directory and instantiates them.
func (m *Migrataur) getAllFromFilesystem() ([]*Migration, error) {
	migrations := []*Migration{}
	files, err := ioutil.ReadDir(m.Options.Directory)

	if err != nil {
		pathErr, ok := err.(*os.PathError)

		if !ok || pathErr.Op != "open" {
			return nil, err
		}

		// If it's an "open" error, maybe that's because the directory does not exists yet
		return migrations, nil
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		existingMigration := &Migration{Name: f.Name()}

		data, err := ioutil.ReadFile(filepath.Join(m.Options.Directory, f.Name()))

		if err != nil {
			return nil, err
		}

		if err = existingMigration.unmarshal(data, m.Options.MarshalOptions); err != nil {
			return nil, err
		}

		migrations = append(migrations, existingMigration)
	}

	return migrations, err
}

// sortMigrations sorts given migrations by their name.
func sortMigrations(migrations []*Migration, direction dir) {
	if direction == dirUp {
		sort.Sort(byName(migrations))
	} else {
		sort.Sort(sort.Reverse(byName(migrations)))
	}
}

// apply given migrations in the given direction
func (m *Migrataur) apply(migrations []*Migration, direction dir) ([]*Migration, error) {
	appliedMigrations := []*Migration{}

	for _, mig := range migrations {
		ok, err := m.applyOne(mig, direction)

		if err != nil {
			return nil, err
		}

		if ok {
			appliedMigrations = append(appliedMigrations, mig)
		}
	}

	if len(appliedMigrations) == 0 {
		m.Printf("\tAll clear, nothing done!")
	}

	return appliedMigrations, nil
}

// getAllMigrationsForRange retrieves all migrations concerned by a range
func (m *Migrataur) getAllMigrationsForRange(start, end string, direction dir) ([]*Migration, error) {
	migrations, err := m.getAllMigrations(direction)

	if err != nil {
		return nil, err
	}

	rangeMigrations := []*Migration{}

	if start == "" {
		return rangeMigrations, nil
	}

	startApplied := false

	for _, mig := range migrations {
		if !startApplied {
			if strings.Contains(mig.Name, start) {
				startApplied = true
				rangeMigrations = append(rangeMigrations, mig)

				// Break early if no end migration has been set or if the end is the same
				if end == "" || start == end {
					break
				}
			}
		} else {
			rangeMigrations = append(rangeMigrations, mig)

			// If we reach the end, break
			if strings.Contains(mig.Name, end) {
				break
			}
		}
	}

	return rangeMigrations, nil
}

// getAllMigrations retrieves all migrations from the filesystem, and from the
// configurated adapter. It will mark them as applied if they are present in the
// adapter.
func (m *Migrataur) getAllMigrations(direction dir) ([]*Migration, error) {

	fileSystemMigrations, err := m.getAllFromFilesystem()

	if err != nil {
		return nil, err
	}

	adapterMigrations, err := m.adapter.GetAll()

	if err != nil {
		return nil, err
	}

	// Constructs the migrations map to easily update them with adapter ones
	migrationsMap := map[string]*Migration{}
	migrationsCount := len(fileSystemMigrations)

	for _, m := range fileSystemMigrations {
		migrationsMap[m.Name] = m
	}

	for _, mig := range adapterMigrations {
		fsMigration, ok := migrationsMap[mig.Name]

		if !ok {
			return nil, fmt.Errorf("the migration %s was not found in the migrations directory", mig.Name)
		}

		fsMigration.hasBeenAppliedAt(*mig.AppliedAt)
	}

	sortMigrations(fileSystemMigrations, direction)

	// Find the initial migration and marks it. This is used primarly by adapters to
	// perform specific behaviors
	if migrationsCount > 0 {
		switch direction {
		case dirUp:
			fileSystemMigrations[0].markAsInitial()
		case dirDown:
			fileSystemMigrations[migrationsCount-1].markAsInitial()
		}
	}

	return fileSystemMigrations, nil
}

func getMigrationRange(rangeStr string) (first, last string) {
	splitted := strings.Split(rangeStr, "..")

	if len(splitted) == 1 {
		return splitted[0], ""
	}

	return splitted[0], splitted[1]
}

// applyOne runs a single migration and returns if it has been applied. If the migration
// did not run because that was not needed, it will returns false.
func (m *Migrataur) applyOne(migration *Migration, direction dir) (bool, error) {

	// Do not execute commands if already applied or not applied at all when rolling back
	if (migration.HasBeenApplied() && direction == dirUp) || (!migration.HasBeenApplied() && direction == dirDown) {
		return false, nil
	}

	command := migration.up

	if direction == dirDown {
		command = migration.down
	}

	if err := m.adapter.Exec(command); err != nil {
		m.Fatalf("✗\t%s: %s", migration.Name, err)
		return false, err
	}

	if direction == dirUp {
		migration.hasBeenAppliedAt(time.Now().UTC())

		if err := m.adapter.MigrationApplied(migration); err != nil {
			m.Fatalf("✗\t%s: %s", migration.Name, err)
			return false, err
		}

	} else {
		migration.hasBeenRolledBack()

		if err := m.adapter.MigrationRollbacked(migration); err != nil {
			m.Fatalf("✗\t%s: %s", migration.Name, err)
			return false, err
		}
	}

	m.Printf("✓\t%s", migration.Name)

	return true, nil
}

func (m *Migrataur) getMigrationFullpath(name string) string {
	return filepath.Join(m.Options.Directory,
		fmt.Sprintf("%s_%s%s", m.Options.SequenceGenerator(), name, m.Options.Extension))
}
