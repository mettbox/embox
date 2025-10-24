package services

import (
	"embox/internal/api/dto"
	"embox/internal/models"
	"embox/internal/repositories"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type MediaService struct {
	storage   *StorageService
	mediaRepo repositories.MediaRepository
	userRepo  repositories.UserRepository
}

var MediaDir = "./media"
var imgMaxSize = 384
var imgQuality float32 = 80

func NewMediaService(storage *StorageService, mediaRepo repositories.MediaRepository, userRepo repositories.UserRepository) *MediaService {
	return &MediaService{storage, mediaRepo, userRepo}
}

// GetMediaList retrieves media items based on order and user permissions.
func (s *MediaService) GetMediaList(userEmail string) ([]dto.MediaResponseDto, error) {
	user, err := s.userRepo.GetByEmail(userEmail)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}

	mediaList, err := s.mediaRepo.Get(user.IsAdmin, user.ID)
	if err != nil {
		return nil, err
	}

	var results []dto.MediaResponseDto
	for _, media := range mediaList {
		// var loc dto.Location
		// if media.Location.Latitude != 0 || media.Location.Longitude != 0 {
		// 	loc = dto.Location{Lat: media.Location.Latitude, Lng: media.Location.Longitude}
		// }
		results = append(results, dto.MediaResponseDto{
			Id:          media.ID,
			IsPublic:    media.IsPublic,
			IsFavourite: media.IsFavourite,
			Caption:     media.Caption,
			Date:        media.Date.Format("2006-01-02"),
			Type:        media.Type,
			// Location:    loc,
			CreatedAt: media.CreatedAt,
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

	// Check permissions
	if !media.IsPublic && (media.UserID == nil || *media.UserID != user.ID) && !user.IsAdmin {
		return nil, fmt.Errorf("access denied")
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

	parsedDate, err := time.Parse("2006-01-02", meta.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	media := &models.Media{
		UserID:    &user.ID,
		Caption:   meta.Caption,
		Type:      getMediaType(meta.Type),
		FileExt:   getFileExt(meta.FileName), // Original-Endung behalten
		Date:      parsedDate,
		IsPublic:  meta.IsPublic,
		CreatedAt: time.Now(),
	}

	if meta.LocationLat != 0 || meta.LocationLng != 0 {
		media.Location = models.Location{
			Latitude:  meta.LocationLat,
			Longitude: meta.LocationLng,
		}
	}
	if meta.Orientation != "" {
		media.Orientation = meta.Orientation
	}

	if err := s.mediaRepo.Create(media); err != nil {
		return nil, err
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err := s.storage.Upload(fileData, media.RemotePath()); err != nil {
		return nil, fmt.Errorf("failed to upload original file: %w", err)
	}

	if media.Type == "image" || media.Type == "video" {
		if media.Type == "image" {
			fileData, err = s.convertToWebP(fileData)
			if err != nil {
				return nil, err
			}
		}

		if media.Type == "video" {
			fileData, err = s.generateVideoPoster(media, fileData)
			if err != nil {
				return nil, err
			}
		}

		if err := s.saveMediaFile(media, fileData); err != nil {
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
			IsPublic:    media.IsPublic,
			IsFavourite: false,
			Caption:     media.Caption,
			Date:        media.Date.Format("2006-01-02"),
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

		if !existingMedia.IsPublic && (existingMedia.UserID == nil || *existingMedia.UserID != user.ID) && !user.IsAdmin {
			updateErrors = append(updateErrors, fmt.Sprintf("access denied for media ID %d", update.ID))
			continue
		}

		if update.Caption != nil {
			existingMedia.Caption = *update.Caption
		}
		if update.IsPublic != nil {
			existingMedia.IsPublic = *update.IsPublic
		}
		if update.Date != nil {
			parsedDate, err := time.Parse("2006-01-02", *update.Date)
			if err != nil {
				updateErrors = append(updateErrors, fmt.Sprintf("invalid date format for media ID %d: %v", update.ID, err))
				continue
			}
			existingMedia.Date = parsedDate
		}
		if update.Location != nil {
			existingMedia.Location = models.Location{
				Latitude:  update.Location.Lat,
				Longitude: update.Location.Lng,
			}
		}

		if err := s.mediaRepo.Update(existingMedia); err != nil {
			updateErrors = append(updateErrors, fmt.Sprintf("failed to update media ID %d: %v", update.ID, err))
			continue
		}

		updatedMedia = append(updatedMedia, dto.MediaResponseDto{
			Id:          existingMedia.ID,
			IsPublic:    existingMedia.IsPublic,
			IsFavourite: existingMedia.IsFavourite,
			Caption:     existingMedia.Caption,
			Date:        existingMedia.Date.Format("2006-01-02"),
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
