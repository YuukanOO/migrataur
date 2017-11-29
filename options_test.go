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
	extended := opts.Extend(DefaultOptions)

	fullpath, _ := filepath.Abs(DefaultOptions.Directory)

	if extended.Directory != fullpath {
		shouldHaveBeenEquals(t, fullpath, extended.Directory)
	}

	if extended.Logger != DefaultOptions.Logger {
		t.Fail()
	}

	if extended.Extension != DefaultOptions.Extension {
		shouldHaveBeenEquals(t, DefaultOptions.Extension, extended.Extension)
	}

	if extended.SequenceGenerator() != DefaultOptions.SequenceGenerator() {
		shouldHaveBeenEquals(t, DefaultOptions.SequenceGenerator, extended.SequenceGenerator)
	}

	if extended.MarshalOptions != DefaultOptions.MarshalOptions {
		shouldHaveBeenEquals(t, DefaultOptions.MarshalOptions, extended.MarshalOptions)
	}
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

	extended := opts.Extend(DefaultOptions)

	fullpath, _ := filepath.Abs("./MigrationsGoesHere")

	if extended.Directory != fullpath {
		shouldHaveBeenEquals(t, fullpath, extended.Directory)
	}

	if extended.Extension != ".myext" {
		shouldHaveBeenEquals(t, ".myext", extended.Extension)
	}

	if extended.MarshalOptions != marshalOpts {
		shouldHaveBeenEquals(t, marshalOpts, extended.MarshalOptions)
	}
}
