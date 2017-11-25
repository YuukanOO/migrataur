package migrataur

import (
	"fmt"
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

// MarshalText serialize this migration
func (m *Migration) MarshalText() (text []byte, err error) {
	content := fmt.Sprintf(`-- Migrations %s
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
	return nil
}
