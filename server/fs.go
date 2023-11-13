package server

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func CreateFsDir() {
	path := os.Getenv("FS")
	err := mkdir(path)
	if err != nil {
		log.Printf("Failed to create fs directory at %s: %v", path, err)
	}
}

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
	// Stat the directory to see if it exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create the dir if it doesn't exist
		if err = mkdir(path); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	err := os.WriteFile(path, data, 0777)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

// Creates a directory at the specified path and returns any errors that occured.
func mkdir(path string) error {
	path = filepath.Dir(path)
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	return nil
}
