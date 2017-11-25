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

// NewAdapterMigration instantiates a new migration. It should be used exclusively
// by adapters
func NewAdapterMigration(name string, appliedAt time.Time) *Migration {
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

// MarshalText serialize this migration
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

// UnmarshalText deserialize a migration
func (m *Migration) UnmarshalText(text []byte) error {
	lines := strings.Split(string(text), "\n")

	upFrom, downFrom := 0, 0

	for i := 0; i < len(lines); i++ {
		switch lines[i] {
		case upStart:
			upFrom = i
			break
		case upEnd:
			m.upStr = strings.Join(lines[upFrom+1:i], "\n")
			break
		case downStart:
			downFrom = i
			break
		case downEnd:
			m.downStr = strings.Join(lines[downFrom+1:i], "\n")
			break
		}
	}

	return nil
}
