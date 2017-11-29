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
func (m *Migrataur) Init() *Migration {

	initialMigration := m.adapter.GetInitialMigration()

	// Updates the name based on the generated path
	fullPath := m.getMigrationFullpath(initialMigration.Name)
	initialMigration.Name = filepath.Base(fullPath)

	if err := initialMigration.WriteTo(fullPath, m.options.MarshalOptions); err != nil {
		m.options.Logger.Panic(err)
	}

	m.options.Logger.Printf("Migrataur initialized with %s", initialMigration.Name)

	return initialMigration
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
		migrationsMap[m.Name] = m
	}

	for _, mig := range adapterMigrations {
		fsMigration, ok := migrationsMap[mig.Name]

		if !ok {
			m.options.Logger.Panicf("The migration %s was not found in the migrations directory!", mig.Name)
		}

		fsMigration.hasBeenAppliedAt(*mig.AppliedAt)
	}

	sortMigrations(fileSystemMigrations)

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

// Migrate migrates the database.
// rangeOrName can be the exact migration name or a range such as <migration>..<another migration name>
func (m *Migrataur) Migrate(rangeOrName string) {
	m.options.Logger.Printf("Applying %s", rangeOrName)

	start, end := getMigrationRange(rangeOrName)

	startApplied := false

	for _, migration := range m.GetAll() {
		if !startApplied {
			if strings.Contains(migration.Name, start) {
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
			if strings.Contains(migration.Name, end) {
				break
			}
		}
	}
}

// Rollback inverts migration
func (m *Migrataur) Rollback(rangeOrName string) {
	m.options.Logger.Printf("Rollbacking %s", rangeOrName)

	end, start := getMigrationRange(rangeOrName)

	endApplied := false
	migrations := m.GetAll()

	sort.Sort(sort.Reverse(ByName(migrations)))

	for _, migration := range migrations {
		if !endApplied {
			if strings.Contains(migration.Name, end) {
				m.rollbackMigration(migration)

				endApplied = true

				if start == "" || start == end {
					break
				}
			}
		} else {
			m.rollbackMigration(migration)

			if strings.Contains(migration.Name, start) {
				break
			}
		}
	}
}

// TODO: We should merge apply and rollback into one function
// and write tests for Rollback
func (m *Migrataur) applyMigration(migration *Migration) bool {
	if migration.HasBeenApplied() {
		return false
	}

	if err := m.adapter.Exec(migration.Up); err != nil {
		m.options.Logger.Panicf("✗\t%s: %s", migration.Name, err)
	}

	now := time.Now().UTC()

	if err := m.adapter.AddMigration(migration.Name, now); err != nil {
		m.options.Logger.Panicf("✗\t%s: %s", migration.Name, err)
	}

	migration.hasBeenAppliedAt(now)

	m.options.Logger.Printf("✓\t%s", migration.Name)

	return true
}

func (m *Migrataur) rollbackMigration(migration *Migration) bool {
	if !migration.HasBeenApplied() {
		return false
	}

	if err := m.adapter.Exec(migration.Down); err != nil {
		m.options.Logger.Panicf("✗\t%s: %s", migration.Name, err)
	}

	if err := m.adapter.RemoveMigration(migration.Name); err != nil {
		m.options.Logger.Panicf("✗\t%s: %s", migration.Name, err)
	}

	migration.hasBeenRolledBack()

	m.options.Logger.Printf("✓\t%s", migration.Name)

	return true
}

// MigrateToLatest migrates the database to the latest version
func (m *Migrataur) MigrateToLatest() {
	m.options.Logger.Print("Applying all pending migrations")

	for _, migration := range m.GetAll() {
		m.applyMigration(migration)
	}
}

// Reset resets the database to its initial state
func (m *Migrataur) Reset() {
	m.options.Logger.Print("Resetting database")

	migrations := m.GetAll()

	sort.Sort(sort.Reverse(ByName(migrations)))

	for _, migration := range migrations {
		m.rollbackMigration(migration)
	}
}
