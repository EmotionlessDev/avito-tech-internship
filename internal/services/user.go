package services

import (
	"context"

	"github.com/EmotionlessDev/avito-tech-internship/internal/models"
	"github.com/EmotionlessDev/avito-tech-internship/internal/repository"
)

type UserService interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error)
}

type userService struct {
	users repository.UserRepository
}

func NewUserService(users repository.UserRepository) UserService {
	return &userService{users: users}
}

func (s *userService) SetIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	return s.users.UpdateIsActive(ctx, userID, isActive)
}
