// Package storage defines a common interface for file storage systems.
// Implementations may include local file storage, Amazon S3, Google Cloud Storage, and others.
package storage

import (
	"context"
	"io"
)

// Storage represents a file storage system.
// It providing a comprehensive API for managing files in various storage systems.
type Storage interface {
	// Close releases any resources held by the storage system.
	// It is intended to be called when the storage system is no longer needed.
	// Implementations that do not hold resources can return nil.
	io.Closer

	// Upload saves the provided data to a storage location specified by the key.
	// The key is a unique identifier within the storage system (e.g., file path, S3 object key).
	// It returns the key or an error if the operation fails.
	Upload(ctx context.Context, key string, data []byte) (string, error)

	// Download retrieves the data from the storage location specified by the key.
	// The key is a unique identifier within the storage system.
	// It returns the data or an error if the operation fails.
	Download(ctx context.Context, key string) ([]byte, error)

	// Delete removes the data from the storage location specified by the key.
	// The key is a unique identifier within the storage system.
	// It returns an error if the operation fails.
	Delete(ctx context.Context, key string) error

	// List returns a list of keys that match the specified prefix within the storage system.
	// The prefix is used to filter the keys (e.g., directory path, S3 prefix).
	// It returns a slice of keys or an error if the operation fails.
	List(ctx context.Context, prefix string) ([]string, error)
}
