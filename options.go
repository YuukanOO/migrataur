package migrataur

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
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
}

func extendOptionsAndSanitize(opts *Options) *Options {

	dir, extension, generator, logger := opts.Directory, opts.Extension, opts.SequenceGenerator, opts.Logger

	if dir == "" {
		dir = "./migrations"
	}

	absPath, err := filepath.Abs(dir)

	if err != nil {
		panic(fmt.Sprintf("Could not retrieve the absolute path for %s", dir))
	}

	if err = os.MkdirAll(absPath, os.ModeDir); err != nil {
		panic(err)
	}

	if extension == "" {
		extension = ".sql"
	}

	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}

	if generator == nil {
		generator = func() string { return strconv.FormatInt(time.Now().Unix(), 10) }
	}

	if logger == nil {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	}

	return &Options{
		Logger:            logger,
		Directory:         absPath,
		Extension:         extension,
		SequenceGenerator: generator,
	}
}
