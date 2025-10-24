package infrastructure

import (
	"embox/internal/config"
	"embox/internal/models"
	"fmt"
	"log"
	"net/url"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDatabase initializes the database connection and performs migrations.
// It also ensures that a system user is created if it does not exist.
func InitDatabase(config *config.DbConfig) (*gorm.DB, error) {
	password := url.QueryEscape(config.DBPass)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser, password, config.DBHost, config.DBPort, config.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{}, &models.Media{}, &models.Favourite{}, &models.Album{}, &models.AlbumMedia{})
	if err != nil {
		return nil, err
	}

	var count int64
	db.Model(&models.User{}).Where("email = ?", config.DBSystemUser).Count(&count)
	if count == 0 {
		systemUser := models.User{
			Name:    "System",
			Email:   config.DBSystemUser,
			IsAdmin: true,
		}
		if err := db.Create(&systemUser).Error; err != nil {
			log.Printf("Failed to create system user (%s): %v", config.DBSystemUser, err)
		} else {
			log.Printf("System user created: %s", config.DBSystemUser)
		}
	}

	return db, nil
}
