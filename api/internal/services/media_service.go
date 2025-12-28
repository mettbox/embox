package services

import (
	"bytes"
	"embox/internal/api/dto"
	"embox/internal/models"
	"embox/internal/repositories"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

type MediaService struct {
	storage   *StorageService
	mediaRepo repositories.MediaRepository
	userRepo  repositories.UserRepository
}

var MediaDir = "./media"
var imgMaxSize = 512
var imgQuality float32 = 80

func NewMediaService(storage *StorageService, mediaRepo repositories.MediaRepository, userRepo repositories.UserRepository) *MediaService {
	return &MediaService{storage, mediaRepo, userRepo}
}

// === public functions ===

// GetMediaList retrieves media items based on order and user permissions.
func (s *MediaService) GetMediaList(userEmail string) ([]dto.MediaResponseDto, error) {
	user, err := s.userRepo.GetByEmail(userEmail)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}

	mediaList, err := s.mediaRepo.Get(user.ID)
	if err != nil {
		return nil, err
	}

	var results []dto.MediaResponseDto
	for _, media := range mediaList {
		results = append(results, dto.MediaResponseDto{
			Id:          media.ID,
			IsFavourite: media.IsFavourite,
			Caption:     media.Caption,
			Date:        media.Date.Format(time.RFC3339),
			Type:        media.Type,
			CreatedAt:   media.CreatedAt,
		})
	}

	return results, nil
}

func (s *MediaService) GetMediaByID(id uint, userEmail string) (*models.Media, error) {
	user, err := s.userRepo.GetByEmail(userEmail)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}

	media, err := s.mediaRepo.GetById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve media: %w", err)
	}
	if media == nil {
		return nil, nil // Media not found
	}

	return media, nil
}

// GetThumbnail retrieves the media thumbnail data and its MIME type by media ID.
func (s *MediaService) GetThumbnail(id uint) ([]byte, string, error) {
	media, err := s.mediaRepo.GetById(id)
	if err != nil {
		return nil, "", fmt.Errorf("failed to find media with id %d: %w", id, err)
	}
	if media == nil {
		return nil, "", fmt.Errorf("media with id %d not found", id)
	}

	data, mimeType, err := s.loadMediaFromDisk(media.Path())
	if err != nil {
		return nil, "", fmt.Errorf("failed to load media thumbnail: %w", err)
	}

	return data, mimeType, nil
}

// GetMediaFile retrieves the full media file data and its MIME type by media ID.
func (s *MediaService) GetMediaFile(id uint) ([]byte, string, error) {
	media, err := s.mediaRepo.GetById(id)
	if err != nil {
		return nil, "", fmt.Errorf("failed to find media with id %d: %w", id, err)
	}
	if media == nil {
		return nil, "", fmt.Errorf("media with id %d not found", id)
	}

	data, mimeType, err := s.storage.Download(media.RemotePath())
	if err != nil {
		return nil, "", fmt.Errorf("failed to download media file: %w", err)
	}

	return data, mimeType, nil
}

// Creates a new media entry from the provided metadata and file data.
func (s *MediaService) CreateFromRequest(meta dto.MediaUploadRequestDto, file io.Reader, userEmail string) (*models.Media, error) {
	user, err := s.userRepo.GetByEmail(userEmail)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Parse the date string
	// First try RFC3339 ISO (with timezone, e.g. from EXIF or toISOString())
	// Then try ISO without timezone (common from ion-datetime)
	parsedDate, err := time.Parse(time.RFC3339, meta.Date)
	if err != nil {
		parsedDate, err = time.Parse("2006-01-02T15:04:05", meta.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %w", err)
		}
	}

	media := &models.Media{
		UserID:    &user.ID,
		Caption:   meta.Caption,
		Type:      getMediaType(meta.Type),
		FileExt:   getFileExt(meta.FileName), // Original-Endung behalten
		Date:      parsedDate,
		CreatedAt: time.Now(),
	}

	if err := s.mediaRepo.Create(media); err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err := s.storage.Upload(bytes, media.RemotePath()); err != nil {
		return nil, fmt.Errorf("failed to upload original file: %w", err)
	}

	if media.Type == "image" || media.Type == "video" {
		if media.Type == "image" {
			bytes, err = s.convertToWebP(bytes)
			if err != nil {
				return nil, err
			}
		}

		if media.Type == "video" {
			bytes, err = s.generateVideoPoster(media, bytes)
			if err != nil {
				return nil, err
			}
		}

		if err := s.saveMediaFile(media, bytes); err != nil {
			return nil, err
		}
	}

	return media, nil
}

// UploadMedia handles the upload of multiple media files with their metadata.
func (s *MediaService) UploadMedia(files []*multipart.FileHeader, metaList []dto.MediaUploadRequestDto, userEmail string) ([]dto.MediaResponseDto, error) {
	user, err := s.userRepo.GetByEmail(userEmail)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if len(files) != len(metaList) {
		return nil, fmt.Errorf("meta and files count mismatch")
	}

	var uploaded []dto.MediaResponseDto

	for i, fileHeader := range files {
		meta := metaList[i]

		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		media, err := s.CreateFromRequest(meta, file, userEmail)
		if err != nil {
			return nil, fmt.Errorf("failed to save media: %w", err)
		}

		uploaded = append(uploaded, dto.MediaResponseDto{
			Id:          media.ID,
			IsFavourite: false,
			Caption:     media.Caption,
			Date:        media.Date.Format(time.RFC3339),
			Type:        media.Type,
			CreatedAt:   media.CreatedAt,
		})
	}

	return uploaded, nil
}

// UpdateMediaBatch updates multiple media items based on the provided update requests.
func (s *MediaService) UpdateMediaBatch(updates []dto.MediaUpdateRequestDto, userEmail string) ([]dto.MediaResponseDto, error) {
	user, err := s.userRepo.GetByEmail(userEmail)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}

	var updatedMedia []dto.MediaResponseDto
	var updateErrors []string

	for _, update := range updates {
		existingMedia, err := s.mediaRepo.GetById(update.ID)
		if err != nil {
			updateErrors = append(updateErrors, fmt.Sprintf("failed to retrieve media ID %d: %v", update.ID, err))
			continue
		}
		if existingMedia == nil {
			updateErrors = append(updateErrors, fmt.Sprintf("media ID %d not found", update.ID))
			continue
		}

		if update.Caption != nil {
			existingMedia.Caption = *update.Caption
		}

		if update.Date != nil {
			parsedDate, err := time.Parse(time.RFC3339, *update.Date)
			if err != nil {
				parsedDate, err = time.Parse("2006-01-02T15:04:05", *update.Date)
				if err != nil {
					updateErrors = append(updateErrors, fmt.Sprintf("invalid date format for media ID %d: %v", update.ID, err))
					continue
				}
			}
			existingMedia.Date = parsedDate
		}

		if err := s.mediaRepo.Update(existingMedia); err != nil {
			updateErrors = append(updateErrors, fmt.Sprintf("failed to update media ID %d: %v", update.ID, err))
			continue
		}

		updatedMedia = append(updatedMedia, dto.MediaResponseDto{
			Id:          existingMedia.ID,
			IsFavourite: existingMedia.IsFavourite,
			Caption:     existingMedia.Caption,
			Date:        existingMedia.Date.Format(time.RFC3339),
			Type:        existingMedia.Type,
			CreatedAt:   existingMedia.CreatedAt,
		})
	}

	if len(updateErrors) > 0 {
		return updatedMedia, fmt.Errorf("some updates failed: %v", updateErrors)
	}

	return updatedMedia, nil
}

// DeleteMedia deletes media items by their IDs, removing both local and remote files.
func (s *MediaService) DeleteMedia(ids []uint) error {
	mediaList, err := s.mediaRepo.GetByIDs(ids)
	if err != nil {
		return err
	}

	for _, media := range mediaList {
		localFilePath := filepath.Join(MediaDir, media.Path())
		_ = os.Remove(localFilePath)

		if err := s.storage.Delete(media.RemotePath()); err != nil {
			fmt.Printf("Failed to delete remote file %s: %v\n", media.RemotePath(), err)
		}
	}

	return s.mediaRepo.Delete(ids)
}

// === private functions ===

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

// convertToWebP converts the given image data to WebP format with specified size and quality.
func (s *MediaService) convertToWebP(data []byte) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// 1. Handle Orientation
	needsProcessing := false
	exifData, err := exif.Decode(bytes.NewReader(data))
	if err == nil {
		orientTag, err := exifData.Get(exif.Orientation)
		if err == nil {
			orient, err := orientTag.Int(0)
			if err == nil && orient > 1 {
				img = applyOrientation(img, orient)
				needsProcessing = true
			}
		}
	}

	// 2. Handle Resizing (Width must be at most imgMaxSize, height automatic)
	if img.Bounds().Dx() > imgMaxSize {
		img = imaging.Resize(img, imgMaxSize, 0, imaging.Lanczos)
		needsProcessing = true
	}

	// 3. Skip re-encoding if it's already WebP and no changes were made
	if !needsProcessing && format == "webp" {
		return data, nil
	}

	// 4. Encode to WebP
	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, &webp.Options{Lossless: false, Quality: imgQuality}); err != nil {
		return nil, fmt.Errorf("failed to encode webp: %w", err)
	}

	return buf.Bytes(), nil
}

func applyOrientation(img image.Image, orient int) image.Image {
	switch orient {
	case 3:
		return imaging.Rotate180(img)
	case 6:
		return imaging.Rotate270(img)
	case 8:
		return imaging.Rotate90(img)
	default:
		return img
	}
}
