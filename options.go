package migrataur

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Options represents migrataur options to give to an instance
type Options struct {
	Directory         string
	Extension         string
	SequenceGenerator func() string
}

func extendOptionsAndSanitize(opts *Options) *Options {

	dir, extension, generator := opts.Directory, opts.Extension, opts.SequenceGenerator

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

	return &Options{
		Directory:         absPath,
		Extension:         extension,
		SequenceGenerator: generator,
	}
}
