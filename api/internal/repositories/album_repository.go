package repositories

import (
	"embox/internal/models"

	"gorm.io/gorm"
)

type albumRepository struct {
	db *gorm.DB
}

func NewAlbumRepository(db *gorm.DB) AlbumRepository {
	return &albumRepository{db}
}

func (r *albumRepository) Create(album *models.Album) error {
	return r.db.Create(album).Error
}

func (r *albumRepository) Update(album *models.Album) error {
	return r.db.Save(album).Error
}

func (r *albumRepository) Get() ([]*AlbumListItem, error) {
	var albums []*AlbumListItem

	err := r.db.
		Preload("Media").
		Preload("AlbumMedia.Media").
		Preload("AlbumMedia", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_cover = ?", true).
				Or("album_media.media_id = (SELECT media_id FROM album_media am WHERE am.album_id = album_media.album_id ORDER BY media_id ASC LIMIT 1)")
		}).
		Select(`
        albums.*,
        (
            SELECT COUNT(*)
            FROM album_media
            WHERE album_media.album_id = albums.id
        ) as media_count
    `).
		Find(&albums).Error

	if err != nil {
		return nil, err
	}

	if albums == nil {
		albums = []*AlbumListItem{}
	}

	return albums, nil
}

func (r *albumRepository) GetById(id uint) (*models.Album, error) {
	{
		var album models.Album

		err := r.db.
			Preload("Media").
			Preload("AlbumMedia.Media").
			Preload("AlbumMedia", func(db *gorm.DB) *gorm.DB {
				return db.Joins("JOIN media ON media.id = album_media.media_id").
					Order("media.date DESC")
			}).
			First(&album, id).Error

		if err != nil {
			return nil, err
		}

		return &album, nil
	}
}

func (r *albumRepository) GetMediaIdsByAlbumId(albumId uint) ([]uint, error) {
	var mediaIds []uint

	err := r.db.
		Table("album_media").
		Where("album_id = ?", albumId).
		Pluck("media_id", &mediaIds).Error

	if err != nil {
		return nil, err
	}

	return mediaIds, nil
}

func (r *albumRepository) Delete(id uint) error {
	return r.db.Delete(&models.Album{}, id).Error
}

func (r *albumRepository) AddMediaToAlbum(albumId uint, mediaIds []uint, isCover bool) error {
	var albumMediaEntries []models.AlbumMedia
	for _, mediaId := range mediaIds {
		albumMediaEntries = append(albumMediaEntries, models.AlbumMedia{
			AlbumID: albumId,
			MediaID: mediaId,
			IsCover: isCover,
		})
	}
	return r.db.Create(&albumMediaEntries).Error
}

func (r *albumRepository) RemoveMediaFromAlbum(albumId uint, mediaIds []uint) error {
	return r.db.Where("album_id = ? AND media_id IN ?", albumId, mediaIds).Delete(&models.AlbumMedia{}).Error
}
