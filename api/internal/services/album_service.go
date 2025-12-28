package services

import (
	"embox/internal/api/dto"
	"embox/internal/models"
	"embox/internal/repositories"
	"fmt"
)

type AlbumService struct {
	userRepo  repositories.UserRepository
	albumRepo repositories.AlbumRepository
}

func NewAlbumService(userRepo repositories.UserRepository, albumRepo repositories.AlbumRepository) *AlbumService {
	return &AlbumService{userRepo, albumRepo}
}

func (s *AlbumService) CreateAlbum(album *dto.CreateAlbumRequestDto, userEmail string) (*dto.AlbumResponseDto, error) {
	user, err := s.userRepo.GetByEmail(userEmail)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}

	newAlbum := &dto.AlbumResponseDto{
		Name:        album.Name,
		Description: album.Description,
		MediaCount:  len(album.MediaIDs),
	}

	albumModel := &models.Album{
		Name:        album.Name,
		Description: album.Description,
		UserID:      &user.ID,
	}

	for _, mediaID := range album.MediaIDs {
		albumModel.AlbumMedia = append(albumModel.AlbumMedia, models.AlbumMedia{
			MediaID: mediaID,
			IsCover: mediaID == album.CoverMediaID,
		})
	}

	err = s.albumRepo.Create(albumModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create album: %w", err)
	}

	newAlbum.Id = albumModel.ID

	return newAlbum, nil
}

func (s *AlbumService) UpdateAlbum(id uint, album *dto.UpdateAlbumRequestDto) (*dto.AlbumResponseDto, error) {
	existingAlbum, err := s.albumRepo.GetById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch album: %w", err)
	}

	if album.Name != "" {
		existingAlbum.Name = album.Name
	}
	if album.Description != "" {
		existingAlbum.Description = album.Description
	}

	if err := s.albumRepo.Update(existingAlbum); err != nil {
		return nil, fmt.Errorf("failed to update album: %w", err)
	}
	// return updated existingAlbum album as AlbumResponseDto
	result := &dto.AlbumResponseDto{
		Id:          existingAlbum.ID,
		Name:        existingAlbum.Name,
		Description: existingAlbum.Description,
		MediaCount:  len(existingAlbum.AlbumMedia),
		Media:       []dto.AlbumMediaResponseDto{},
	}

	return result, nil

}

func (s *AlbumService) GetAlbumList(userEmail string) ([]dto.AlbumResponseDto, error) {
	user, err := s.userRepo.GetByEmail(userEmail)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}

	albums, err := s.albumRepo.Get()
	if err != nil {
		fmt.Printf("Error fetching albums: %v\n", err)
		return nil, err
	}

	var result []dto.AlbumResponseDto
	for _, album := range albums {
		var mediaDtos []dto.AlbumMediaResponseDto
		for _, albumMedia := range album.AlbumMedia {
			mediaDto := dto.AlbumMediaResponseDto{
				Id:          albumMedia.MediaID,
				IsFavourite: albumMedia.Media.IsFavourite,
				Caption:     albumMedia.Media.Caption,
				Date:        albumMedia.Media.Date.Format("2006-01-02"),
				Type:        albumMedia.Media.Type,
				IsCover:     albumMedia.IsCover,
				CreatedAt:   albumMedia.Media.CreatedAt,
			}
			mediaDtos = append(mediaDtos, mediaDto)
		}

		albumDto := dto.AlbumResponseDto{
			Id:          album.ID,
			Name:        album.Name,
			Description: album.Description,
			MediaCount:  album.MediaCount,
			Media:       mediaDtos,
		}
		result = append(result, albumDto)
	}

	return result, nil
}

func (s *AlbumService) GetAlbumByID(id uint) (*dto.AlbumResponseDto, error) {
	album, err := s.albumRepo.GetById(id)
	if err != nil {
		fmt.Printf("Error fetching album by ID: %v\n", err)
		return nil, err
	}

	var mediaDtos []dto.AlbumMediaResponseDto
	for _, albumMedia := range album.AlbumMedia {
		mediaDto := dto.AlbumMediaResponseDto{
			Id:          albumMedia.MediaID,
			IsFavourite: albumMedia.Media.IsFavourite,
			Caption:     albumMedia.Media.Caption,
			Date:        albumMedia.Media.Date.Format("2006-01-02"),
			Type:        albumMedia.Media.Type,
			IsCover:     albumMedia.IsCover,
			CreatedAt:   albumMedia.Media.CreatedAt,
		}
		mediaDtos = append(mediaDtos, mediaDto)
	}

	albumDto := &dto.AlbumResponseDto{
		Id:          album.ID,
		Name:        album.Name,
		Description: album.Description,
		MediaCount:  len(mediaDtos),
		Media:       mediaDtos,
	}

	return albumDto, nil
}

func (s *AlbumService) DeleteAlbum(id uint) error {
	return s.albumRepo.Delete(id)
}

func (s *AlbumService) AddMediaToAlbum(albumId uint, mediaIds []uint, isCover bool) error {
	existingMediaIds, err := s.albumRepo.GetMediaIdsByAlbumId(albumId)
	if err != nil {
		return fmt.Errorf("failed to fetch existing media IDs: %w", err)
	}

	var newMediaIds []uint
	existingMediaSet := make(map[uint]bool)
	for _, id := range existingMediaIds {
		existingMediaSet[id] = true
	}

	for _, id := range mediaIds {
		if !existingMediaSet[id] {
			newMediaIds = append(newMediaIds, id)
		}
	}

	if len(newMediaIds) == 0 {
		return nil
	}

	return s.albumRepo.AddMediaToAlbum(albumId, newMediaIds, isCover)
}

func (s *AlbumService) RemoveMediaFromAlbum(albumId uint, mediaIds []uint) error {
	return s.albumRepo.RemoveMediaFromAlbum(albumId, mediaIds)
}
