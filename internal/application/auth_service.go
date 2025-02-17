package application

import (
	"context"

	"github.com/christhianjesus/crabi-challenge/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Signin(ctx context.Context, user *domain.User) error
	Login(ctx context.Context, email, password string) (string, error)
}

type authService struct {
	repo    domain.AuthRepository
	userSrv UserService
}

func NewAuthService(repo domain.AuthRepository, userSrv UserService) AuthService {
	return &authService{repo, userSrv}
}

func (a *authService) Signin(ctx context.Context, user *domain.User) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedBytes)

	return a.userSrv.CreateUser(ctx, user)
}

func (a *authService) Login(ctx context.Context, email string, password string) (string, error) {
	userID, hash, err := a.repo.GetIdAndHash(ctx, email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return "", err
	}

	return userID, nil
}
