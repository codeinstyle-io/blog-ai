package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalProvider implements Provider interface for local filesystem storage
type LocalProvider struct {
	baseDir string
}

// NewLocalProvider creates a new LocalProvider
func NewLocalProvider(baseDir string) (*LocalProvider, error) {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %v", err)
	}
	return &LocalProvider{baseDir: baseDir}, nil
}

// Save implements Provider.Save
func (p *LocalProvider) Save(filename string, reader io.Reader) (string, error) {
	// Generate unique filename with slugified name
	// ext := filepath.Ext(filename)
	// name := filename[:len(filename)-len(ext)]
	// filename := fmt.Sprintf("%d-%s%s", time.Now().Unix(), slugify(name), ext)
	path := filepath.Join(p.baseDir, filename)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", fmt.Errorf("failed to create directories: %w", err)
	}

	dst, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, reader); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	return filename, nil
}

// Delete implements Provider.Delete
func (p *LocalProvider) Delete(path string) error {
	fullPath := filepath.Join(p.baseDir, path)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// Get implements Provider.Get
func (p *LocalProvider) Get(path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(p.baseDir, path)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	return file, nil
}
