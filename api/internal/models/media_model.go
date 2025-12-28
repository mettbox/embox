package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Media struct {
	ID        uint       `gorm:"type:int;primaryKey"`
	Date      time.Time  `gorm:"type:date;not null"`
	UserID    *uuid.UUID `gorm:"type:char(36);null"`                             // Foreign Key, nullable
	User      User       `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"` // Relation
	FileExt   string     `gorm:"type:varchar(8);not null"`
	Type      string     `gorm:"type:varchar(8);not null"`
	Caption   string     `gorm:"type:varchar(255);null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Computed fields, ignored by GORM for DB operations
	IsFavourite       bool   `gorm:"-" json:"isFavourite"`
	FavouriteUserID   string `gorm:"-" json:"favourite_user_id"`
	FavouriteUserName string `gorm:"-" json:"favourite_user_name"`

	// Many-to-Many Relation with Album
	Albums []Album `gorm:"many2many:album_media;constraint:OnDelete:CASCADE;"`
}

func (Media) TableName() string {
	return "media"
}

// Path returns the local path of the media item.
// For images: yyyy/mm/dd_Id.webp
// For audio/video: yyyy/mm/dd_Id.FileExt
func (m *Media) Path() string {
	dateStr := m.Date.Format("2006/01/02") // yyyy/mm/dd
	ext := m.FileExt
	mediaType := strings.ToLower(m.Type)
	if mediaType == "image" || mediaType == "video" {
		ext = "webp"
	}
	return fmt.Sprintf("%s_%d.%s", dateStr, m.ID, ext)
}

// RemotePath always returns: yyyy/mm/dd_Id.FileExt
func (m *Media) RemotePath() string {
	dateStr := m.Date.Format("2006/01/02") // yyyy/mm/dd
	return fmt.Sprintf("%s_%d.%s", dateStr, m.ID, m.FileExt)
}
