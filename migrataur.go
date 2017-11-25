package migrataur

import (
	"fmt"
	"os"
	"path"
)

const (
	upStart   = "-- +migrataur up"
	upEnd     = "-- -migrataur up"
	downStart = "-- +migrataur down"
	downEnd   = "-- -migrataur down"
)

// Migrataur represents an instance configurated for a particular use
type Migrataur struct {
	options *Options
}

// New instantiates a new Migrataur instance for the given options
func New(opts *Options) *Migrataur {
	return &Migrataur{options: extendOptions(opts)}
}

// NewMigration creates a new migration in the configured folder and returns the instance of the migration
// attached to the newly created file
func (m *Migrataur) NewMigration(name string) *Migration {

	fullPath := path.Join(m.options.Directory,
		fmt.Sprintf("%s_%s%s", m.options.UnicityGenerator(), name, m.options.Extension))

	content := fmt.Sprintf(`-- Migrations %s
%s

%s


%s

%s
`, path.Base(fullPath), upStart, upEnd, downStart, downEnd)

	file, err := os.Create(fullPath)

	if err != nil {
		panic(err)
	}

	_, err = file.WriteString(content)

	if err != nil {
		panic(err)
	}

	return nil
}
