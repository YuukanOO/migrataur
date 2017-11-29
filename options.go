package migrataur

import (
	"log"
	"os"
	"path/filepath"
	"strings"
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
	Logger            Logger
	Directory         string
	Extension         string
	SequenceGenerator func() string
	MarshalOptions    *MarshalOptions
}

// GetCurrentTimeFormatted retrieves the current time formatted. It's used as
// the default sequence generator since a linux timestamp goes up to 2038 :)
func GetCurrentTimeFormatted() string {
	return time.Now().UTC().Format("20060102150405")
}

func extendOptionsAndSanitize(opts *Options) *Options {

	dir, extension, generator, logger, marshalOpts :=
		opts.Directory,
		opts.Extension,
		opts.SequenceGenerator,
		opts.Logger, opts.MarshalOptions

	if logger == nil {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	}

	if dir == "" {
		dir = "./migrations"
	}

	absPath, err := filepath.Abs(dir)

	if err != nil {
		logger.Panicf("Could not retrieve the absolute path for %s", dir)
	}

	if err = os.MkdirAll(absPath, os.ModeDir); err != nil {
		logger.Panic(err)
	}

	if extension == "" {
		extension = ".sql"
	}

	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}

	if generator == nil {
		generator = GetCurrentTimeFormatted
	}

	if marshalOpts == nil {
		marshalOpts = &DefaultMarshalOptions
	}

	return &Options{
		Logger:            logger,
		Directory:         absPath,
		Extension:         extension,
		SequenceGenerator: generator,
		MarshalOptions:    marshalOpts,
	}
}
