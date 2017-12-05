// Package migrataur is a simple migration tool for the Go language. It's written as
// a library that needs an adapter to work. It has been build with simplicity in mind.
package migrataur

import (
	"fmt"
	"io/ioutil"
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

func (m *Migrataur) getMigrationFullpath(name string) string {
	return filepath.Join(m.Options.Directory,
		fmt.Sprintf("%s_%s%s", m.Options.SequenceGenerator(), name, m.Options.Extension))
}

// Init writes the initial migration provided by the adapter to create the needed
// migrations table, you should call it at the start of your project.
func (m *Migrataur) Init() (*Migration, error) {

	fullPath := m.getMigrationFullpath(m.Options.InitialMigrationName)
	initialMigration := m.adapter.GetInitialMigration(filepath.Base(fullPath))

	if err := initialMigration.WriteTo(fullPath, m.Options.MarshalOptions); err != nil {
		return nil, err
	}

	m.Options.Logger.Printf("Migrataur initialized with %s", initialMigration.Name)

	return initialMigration, nil
}

// NewMigration creates a new migration in the configured folder and returns the instance of the migration
// attached to the newly created file
func (m *Migrataur) NewMigration(name string) (*Migration, error) {
	fullPath := m.getMigrationFullpath(name)
	migration := &Migration{Name: filepath.Base(fullPath)}

	if err := migration.WriteTo(fullPath, m.Options.MarshalOptions); err != nil {
		return nil, err
	}

	m.Options.Logger.Printf("%s created", migration.Name)

	return migration, nil
}

// getAllFromFilesystem reads all migrations in the directory and instantiates them.
func (m *Migrataur) getAllFromFilesystem() ([]*Migration, error) {
	files, err := ioutil.ReadDir(m.Options.Directory)

	if err != nil {
		return nil, err
	}

	migrations := []*Migration{}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		existingMigration := &Migration{Name: f.Name()}

		data, err := ioutil.ReadFile(filepath.Join(m.Options.Directory, f.Name()))

		if err != nil {
			return nil, err
		}

		if err = existingMigration.Unmarshal(data, m.Options.MarshalOptions); err != nil {
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

	return fileSystemMigrations, nil
}

func getMigrationRange(rangeStr string) (first, last string) {
	splitted := strings.Split(rangeStr, "..")

	if len(splitted) == 1 {
		return splitted[0], ""
	}

	return splitted[0], splitted[1]
}

func (m *Migrataur) runRange(start, end string, direction dir) ([]*Migration, error) {
	appliedMigrations := []*Migration{}

	if start == "" {
		return appliedMigrations, nil
	}

	startApplied := false

	migrations, err := m.getAllMigrations(direction)

	if err != nil {
		return nil, err
	}

	for _, migration := range migrations {
		if !startApplied {
			if strings.Contains(migration.Name, start) {

				ok, err := m.runStep(migration, direction)

				if err != nil {
					return nil, err
				}

				if ok {
					appliedMigrations = append(appliedMigrations, migration)
				}

				startApplied = true

				// Break early if no end migration has been set or if the end is the same
				if end == "" || start == end {
					break
				}
			}
		} else {
			ok, err := m.runStep(migration, direction)

			if err != nil {
				return nil, err
			}

			if ok {
				appliedMigrations = append(appliedMigrations, migration)
			}

			// If we reach the end, break
			if strings.Contains(migration.Name, end) {
				break
			}
		}
	}

	return appliedMigrations, nil
}

func (m *Migrataur) run(rangeOrName string, direction dir) ([]*Migration, error) {
	start, end := getMigrationRange(rangeOrName)

	return m.runRange(start, end, direction)
}

// runStep runs a single migration and returns if it has been applied. If the migration
// did not run because that was not needed, it will returns false.
func (m *Migrataur) runStep(migration *Migration, direction dir) (bool, error) {

	shouldSkip := false

	// Do not execute commands if already applied or not applied at all when rolling back
	if (migration.HasBeenApplied() && direction == dirUp) || (!migration.HasBeenApplied() && direction == dirDown) {
		shouldSkip = true
	}

	if shouldSkip {
		m.Options.Logger.Printf("—\t%s", migration.Name)

		return false, nil
	}

	command := migration.Up

	if direction == dirDown {
		command = migration.Down
	}

	if err := m.adapter.Exec(command); err != nil {
		m.Options.Logger.Fatalf("✗\t%s: %s", migration.Name, err)
		return false, err
	}

	if direction == dirUp {
		now := time.Now().UTC()

		if err := m.adapter.AddMigration(migration.Name, now); err != nil {
			m.Options.Logger.Fatalf("✗\t%s: %s", migration.Name, err)
			return false, err
		}

		migration.hasBeenAppliedAt(now)
	} else {
		if err := m.adapter.RemoveMigration(migration.Name); err != nil {
			m.Options.Logger.Fatalf("✗\t%s: %s", migration.Name, err)
			return false, err
		}

		migration.hasBeenRolledBack()
	}

	m.Options.Logger.Printf("✓\t%s", migration.Name)

	return true, nil
}

// GetAll retrieve all migrations for the current instance. It will list applied and pending migrations
func (m *Migrataur) GetAll() ([]*Migration, error) {
	m.Options.Logger.Printf("Fetching migrations in %s", m.Options.Directory)

	return m.getAllMigrations(dirUp)
}

// Migrate migrates the database and returns an array of effectively applied migrations (it will
// not contains those that were already applied.
// rangeOrName can be the exact migration name or a range such as <migration>..<another migration name>
func (m *Migrataur) Migrate(rangeOrName string) ([]*Migration, error) {
	m.Options.Logger.Printf("Applying %s", rangeOrName)

	return m.run(rangeOrName, dirUp)
}

// MigrateToLatest migrates the database to the latest version
func (m *Migrataur) MigrateToLatest() ([]*Migration, error) {
	m.Options.Logger.Print("Applying all pending migrations")

	migrations, err := m.getAllMigrations(dirUp)

	if err != nil {
		return nil, err
	}

	if len(migrations) == 0 {
		return []*Migration{}, nil
	}

	return m.runRange(migrations[0].Name, migrations[len(migrations)-1].Name, dirUp)
}

// Rollback inverts migrations and return an array of effectively rollbacked migrations
// (it will not contains those that were not applied).
func (m *Migrataur) Rollback(rangeOrName string) ([]*Migration, error) {
	m.Options.Logger.Printf("Rollbacking %s", rangeOrName)

	return m.run(rangeOrName, dirDown)
}

// Reset resets the database to its initial state
func (m *Migrataur) Reset() ([]*Migration, error) {
	m.Options.Logger.Print("Resetting database")

	migrations, err := m.getAllMigrations(dirDown)

	if err != nil {
		return nil, err
	}

	if len(migrations) == 0 {
		return []*Migration{}, nil
	}

	return m.runRange(migrations[0].Name, migrations[len(migrations)-1].Name, dirDown)
}
