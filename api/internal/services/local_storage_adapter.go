package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type LocalStorageAdapter struct {
	dir string
}

func NewLocalStorageAdapter(dir string) *LocalStorageAdapter {
	_ = os.MkdirAll(dir, 0755)
	return &LocalStorageAdapter{dir: dir}
}

func (a *LocalStorageAdapter) Upload(data []byte, filePath string) error {
	dest := filepath.Join(a.dir, filePath)
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}
	return os.WriteFile(dest, data, 0644)
}

func (a *LocalStorageAdapter) Download(path string) ([]byte, string, error) {
	resp, err := a.DownloadStream(path, nil)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read file: %w", err)
	}

	return data, resp.Header.Get("Content-Type"), nil
}

func (a *LocalStorageAdapter) DownloadStream(path string, _ http.Header) (*http.Response, error) {
	fullPath := filepath.Join(a.dir, path)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read local file: %w", err)
	}

	mimeType := http.DetectContentType(data)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewReader(data)),
	}
	resp.Header.Set("Content-Type", mimeType)
	resp.Header.Set("Content-Length", strconv.Itoa(len(data)))

	return resp, nil
}

func (a *LocalStorageAdapter) Delete(filePath string) error {
	return os.Remove(filepath.Join(a.dir, filePath))
}
