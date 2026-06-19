package services

import "net/http"

type Storage interface {
	Upload(data []byte, filePath string) error
	Download(path string) ([]byte, string, error)
	DownloadStream(path string, headers http.Header) (*http.Response, error)
	Delete(filePath string) error
}
