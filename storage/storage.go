package storage

import (
	"io"
	"mime/multipart"
)

// Provider defines the interface for storage providers
type Provider interface {
	// Save stores a file and returns its path and any error
	Save(file *multipart.FileHeader) (string, error)

	// Delete removes a file and returns any error
	Delete(path string) error

	// Get retrieves a file and returns a ReadCloser and any error
	Get(path string) (io.ReadCloser, error)
}
