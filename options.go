package migrataur

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Options represents migrataur options to give to an instance
type Options struct {
	Directory        string
	Extension        string
	UnicityGenerator func() string
}

func extendOptions(opts *Options) *Options {

	dir, extension, generator := opts.Directory, opts.Extension, opts.UnicityGenerator

	if dir == "" {
		dir = "./migrations"
	}

	absPath, err := filepath.Abs(dir)

	if err != nil {
		panic(fmt.Sprintf("Could not retrieve the absolute path for %s", dir))
	}

	if extension == "" {
		extension = ".sql"
	}

	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}

	if generator == nil {
		generator = currentTimestamp
	}

	return &Options{
		Directory:        absPath,
		Extension:        extension,
		UnicityGenerator: generator,
	}
}
