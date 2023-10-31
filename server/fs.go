package server

import (
	"os"
)

// Reads and returns the file contents of the specified path and any errors that occured.
func GetFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Writes a file to the specified path and returns any errors that occured.
func WriteFile(path string, data []byte) error {
	err := os.WriteFile(path, data, 0777) // FIXME: bad perms
	if err != nil {
		return err
	}

	return nil
}
