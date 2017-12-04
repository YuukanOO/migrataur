package migrataur

import "testing"
import "path/filepath"

func TestMarshalOptions(t *testing.T) {
	empty := MarshalOptions{}

	if empty != emptyMarshalOptions {
		t.Fail()
	}

	anotherOne := MarshalOptions{UpStart: "something"}

	if anotherOne == emptyMarshalOptions {
		t.Fail()
	}
}

func TestExtendEmptyOptions(t *testing.T) {
	opts := Options{}
	extended := opts.ExtendWith(DefaultOptions)

	fullpath, _ := filepath.Abs(DefaultOptions.Directory)

	assertEquals(t, fullpath, extended.Directory)

	if extended.Logger != DefaultOptions.Logger {
		t.Fail()
	}

	assertEquals(t, DefaultOptions.Extension, extended.Extension)
	assertEquals(t, DefaultOptions.SequenceGenerator(), extended.SequenceGenerator())
	assertEquals(t, DefaultOptions.MarshalOptions, extended.MarshalOptions)
	assertEquals(t, DefaultOptions.InitialMigrationName, extended.InitialMigrationName)
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

	assertEquals(t, fullpath, extended.Directory)
	assertEquals(t, ".myext", extended.Extension)
	assertEquals(t, marshalOpts, extended.MarshalOptions)
}
