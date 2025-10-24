package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                  uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name                string    `gorm:"type:varchar(48);null"`
	Email               string    `gorm:"type:varchar(128);not null;unique"`
	IsAdmin             bool      `gorm:"default:false"`
	HasPublicFavourites *bool     `gorm:"default:true"`
	Token               string    `gorm:"type:varchar(80);null"`
	TokenCreatedAt      time.Time `gorm:"default:null"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Favourites          []Favourite `gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
