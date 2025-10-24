package models

import (
	"time"

	"github.com/google/uuid"
)

type Favourite struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uuid.UUID `gorm:"type:char(36);not null;index;uniqueIndex:idx_user_media"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	MediaID   uint      `gorm:"not null;index;uniqueIndex:idx_user_media"`
	Media     Media     `gorm:"foreignKey:MediaID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `gorm:"type:datetime(3);not null;index"`
}

func (Favourite) TableName() string {
	return "favourites"
}
