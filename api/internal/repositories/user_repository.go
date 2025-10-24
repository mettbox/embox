package repositories

import (
	"embox/internal/models"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetById(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByToken(token string, validDuration time.Duration) (*models.User, error) {
	var user models.User
	if err := r.db.Where("token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}

	// Check if the token is expired
	if time.Since(user.TokenCreatedAt) > validDuration {
		if err := r.removeToken(&user); err != nil {
			return nil, err
		}

		return nil, errors.New("token expired")
	}

	if err := r.removeToken(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) removeToken(user *models.User) error {
	user.Token = ""
	user.TokenCreatedAt = time.Time{}
	if err := r.db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetAll() ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *userRepository) GenerateToken(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	token := generateUniqueToken(r.db)
	if token == "" {
		return nil, errors.New("failed to generate unique token")
	}

	user.Token = token
	user.TokenCreatedAt = time.Now()

	if err := r.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func generateUniqueToken(db *gorm.DB) string {
	for range 10 { // Try up to 10 times to generate a unique token
		token := generateToken()
		var count int64
		db.Model(&models.User{}).Where("token = ?", token).Count(&count)
		if count == 0 {
			return token
		}
	}
	return ""
}

func generateToken() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", r.Intn(1000000))
}
