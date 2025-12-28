package services

import (
	"embox/internal/api/dto"
	"embox/internal/repositories"
	"fmt"
)

type FavouriteService struct {
	userRepo repositories.UserRepository
	favRepo  repositories.FavouriteRepository
}

func NewFavouriteService(userRepo repositories.UserRepository, favRepo repositories.FavouriteRepository) *FavouriteService {
	return &FavouriteService{userRepo, favRepo}
}

// Add one or more media to the user's favourites
func (s *FavouriteService) AddFavourites(userEmail string, mediaIds []uint) error {
	user, err := s.userRepo.GetByEmail(userEmail)
	if err != nil || user == nil {
		return fmt.Errorf("user not found")
	}
	return s.favRepo.Add(user.ID, mediaIds)
}

// Remove one or more media from the user's favourites
func (s *FavouriteService) RemoveFavourites(userEmail string, mediaIds []uint) error {
	user, err := s.userRepo.GetByEmail(userEmail)
	if err != nil || user == nil {
		return fmt.Errorf("user not found")
	}
	return s.favRepo.Remove(user.ID, mediaIds)
}

// GetUsersWithLatestFavourite retrieves all users and their latest favourite media.
func (s *FavouriteService) GetUsersWithLatestFavourite(userEmail string) ([]dto.UserWithLatestFavouriteDto, error) {
	user, err := s.userRepo.GetByEmail(userEmail)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}

	results, err := s.favRepo.GetUsersWithLatestFavourite()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users with latest favourites: %w", err)
	}

	responseDtos := make([]dto.UserWithLatestFavouriteDto, len(results))

	for i, result := range results {
		responseDtos[i] = dto.UserWithLatestFavouriteDto{
			User: struct {
				ID         string `json:"id"`
				Name       string `json:"name"`
				MediaCount int    `json:"media_count"`
			}{
				ID:         result.UserID,
				Name:       result.UserName,
				MediaCount: result.MediaCount,
			},
			Media: struct {
				ID      uint   `json:"id"`
				Caption string `json:"caption"`
				Type    string `json:"type"`
			}{
				ID:      result.MediaID,
				Caption: result.MediaCaption,
				Type:    result.MediaType,
			},
		}
	}

	return responseDtos, nil
}

func (s *FavouriteService) GetFavouritesByUserID(favUserID, userEmail string) (dto.FavouritesResponseDto, error) {
	user, err := s.userRepo.GetByEmail(userEmail)
	if err != nil || user == nil {
		return dto.FavouritesResponseDto{}, fmt.Errorf("user not found")
	}

	results, err := s.favRepo.GetFavouritesByUserID(user.ID, favUserID)
	if err != nil {
		return dto.FavouritesResponseDto{}, fmt.Errorf("failed to fetch favourites for user %s: %w", favUserID, err)
	}

	if len(results) == 0 {
		return dto.FavouritesResponseDto{
			User: dto.MediaUserResponseDto{
				ID:   favUserID,
				Name: "",
			},
			Media: []dto.MediaResponseDto{},
		}, nil
	}

	userDto := dto.MediaUserResponseDto{
		ID:   results[0].FavouriteUserID,
		Name: results[0].FavouriteUserName,
	}

	mediaDtos := make([]dto.MediaResponseDto, len(results))
	for i, fav := range results {
		mediaDtos[i] = dto.MediaResponseDto{
			Id:          fav.ID,
			IsFavourite: fav.IsFavourite,
			Caption:     fav.Caption,
			Date:        fav.Date.Format("2006-01-02"),
			Type:        fav.Type,
			CreatedAt:   fav.CreatedAt,
		}
	}

	return dto.FavouritesResponseDto{
		User:  userDto,
		Media: mediaDtos,
	}, nil
}
