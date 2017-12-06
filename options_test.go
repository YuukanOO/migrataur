package migrataur

import "testing"
import "path/filepath"

func TestMarshalOptions(t *testing.T) {
	assert := assert(t)
	empty := MarshalOptions{}

	assert.equals(empty, emptyMarshalOptions)

	anotherOne := MarshalOptions{UpStart: "something"}

	assert.notEquals(emptyMarshalOptions, anotherOne)
}

func TestExtendEmptyOptions(t *testing.T) {
	opts := Options{}
	extended := opts.ExtendWith(DefaultOptions)

	fullpath, _ := filepath.Abs(DefaultOptions.Directory)

	assert(t).
		equals(fullpath, extended.Directory).
		equals(nil, extended.Logger).
		equals(DefaultOptions.Extension, extended.Extension).
		equals(DefaultOptions.SequenceGenerator(), extended.SequenceGenerator()).
		equals(DefaultOptions.MarshalOptions, extended.MarshalOptions).
		equals(DefaultOptions.InitialMigrationName, extended.InitialMigrationName)
}

func TestExtendOptions(t *testing.T) {

	marshalOpts := MarshalOptions{
		UpStart:   "-- up",
		UpEnd:     "-- /up",
		DownStart: "-- down",
		DownEnd:   "-- /down",
	}

	opts := Options{
		Directory:      "./MigrationsGoesHere",
		Extension:      "myext",
		MarshalOptions: marshalOpts,
	}

	extended := opts.ExtendWith(DefaultOptions)

	fullpath, _ := filepath.Abs("./MigrationsGoesHere")

	assert(t).
		equals(fullpath, extended.Directory).
		equals(".myext", extended.Extension).
		equals(marshalOpts, extended.MarshalOptions)
}
