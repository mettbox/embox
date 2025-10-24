package repositories

import (
	"embox/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uuid.UUID) error
	GetAll() ([]*models.User, error)
	GetById(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByToken(token string, validDuration time.Duration) (*models.User, error)
	GenerateToken(email string) (*models.User, error)
}

type MediaRepository interface {
	Create(media *models.Media) error
	Update(media *models.Media) error
	Delete(ids []uint) error
	Get(isAdmin bool, userId uuid.UUID) ([]*models.Media, error)
	GetById(id uint) (*models.Media, error)
	GetByIDs(ids []uint) ([]*models.Media, error)
}

type FavouriteRepository interface {
	Add(userId uuid.UUID, mediaIds []uint) error
	Remove(userId uuid.UUID, mediaIds []uint) error
	GetUsersWithLatestFavourite(isAdmin bool) ([]UserWithLatestFavourite, error)
	GetFavouritesByUserID(isAdmin bool, userId uuid.UUID, favUserID string) ([]*models.Media, error)
}

type AlbumRepository interface {
	Create(album *models.Album) error
	Update(album *models.Album) error
	Delete(id uint) error
	Get(isAdmin bool) ([]*models.Album, error)
	GetById(id uint) (*models.Album, error)
	GetMediaIdsByAlbumId(albumId uint) ([]uint, error)
	AddMediaToAlbum(albumId uint, mediaIds []uint, isCover bool) error
	RemoveMediaFromAlbum(albumId uint, mediaIds []uint) error
}

type Repositories struct {
	User      UserRepository
	Media     MediaRepository
	Favourite FavouriteRepository
	Album     AlbumRepository
}

// Repository responses

type UserWithLatestFavourite struct {
	UserID       string `json:"user_id"`
	UserName     string `json:"user_name"`
	MediaID      uint   `json:"media_id"`
	MediaCaption string `json:"caption"`
	MediaType    string `json:"type"`
	MediaCount   int    `json:"media_count"`
}

// Init initializes the repositories with the provided database connection.
// It returns a Repositories struct containing all the repositories.
func Init(db *gorm.DB) *Repositories {
	return &Repositories{
		User:      NewUserRepository(db),
		Media:     NewMediaRepository(db),
		Favourite: NewFavouriteRepository(db),
		Album:     NewAlbumRepository(db),
	}
}
