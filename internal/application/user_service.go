package application

import (
	"context"

	"github.com/christhianjesus/crabi-challenge/internal/domain"
)

type UserService interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUser(ctx context.Context, userID string) (*domain.User, error)
}

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) UserService {
	return &userService{repo}
}

func (u *userService) CreateUser(ctx context.Context, user *domain.User) error {
	return u.repo.CreateUser(ctx, user)
}

func (u *userService) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	return u.repo.GetUser(ctx, userID)
}
