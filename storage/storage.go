package storage

import (
	"fmt"
	"io"

	"github.com/captain-corp/captain/config"
)

// Provider defines the interface for storage providers
type Provider interface {
	// Save stores a file from a reader and returns its path and any error
	Save(filename string, reader io.Reader) (string, error)

	// Delete removes a file and returns any error
	Delete(path string) error

	// Get retrieves a file and returns a ReadCloser and any error
	Get(path string) (io.ReadCloser, error)
}

// Storage wraps a Provider with its name
type Storage struct {
	name string
	Provider
}

func NewStorage(cfg *config.Config) (*Storage, error) {
	var provider Provider
	var err error
	name := cfg.Storage.Provider

	switch name {
	case "s3":
		provider, err = NewS3Provider(cfg.Storage.S3.Bucket, cfg.Storage.S3.Region, cfg.Storage.S3.Endpoint, cfg.Storage.S3.AccessKey, cfg.Storage.S3.SecretKey)
	default:
		provider, err = NewLocalProvider(cfg.Storage.LocalPath)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage provider: %w", err)
	}

	return &Storage{
		name:     name,
		Provider: provider,
	}, nil
}
