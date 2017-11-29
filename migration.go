package migrataur

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	upStart   = "-- +migrataur up"
	upEnd     = "-- -migrataur up"
	downStart = "-- +migrataur down"
	downEnd   = "-- -migrataur down"
)

// MarshalOptions holds configuration for the migration marshaling & unmarshaling
type MarshalOptions struct {
	UpStart   string
	UpEnd     string
	DownStart string
	DownEnd   string
}

// DefaultMarshalOptions holds default marshal options for the migration
var DefaultMarshalOptions = MarshalOptions{
	UpStart:   "-- +migrataur up",
	UpEnd:     "-- -migrataur up",
	DownStart: "-- +migrataur down",
	DownEnd:   "-- -migrataur down",
}

// Migration represents a database migration :)
type Migration struct {
	Name      string
	Up        string
	Down      string
	AppliedAt *time.Time
}

// ByName sort an array of migrations by their name, use it with sort.Sort and the like
type ByName []*Migration

func (m ByName) Len() int           { return len(m) }
func (m ByName) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m ByName) Less(i, j int) bool { return m[i].Name < m[j].Name }

func (m *Migration) String() string {
	ticked := " "

	if m.AppliedAt != nil {
		ticked = "✓"
	}

	return fmt.Sprintf("[%s]\t%s", ticked, m.Name)
}

func (m *Migration) hasBeenAppliedAt(time time.Time) {
	m.AppliedAt = &time
}

// HasBeenApplied checks if the migration has already been applied in the database
func (m *Migration) HasBeenApplied() bool {
	return m.AppliedAt != nil
}

// MarshalText serializes this migration
func (m *Migration) MarshalText() (text []byte, err error) {
	content := fmt.Sprintf(`-- migration %s
%s
%s
%s


%s
%s
%s
`, m.Name, upStart, m.Up, upEnd, downStart, m.Down, downEnd)

	return []byte(content), nil
}

// UnmarshalText deserializes a migration
func (m *Migration) UnmarshalText(text []byte) error {
	lines := strings.Split(string(text), "\n")

	upFrom, downFrom := 0, 0

	for i := 0; i < len(lines); i++ {
		switch lines[i] {
		case upStart:
			upFrom = i
		case upEnd:
			m.Up = strings.Join(lines[upFrom+1:i], "\n")
		case downStart:
			downFrom = i
		case downEnd:
			m.Down = strings.Join(lines[downFrom+1:i], "\n")
		}
	}

	return nil
}

// WriteTo writes this migration to the filesystem
func (m *Migration) WriteTo(path string, options MarshalOptions) error {
	file, err := os.Create(path)

	if err != nil {
		return err
	}

	defer file.Close()

	data, err := m.MarshalText()

	if err != nil {
		return err
	}

	_, err = file.Write(data)

	return err
}
