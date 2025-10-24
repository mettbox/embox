package repositories

import (
	"embox/internal/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type favouriteRepository struct {
	db *gorm.DB
}

func NewFavouriteRepository(db *gorm.DB) FavouriteRepository {
	return &favouriteRepository{db}
}

func (r *favouriteRepository) Add(userId uuid.UUID, mediaIds []uint) error {
	for _, mediaId := range mediaIds {
		fav := models.Favourite{
			UserID:    userId,
			MediaID:   mediaId,
			CreatedAt: time.Now(),
		}
		if err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&fav).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *favouriteRepository) Remove(userId uuid.UUID, mediaIds []uint) error {
	return r.db.Where("user_id = ? AND media_id IN ?", userId, mediaIds).Delete(&models.Favourite{}).Error
}

func (r *favouriteRepository) GetUsersWithLatestFavourite(isAdmin bool) ([]UserWithLatestFavourite, error) {
	var results []UserWithLatestFavourite

	// Subquery: Get latest favourite per user firstly
	subQuery := r.db.
		Table("favourites").
		Select("user_id, MAX(created_at) as latest_created_at").
		Group("user_id")

	// Main query: Join with users and media to get the required details
	query := r.db.
		Table("favourites").
		Select(`
      users.id as user_id,
      users.name as user_name,
      media.id as media_id,
      media.caption as media_caption,
      media.type as media_type,
      (
          SELECT COUNT(*)
          FROM favourites
          WHERE favourites.user_id = users.id
      ) as media_count
    `).
		Joins("JOIN users ON users.id = favourites.user_id").
		Joins("JOIN media ON media.id = favourites.media_id").
		Joins("JOIN (?) as latest_favourites ON favourites.user_id = latest_favourites.user_id AND favourites.created_at = latest_favourites.latest_created_at", subQuery).
		Where("users.has_public_favourites = ?", true)

	if !isAdmin {
		query = query.Where("media.is_public = ?", true)
	}

	query = query.
		Group("users.id, users.name, media.id, media.caption, media.type").
		Order("latest_favourites.latest_created_at DESC")

	err := query.Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users with latest favourites: %w", err)
	}

	return results, nil
}

func (r *favouriteRepository) GetFavouritesByUserID(isAdmin bool, userId uuid.UUID, favUserID string) ([]*models.Media, error) {
	var media []*models.Media

	query := r.db.
		Table("favourites AS f").
		Select(`
            m.*,
            f.user_id as favourite_user_id,
            u.name as favourite_user_name,
            CASE WHEN fav.user_id IS NOT NULL THEN true ELSE false END AS is_favourite
        `).
		Joins("JOIN media AS m ON m.id = f.media_id").
		Joins("JOIN users AS u ON u.id = f.user_id").
		Joins("LEFT JOIN favourites AS fav ON fav.media_id = m.id AND fav.user_id = ?", userId).
		Where("f.user_id = ?", favUserID).
		Order("m.date DESC")

	if !isAdmin {
		query = query.Where("m.is_public = ?", true)
	}

	err := query.Scan(&media).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch favourites for user %s: %w", favUserID, err)
	}

	return media, nil
}
