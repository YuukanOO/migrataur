package migrataur

import (
	"fmt"
	"strings"
	"time"
)

const (
	upStart   = "-- +migrataur up"
	upEnd     = "-- -migrataur up"
	downStart = "-- +migrataur down"
	downEnd   = "-- -migrataur down"
)

// Migration represents a database migration :)
type Migration struct {
	name      string
	upStr     string
	downStr   string
	appliedAt *time.Time
}

func newMigration(name string) *Migration {
	return &Migration{
		name: name,
	}
}

// ByName sort an array of migrations by their name, use it with sort.Sort and the like
type ByName []*Migration

func (m ByName) Len() int           { return len(m) }
func (m ByName) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m ByName) Less(i, j int) bool { return m[i].name < m[j].name }

// NewMigration instantiates a new migration. It should be used exclusively
// by adapters.
func NewMigration(name string, appliedAt time.Time) *Migration {
	return &Migration{
		name:      name,
		appliedAt: &appliedAt,
	}
}

func (m *Migration) String() string {
	ticked := " "

	if m.appliedAt != nil {
		ticked = "✓"
	}

	return fmt.Sprintf("[%s]\t%s", ticked, m.name)
}

func (m *Migration) hasBeenAppliedAt(time time.Time) {
	m.appliedAt = &time
}

// HasBeenApplied checks if the migration has already been applied in the database
func (m *Migration) HasBeenApplied() bool {
	return m.appliedAt != nil
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
`, m.name, upStart, m.upStr, upEnd, downStart, m.downStr, downEnd)

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
			m.upStr = strings.Join(lines[upFrom+1:i], "\n")
		case downStart:
			downFrom = i
		case downEnd:
			m.downStr = strings.Join(lines[downFrom+1:i], "\n")
		}
	}

	return nil
}
