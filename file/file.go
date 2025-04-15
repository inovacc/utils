package file

import (
	"bytes"
	"os"
)

// WriteToFile writes the given byte data to a specified file.
// It creates or truncates the file, writes the data, flushes it to disk, and closes the file.
// Returns an error if any step fails.
func WriteToFile(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err // Failed to create or truncate the file
	}

	if _, err := file.Write(data); err != nil {
		return err // Failed to write data
	}

	if err := file.Sync(); err != nil {
		return err // Failed to flush to disk
	}

	if err := file.Close(); err != nil {
		return err // Failed to close file
	}
	return nil
}

// ReadFromFile reads the contents of a specified file into memory and returns it as a byte slice.
// Returns an error if the file cannot be opened, read, or closed.
func ReadFromFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err // Failed to open file
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return nil, err // Failed to read content into buffer
	}

	if err := file.Close(); err != nil {
		return nil, err // Failed to close file
	}
	return buf.Bytes(), nil
}
