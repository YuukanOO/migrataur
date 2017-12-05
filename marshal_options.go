package migrataur

// MarshalOptions holds configuration for the migration marshaling & unmarshaling
type MarshalOptions struct {
	Header    string
	UpStart   string
	UpEnd     string
	DownStart string
	DownEnd   string
}

// DefaultMarshalOptions holds default marshal options for the migration used when
// writing or reading migration files to the filesystem.
var DefaultMarshalOptions = MarshalOptions{
	Header:    "-- migration",
	UpStart:   "-- +migrataur up",
	UpEnd:     "-- -migrataur up",
	DownStart: "-- +migrataur down",
	DownEnd:   "-- -migrataur down",
}

var emptyMarshalOptions = MarshalOptions{}
