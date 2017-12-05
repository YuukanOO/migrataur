package migrataur

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Migration represents a database migration, nothing more.
type Migration struct {
	Name      string
	Up        string
	Down      string
	AppliedAt *time.Time
}

// byName sort an array of migrations by their name, use it with sort.Sort and the like
type byName []*Migration

func (m byName) Len() int           { return len(m) }
func (m byName) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m byName) Less(i, j int) bool { return m[i].Name < m[j].Name }

func (m *Migration) String() string {
	ticked := " "

	if m.HasBeenApplied() {
		ticked = "✓"
	}

	return fmt.Sprintf("[%s]\t%s", ticked, m.Name)
}

func (m *Migration) hasBeenAppliedAt(time time.Time) {
	m.AppliedAt = &time
}

func (m *Migration) hasBeenRolledBack() {
	m.AppliedAt = nil
}

// HasBeenApplied checks if the migration has already been applied in the database
func (m *Migration) HasBeenApplied() bool {
	return m.AppliedAt != nil
}

// Marshal serializes this migration
func (m *Migration) Marshal(options MarshalOptions) (text []byte, err error) {
	content := fmt.Sprintf(`%s %s
%s
%s
%s


%s
%s
%s
`, options.Header, m.Name, options.UpStart, m.Up, options.UpEnd, options.DownStart, m.Down, options.DownEnd)

	return []byte(content), nil
}

// Unmarshal deserializes a migration
func (m *Migration) Unmarshal(text []byte, options MarshalOptions) error {
	lines := strings.Split(string(text), "\n")

	upFrom, downFrom := 0, 0

	for i := 0; i < len(lines); i++ {
		switch lines[i] {
		case options.UpStart:
			upFrom = i
		case options.UpEnd:
			m.Up = strings.Join(lines[upFrom+1:i], "\n")
		case options.DownStart:
			downFrom = i
		case options.DownEnd:
			m.Down = strings.Join(lines[downFrom+1:i], "\n")
		}
	}

	return nil
}

// WriteTo writes this migration to the filesystem
func (m *Migration) WriteTo(path string, options MarshalOptions) error {

	// Make sure the directory exists
	if err := os.MkdirAll(filepath.Dir(path), os.ModeDir); err != nil {
		return err
	}

	file, err := os.Create(path)

	if err != nil {
		return err
	}

	defer file.Close()

	data, err := m.Marshal(options)

	if err != nil {
		return err
	}

	_, err = file.Write(data)

	return err
}
