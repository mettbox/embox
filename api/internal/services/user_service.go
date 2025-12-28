package services

import (
	"embox/internal/api/dto"
	"embox/internal/models"
	"embox/internal/repositories"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{userRepo}
}

// GenerateToken generates a token for the user with the given email
func (s *UserService) GenerateToken(email string) (*models.User, error) {
	user, err := s.userRepo.GenerateToken(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) CreateUser(user *models.User) (*dto.UserResponse, error) {
	err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	userResponse := &dto.UserResponse{
		ID:          user.ID.String(),
		Name:        user.Name,
		Email:       user.Email,
		IsAdmin:     user.IsAdmin,
		LastLoginAt: user.LastLoginAt,
	}

	return userResponse, nil
}

// GetAllUsers returns all users
func (s *UserService) GetAllUsers() ([]dto.UserResponse, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}

	var results []dto.UserResponse
	for _, user := range users {
		results = append(results, dto.UserResponse{
			ID:          user.ID.String(),
			Name:        user.Name,
			Email:       user.Email,
			IsAdmin:     user.IsAdmin,
			LastLoginAt: user.LastLoginAt,
		})
	}

	return results, nil
}

// GetByToken returns a user by token if the token is valid and not expired
func (s *UserService) GetByToken(token string, validDuration time.Duration) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByToken(token, validDuration)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	result := &dto.UserResponse{
		ID:          user.ID.String(),
		Name:        user.Name,
		Email:       user.Email,
		IsAdmin:     user.IsAdmin,
		LastLoginAt: user.LastLoginAt,
	}

	return result, nil
}

// GetUserById returns a user by ID
func (s *UserService) GetUserById(id uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetById(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	result := &dto.UserResponse{
		ID:          user.ID.String(),
		Name:        user.Name,
		Email:       user.Email,
		IsAdmin:     user.IsAdmin,
		LastLoginAt: user.LastLoginAt,
	}

	return result, nil
}

// GetUserByEmail returns a user by email
func (s *UserService) GetUserByEmail(email string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	result := &dto.UserResponse{
		ID:          user.ID.String(),
		Name:        user.Name,
		Email:       user.Email,
		IsAdmin:     user.IsAdmin,
		LastLoginAt: user.LastLoginAt,
	}

	return result, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	existingUser, err := s.userRepo.GetById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}
	if existingUser == nil {
		return nil, fmt.Errorf("user not found")
	}

	if req.Name != nil {
		existingUser.Name = *req.Name
	}
	if req.Email != nil {
		existingUser.Email = *req.Email
	}
	if req.IsAdmin != nil {
		existingUser.IsAdmin = *req.IsAdmin
	}
	if req.LastLoginAt != nil {
		existingUser.LastLoginAt = req.LastLoginAt
	}

	if err := s.userRepo.Update(existingUser); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	result := &dto.UserResponse{
		ID:          existingUser.ID.String(),
		Name:        existingUser.Name,
		Email:       existingUser.Email,
		IsAdmin:     existingUser.IsAdmin,
		LastLoginAt: existingUser.LastLoginAt,
	}

	return result, nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(id uuid.UUID) error {
	return s.userRepo.Delete(id)
}
