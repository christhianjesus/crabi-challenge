package application

import (
	"context"
	"errors"

	"github.com/christhianjesus/crabi-challenge/internal/domain"
)

type UserService interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUser(ctx context.Context, userID string) (*domain.User, error)
}

type userService struct {
	repo    domain.UserRepository
	pldRepo domain.PLDRepository
}

func NewUserService(repo domain.UserRepository, pldRepo domain.PLDRepository) UserService {
	return &userService{repo, pldRepo}
}

func (u *userService) CreateUser(ctx context.Context, user *domain.User) error {
	valid, err := u.pldRepo.IsValidUser(ctx, user)
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("User is in blacklist")
	}

	return u.repo.CreateUser(ctx, user)
}

func (u *userService) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	return u.repo.GetUser(ctx, userID)
}
