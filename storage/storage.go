package storage

import (
	"fmt"
	"io"
	"mime/multipart"

	"codeinstyle.io/captain/config"
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

type Storage struct {
	name string
	Provider
}

func NewStorage(cfg *config.Config) *Storage {
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
		panic(fmt.Sprintf("Failed to initialize storage provider: %v", err))
	}

	return &Storage{
		name:     name,
		Provider: provider,
	}
}
