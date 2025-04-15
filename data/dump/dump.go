package dump

import (
	"github.com/inovacc/gopickle"
)

// Dump serializes the provided data and writes it to the given filename.
// It uses GoPickle format for storing structured data persistently.
// Returns an error if serialization or file write fails.
func Dump(filename string, data any) error {
	return gopickle.Dump(filename, data)
}

// Load deserializes data from the provided filename into the passed pointer.
// Expects the data to have been previously stored with Dump using the same structure.
// Returns an error if file read or unmarshalling fails.
func Load(filename string, data any) error {
	return gopickle.Load(filename, data)
}
