package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"embox/internal/config"
)

type StorageService struct {
	token     string
	tokenLock sync.Mutex
	config    *config.StorageConfig
	client    *http.Client
}

func NewStorageService(cfg *config.StorageConfig) *StorageService {
	return &StorageService{
		config: cfg,
		client: &http.Client{},
	}
}

// Retrieves and caches an authentication token.
func (s *StorageService) Auth() error {
	s.tokenLock.Lock()
	defer s.tokenLock.Unlock()

	if s.token != "" {
		return nil // Token is already available
	}

	body := strings.NewReader(fmt.Sprintf("username=%s&password=%s", s.config.Username, s.config.Password))
	resp, err := http.Post(s.config.Url+"/auth-token/", "application/x-www-form-urlencoded", body)
	if err != nil {
		return fmt.Errorf("failed to request auth token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("auth failed: %s", string(b))
	}

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode auth response: %w", err)
	}

	s.token = result.Token
	return nil
}

// Uploads a file to the repository.
// path is a full path in the repo, e.g. "2025/09/15_123.jpg"
func (s *StorageService) Upload(data []byte, filePath string) error {
	if err := s.Auth(); err != nil {
		return err
	}

	parentDir := "/"
	relativePath := filepath.Dir(filePath) + "/"
	filename := filepath.Base(filePath)

	uploadURL, err := s.getUploadUrl(parentDir)
	if err != nil {
		return fmt.Errorf("failed to get upload URL: %w", err)
	}

	// Prepare the multipart body
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add the file
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := part.Write(data); err != nil {
		return fmt.Errorf("failed to write file data: %w", err)
	}

	// Add Params parent_dir, relative_path and replace
	if err := writer.WriteField("parent_dir", parentDir); err != nil {
		return fmt.Errorf("failed to write parent_dir field: %w", err)
	}
	if relativePath != "" {
		if err := writer.WriteField("relative_path", relativePath); err != nil {
			return fmt.Errorf("failed to write relative_path field: %w", err)
		}
	}
	if err := writer.WriteField("replace", "1"); err != nil {
		return fmt.Errorf("failed to write replace field: %w", err)
	}

	// Close the writer to finalize the multipart body
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Upload the file
	upReq, err := http.NewRequest("POST", uploadURL, &buf)
	if err != nil {
		return fmt.Errorf("failed to create upload request: %w", err)
	}
	upReq.Header.Set("Authorization", "Token "+s.token)
	upReq.Header.Set("Content-Type", writer.FormDataContentType())

	upResp, err := s.client.Do(upReq)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	defer upResp.Body.Close()

	if upResp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(upResp.Body)
		return fmt.Errorf("upload failed: %s", string(b))
	}

	return nil
}

// Download loads a file and returns the data + MIME type.
func (s *StorageService) Download(path string) ([]byte, string, error) {
	resp, err := s.DownloadStream(path, nil)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read stream data: %w", err)
	}

	return data, resp.Header.Get("Content-Type"), nil
}

// DownloadStream returns an http.Response for the file at path.
// The caller is responsible for closing the response body.
func (s *StorageService) DownloadStream(path string, headers http.Header) (*http.Response, error) {
	if err := s.Auth(); err != nil {
		return nil, err
	}

	// Get the file download link
	fileURL := fmt.Sprintf("%s/repos/%s/file/?p=/%s", s.config.Url, s.config.RepoID, path)
	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create download request: %w", err)
	}
	req.Header.Set("Authorization", "Token "+s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request file url: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("file url request failed: %s", string(b))
	}

	var realURL string
	if err := json.NewDecoder(resp.Body).Decode(&realURL); err != nil {
		return nil, fmt.Errorf("failed to decode file url: %w", err)
	}

	// Load the actual stream
	fileReq, err := http.NewRequest("GET", realURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create stream request: %w", err)
	}

	// Pass through headers (important for Range requests)
	for key, values := range headers {
		for _, value := range values {
			fileReq.Header.Add(key, value)
		}
	}

	fileResp, err := s.client.Do(fileReq)
	if err != nil {
		return nil, fmt.Errorf("failed to download stream: %w", err)
	}

	if fileResp.StatusCode != http.StatusOK && fileResp.StatusCode != http.StatusPartialContent {
		defer fileResp.Body.Close()
		b, _ := io.ReadAll(fileResp.Body)
		return nil, fmt.Errorf("stream download failed with status %d: %s", fileResp.StatusCode, string(b))
	}

	return fileResp, nil
}

// Delete entfernt eine Datei oder ein Verzeichnis aus dem Repo.
// filePath ist ein vollständiger Pfad im Repo, z.B. "2025/09/15_123.jpg"
func (s *StorageService) Delete(filePath string) error {
	if err := s.Auth(); err != nil {
		return err
	}

	deleteURL := fmt.Sprintf("%s/repos/%s/file/?p=/%s", s.config.Url, s.config.RepoID, filePath)

	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}
	req.Header.Set("Authorization", "Token "+s.token)
	req.Header.Set("Accept", "application/json; charset=utf-8; indent=4")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute delete request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed: %s", string(b))
	}

	// laut API kommt einfach "success" als Antwort
	var result string
	if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
		if result != "success" {
			return fmt.Errorf("delete returned unexpected response: %s", result)
		}
	}

	return nil
}

func (s *StorageService) getUploadUrl(parentDir string) (string, error) {
	// Erstelle die URL für den Upload-Link
	uploadLinkURL := fmt.Sprintf("%s/repos/%s/upload-link/?p=%s", s.config.Url, s.config.RepoID, parentDir)

	// Erstelle die HTTP-Anfrage
	req, err := http.NewRequest("GET", uploadLinkURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create upload-link request: %w", err)
	}
	req.Header.Set("Authorization", "Token "+s.token)

	// Führe die Anfrage aus
	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to request upload link: %w", err)
	}
	defer resp.Body.Close()

	// Überprüfe den HTTP-Statuscode
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload link request failed: %s", string(b))
	}

	// Dekodiere die Antwort
	var uploadURL string
	if err := json.NewDecoder(resp.Body).Decode(&uploadURL); err != nil {
		return "", fmt.Errorf("failed to decode upload-link response: %w", err)
	}

	return uploadURL, nil
}
