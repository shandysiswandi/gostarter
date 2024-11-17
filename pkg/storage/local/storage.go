// Package local provides a local filesystem implementation of the Storage interface.
package local

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
)

// OwnerReadWrite is expression is a Unix file permission mode
// that determines the read and write permissions for the file.
const OwnerReadWrite fs.FileMode = 0o644

// Storage implements the Storage interface using the local filesystem.
type Storage struct {
	basePath string
}

// NewStorage creates a new instance of Storage using the provided base path.
func NewStorage(basePath string) *Storage {
	return &Storage{basePath: basePath}
}

// Close does nothing in the local storage implementation, as there are no resources to close.
func (s *Storage) Close() error {
	return nil
}

// Upload saves the provided data to a file specified by the key within the basePath.
func (s *Storage) Upload(_ context.Context, key string, data []byte) (string, error) {
	path := filepath.Join(s.basePath, key)

	if err := os.WriteFile(path, data, OwnerReadWrite); err != nil {
		return "", err
	}

	return path, nil
}

// Download retrieves the data from a file specified by the key within the basePath.
func (s *Storage) Download(_ context.Context, key string) ([]byte, error) {
	return os.ReadFile(filepath.Join(s.basePath, key))
}

// Delete removes the file specified by the key within the basePath.
func (s *Storage) Delete(_ context.Context, key string) error {
	return os.Remove(filepath.Join(s.basePath, key))
}

// List returns a list of file paths within the basePath that match the specified prefix.
func (s *Storage) List(_ context.Context, prefix string) ([]string, error) {
	var files []string
	path := filepath.Join(s.basePath, prefix)

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			files = append(files, path)
		}

		return err
	})

	return files, err
}
