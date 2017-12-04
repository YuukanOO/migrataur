package migrataur

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger is the interface to be implemented by a logger since golang does not
// exposes a common interface. This interface is taken from the logrus library.
type Logger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}

// Options represents migrataur options to give to an instance
type Options struct {
	Logger               Logger
	Directory            string
	Extension            string
	InitialMigrationName string
	SequenceGenerator    func() string
	MarshalOptions       MarshalOptions
}

// DefaultOptions represents the default migrataur options
var DefaultOptions = Options{
	Logger:               log.New(os.Stdout, "", log.LstdFlags),
	Directory:            "./migrations",
	Extension:            ".sql",
	InitialMigrationName: "initMigrataur",
	SequenceGenerator:    GetCurrentTimeFormatted,
	MarshalOptions:       DefaultMarshalOptions,
}

// ExtendWith self options with the given one. It means that if a field is not
// present in this option, it will be replaced by the one in the other Options.
func (opts Options) ExtendWith(other Options) Options {
	result := opts

	if result.Logger == nil {
		result.Logger = other.Logger
	}

	if result.Directory == "" {
		result.Directory = other.Directory
	}

	// Sanitize path
	absPath, err := filepath.Abs(result.Directory)

	if err != nil {
		panic(err)
	}

	result.Directory = absPath

	if result.Extension == "" {
		result.Extension = other.Extension
	}

	if result.InitialMigrationName == "" {
		result.InitialMigrationName = other.InitialMigrationName
	}

	// Sanitize extension
	if result.Extension[0] != '.' {
		result.Extension = "." + result.Extension
	}

	if result.SequenceGenerator == nil {
		result.SequenceGenerator = other.SequenceGenerator
	}

	if result.MarshalOptions == emptyMarshalOptions {
		result.MarshalOptions = other.MarshalOptions
	}

	return result
}

// GetCurrentTimeFormatted retrieves the current time formatted. It's used as
// the default sequence generator since a linux timestamp goes up to 2038 :)
func GetCurrentTimeFormatted() string {
	return time.Now().UTC().Format("20060102150405")
}
