package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Location struct {
	Latitude  float64 `gorm:"type:double;not null"`
	Longitude float64 `gorm:"type:double;not null"`
}

type Media struct {
	ID          uint       `gorm:"type:int;primaryKey"`
	IsPublic    bool       `gorm:"not null;default:false" json:"isPublic"`
	Date        time.Time  `gorm:"type:date;not null"`
	UserID      *uuid.UUID `gorm:"type:char(36);null"`                             // Foreign Key, nullable
	User        User       `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"` // Relation
	FileExt     string     `gorm:"type:varchar(8);not null"`
	Type        string     `gorm:"type:varchar(8);not null"`
	Caption     string     `gorm:"type:varchar(255);null"`
	Location    Location   `gorm:"embedded;embeddedPrefix:location_"`
	Orientation string     `gorm:"type:varchar(32);null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Computed field, nur für SELECT
	// Erläuterung:
	// -> bedeutet: read-only (nur beim SELECT berücksichtigen, nicht beim INSERT/UPDATE).
	// column:is_favourite mappt korrekt auf den SQL-Alias aus deinem Query.
	IsFavourite bool `gorm:"->;column:is_favourite" json:"isFavourite"`

	// Addiotional Read-only fields for the Favorit-User
	FavouriteUserID   string `gorm:"->;column:favourite_user_id" json:"favourite_user_id"`
	FavouriteUserName string `gorm:"->;column:favourite_user_name" json:"favourite_user_name"`

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
