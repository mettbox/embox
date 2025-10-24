package models

import (
	"time"

	"github.com/google/uuid"
)

type Album struct {
	ID          uint       `gorm:"type:int;primaryKey"`
	IsPublic    bool       `gorm:"not null;default:false" json:"isPublic"`
	Name        string     `gorm:"type:varchar(255);not null"`
	Description string     `gorm:"type:text;null"`
	UserID      *uuid.UUID `gorm:"type:char(36);null"`                             // Foreign Key, nullable
	User        User       `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"` // Relation
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Many-to-Many Relation with Media
	Media      []Media      `gorm:"many2many:album_media;constraint:OnDelete:CASCADE;"`
	AlbumMedia []AlbumMedia `gorm:"foreignKey:AlbumID"`

	// Computed field for media count
	MediaCount int `json:"mediaCount"`
}

type AlbumMedia struct {
	AlbumID uint  `gorm:"primaryKey"`
	MediaID uint  `gorm:"primaryKey"`
	IsCover bool  `gorm:"default:false"`
	Album   Album `gorm:"foreignKey:AlbumID;constraint:OnDelete:CASCADE;"`
	Media   Media `gorm:"foreignKey:MediaID;constraint:OnDelete:CASCADE;"`
}
