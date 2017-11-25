package migrataur

import "time"

// Migration represents a database migration :)
type Migration struct {
	name      string
	appliedAt *time.Time
}
