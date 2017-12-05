package migrataur

// MarshalOptions holds configuration for the migration marshaling & unmarshaling
type MarshalOptions struct {
	UpStart   string
	UpEnd     string
	DownStart string
	DownEnd   string
}

// DefaultMarshalOptions holds default marshal options for the migration used when
// writing or reading migration files to the filesystem.
var DefaultMarshalOptions = MarshalOptions{
	UpStart:   "-- +migrataur up",
	UpEnd:     "-- -migrataur up",
	DownStart: "-- +migrataur down",
	DownEnd:   "-- -migrataur down",
}

var emptyMarshalOptions = MarshalOptions{}
