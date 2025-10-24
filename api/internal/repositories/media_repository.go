package repositories

import (
	"embox/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) MediaRepository {
	return &mediaRepository{db}
}

func (r *mediaRepository) Create(media *models.Media) error {
	return r.db.Create(media).Error
}

func (r *mediaRepository) Update(media *models.Media) error {
	return r.db.Save(media).Error
}

func (r *mediaRepository) Delete(ids []uint) error {
	return r.db.Where("id IN ?", ids).Delete(&models.Media{}).Error
}

func (r *mediaRepository) Get(isAdmin bool, userId uuid.UUID) ([]*models.Media, error) {
	var media []*models.Media

	query := r.db.
		Table("media AS m").
		Select(`
            m.*,
            CASE WHEN fav.user_id IS NOT NULL THEN true ELSE false END AS is_favourite
        `).
		Joins("LEFT JOIN favourites AS fav ON fav.media_id = m.id AND fav.user_id = ?", userId).
		Order("m.date DESC")

	if !isAdmin {
		query = query.Where("m.is_public = ?", true)
	}

	err := query.Find(&media).Error

	return media, err
}

func (r *mediaRepository) GetById(id uint) (*models.Media, error) {
	var media models.Media
	if err := r.db.First(&media, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No record found
		}
		return nil, err // Query error
	}
	return &media, nil
}

func (r *mediaRepository) GetByIDs(ids []uint) ([]*models.Media, error) {
	var media []*models.Media
	if err := r.db.Where("id IN ?", ids).Find(&media).Error; err != nil {
		return nil, err
	}
	return media, nil
}
