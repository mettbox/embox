package services

import (
	"embox/internal/models"
	"embox/pkg/webpconv"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// loadMediaFromDisk reads the media file from disk and returns its data and MIME type.
func (s *MediaService) loadMediaFromDisk(path string) ([]byte, string, error) {
	filePath := filepath.Join(MediaDir, path)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", err
	}
	mimeType := http.DetectContentType(data)
	return data, mimeType, nil
}

// getMediaType determines the media type based on the MIME type.
func getMediaType(mime string) string {
	if strings.HasPrefix(mime, "image") {
		return "image"
	}
	if strings.HasPrefix(mime, "video") {
		return "video"
	}
	if strings.HasPrefix(mime, "audio") {
		return "audio"
	}
	return "other"
}

// saveMediaFile saves the media file data to disk at the specified path.
func (s *MediaService) saveMediaFile(media *models.Media, data []byte) error {
	savePath := filepath.Join(MediaDir, media.Path())
	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		return err
	}
	return os.WriteFile(savePath, data, 0644)
}

// getFileExt extracts the file extension from the given file name.
func getFileExt(fileName string) string {
	parts := strings.Split(fileName, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return ""
}

// convertToWebP converts the given image data to WebP format with specified size and quality.
func (s *MediaService) convertToWebP(data []byte) ([]byte, error) {
	bytes, err := webpconv.ConvertToWebP(data, imgMaxSize, imgQuality)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// generateVideoPoster extracts a poster image from the video and converts it to WebP format.
func (s *MediaService) generateVideoPoster(media *models.Media, videoData []byte) ([]byte, error) {
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("media_%d_orig.%s", media.ID, media.FileExt))
	if err := os.WriteFile(tmpFile, videoData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write temp video: %w", err)
	}
	defer os.Remove(tmpFile)

	tmpPoster := tmpFile + "_poster.png"
	cmd := exec.Command("ffmpeg",
		"-i", tmpFile,
		"-vf", fmt.Sprintf("scale='min(%d,iw)':-2", imgMaxSize),
		"-vframes", "1",
		tmpPoster,
	)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to extract video poster: %w", err)
	}
	defer os.Remove(tmpPoster)

	posterData, err := os.ReadFile(tmpPoster)
	if err != nil {
		return nil, fmt.Errorf("failed to read poster image: %w", err)
	}

	bytes, err := s.convertToWebP(posterData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert poster to webp: %w", err)
	}

	return bytes, nil
}
