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
	dirUp   dir = iota
	dirDown     = iota
)

// Migrataur represents an instance configurated for a particular use
type Migrataur struct {
	options Options
	adapter Adapter
}

// New instantiates a new Migrataur instance for the given options
func New(adapter Adapter, opts Options) *Migrataur {
	return &Migrataur{
		adapter: adapter,
		options: opts.Extend(DefaultOptions),
	}
}

func (m *Migrataur) getMigrationFullpath(name string) string {
	return filepath.Join(m.options.Directory,
		fmt.Sprintf("%s_%s%s", m.options.SequenceGenerator(), name, m.options.Extension))
}

// Init writes the initial migration provided by the adapter to create the needed
// migrations table, you should call it at the start of your project.
func (m *Migrataur) Init() (*Migration, error) {

	fullPath := m.getMigrationFullpath(m.options.InitialMigrationName)
	initialMigration := m.adapter.GetInitialMigration(filepath.Base(fullPath))

	if err := initialMigration.WriteTo(fullPath, m.options.MarshalOptions); err != nil {
		return nil, err
	}

	m.options.Logger.Printf("Migrataur initialized with %s", initialMigration.Name)

	return initialMigration, nil
}

// NewMigration creates a new migration in the configured folder and returns the instance of the migration
// attached to the newly created file
func (m *Migrataur) NewMigration(name string) *Migration {
	fullPath := m.getMigrationFullpath(name)
	migration := &Migration{Name: filepath.Base(fullPath)}

	if err := migration.WriteTo(fullPath, m.options.MarshalOptions); err != nil {
		m.options.Logger.Panic(err)
	}

	m.options.Logger.Printf("%s created", migration.Name)

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

		existingMigration := &Migration{Name: f.Name()}

		data, err := ioutil.ReadFile(filepath.Join(m.options.Directory, f.Name()))

		if err != nil {
			m.options.Logger.Panic(err)
		}

		if err = existingMigration.Unmarshal(data, m.options.MarshalOptions); err != nil {
			m.options.Logger.Panic(err)
		}

		migrations = append(migrations, existingMigration)
	}

	return migrations
}

func sortMigrations(migrations []*Migration, direction dir) {
	if direction == dirUp {
		sort.Sort(ByName(migrations))
	} else {
		sort.Sort(sort.Reverse(ByName(migrations)))
	}
}

func (m *Migrataur) getAllMigrations(direction dir) []*Migration {

	fileSystemMigrations := m.getAllFromFilesystem()
	adapterMigrations, err := m.adapter.GetAll()

	if err != nil {
		m.options.Logger.Panic(err)
	}

	// Constructs the migrations map to easily update them with adapter ones
	migrationsMap := map[string]*Migration{}

	for _, m := range fileSystemMigrations {
		migrationsMap[m.Name] = m
	}

	for _, mig := range adapterMigrations {
		fsMigration, ok := migrationsMap[mig.Name]

		if !ok {
			m.options.Logger.Panicf("The migration %s was not found in the migrations directory!", mig.Name)
		}

		fsMigration.hasBeenAppliedAt(*mig.AppliedAt)
	}

	sortMigrations(fileSystemMigrations, direction)

	return fileSystemMigrations
}

func getMigrationRange(rangeStr string) (first, last string) {
	if rangeStr == "" {
		return "", ""
	}

	splitted := strings.Split(rangeStr, "..")

	if len(splitted) == 1 {
		return splitted[0], ""
	}

	return splitted[0], splitted[1]
}

func (m *Migrataur) runFrom(start, end string, direction dir) ([]*Migration, error) {
	appliedMigrations := []*Migration{}

	if start == "" {
		return appliedMigrations, nil
	}

	startApplied := false

	for _, migration := range m.getAllMigrations(direction) {
		if !startApplied {
			if strings.Contains(migration.Name, start) {

				if m.runStep(migration, direction) {
					appliedMigrations = append(appliedMigrations, migration)
				}

				startApplied = true

				// Break early if no end migration has been set or if the end is the same
				if end == "" || start == end {
					break
				}
			}
		} else {
			if m.runStep(migration, direction) {
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

func (m *Migrataur) run(rangeOrName string, direction dir) []*Migration {
	start, end := getMigrationRange(rangeOrName)

	appliedMigrations, _ := m.runFrom(start, end, direction)

	return appliedMigrations
}

func (m *Migrataur) runStep(migration *Migration, direction dir) bool {

	// Do not execute commands if already applied or not applied at all when rolling back
	if migration.HasBeenApplied() && direction == dirUp {
		return false
	} else if !migration.HasBeenApplied() && direction == dirDown {
		return false
	}

	command := migration.Up

	if direction == dirDown {
		command = migration.Down
	}

	if err := m.adapter.Exec(command); err != nil {
		m.options.Logger.Panicf("✗\t%s: %s", migration.Name, err)
	}

	if direction == dirUp {
		now := time.Now().UTC()

		if err := m.adapter.AddMigration(migration.Name, now); err != nil {
			m.options.Logger.Panicf("✗\t%s: %s", migration.Name, err)
		}

		migration.hasBeenAppliedAt(now)
	} else {
		if err := m.adapter.RemoveMigration(migration.Name); err != nil {
			m.options.Logger.Panicf("✗\t%s: %s", migration.Name, err)
		}

		migration.hasBeenRolledBack()
	}

	m.options.Logger.Printf("✓\t%s", migration.Name)

	return true
}

// GetAll retrieve all migrations for the current instance. It will list applied and pending migrations
func (m *Migrataur) GetAll() []*Migration {
	m.options.Logger.Print("Fetching migrations")

	return m.getAllMigrations(dirUp)
}

// Migrate migrates the database and returns an array of effectively applied migrations (it will
// not contains those that were already applied.
// rangeOrName can be the exact migration name or a range such as <migration>..<another migration name>
func (m *Migrataur) Migrate(rangeOrName string) []*Migration {
	m.options.Logger.Printf("Applying %s", rangeOrName)

	return m.run(rangeOrName, dirUp)
}

// MigrateToLatest migrates the database to the latest version
func (m *Migrataur) MigrateToLatest() []*Migration {
	m.options.Logger.Print("Applying all pending migrations")

	migrations := m.getAllMigrations(dirUp)

	if len(migrations) == 0 {
		return []*Migration{}
	}

	appliedMigrations, _ := m.runFrom(migrations[0].Name, migrations[len(migrations)-1].Name, dirUp)

	return appliedMigrations
}

// Rollback inverts migrations and return an array of effectively rollbacked migrations
// (it will not contains those that were not applied).
func (m *Migrataur) Rollback(rangeOrName string) []*Migration {
	m.options.Logger.Printf("Rollbacking %s", rangeOrName)

	return m.run(rangeOrName, dirDown)
}

// Reset resets the database to its initial state
func (m *Migrataur) Reset() []*Migration {
	m.options.Logger.Print("Resetting database")

	migrations := m.getAllMigrations(dirDown)

	if len(migrations) == 0 {
		return []*Migration{}
	}

	appliedMigrations, _ := m.runFrom(migrations[0].Name, migrations[len(migrations)-1].Name, dirDown)

	return appliedMigrations
}
