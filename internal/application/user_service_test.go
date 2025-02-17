package application

import (
	"context"
	"testing"

	"github.com/christhianjesus/crabi-challenge/internal/domain"
	"github.com/christhianjesus/crabi-challenge/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type userServiceMock struct {
	repo    *mocks.UserRepository
	service UserService
}

func setupUserService(t *testing.T) *userServiceMock {
	mockUserRepository := mocks.NewUserRepository(t)

	return &userServiceMock{
		repo:    mockUserRepository,
		service: NewUserService(mockUserRepository),
	}
}

func TestCreateUser_OK(t *testing.T) {
	usm := setupUserService(t)
	usm.repo.On("CreateUser", mock.IsType(nil), mock.AnythingOfType("*domain.User")).Return(nil)

	err := usm.service.CreateUser(context.Context(nil), &domain.User{})

	assert.NoError(t, err)
}

func TestCreateUser_Error(t *testing.T) {
	usm := setupUserService(t)
	usm.repo.On("CreateUser", mock.IsType(nil), mock.AnythingOfType("*domain.User")).Return(assert.AnError)

	err := usm.service.CreateUser(context.Context(nil), &domain.User{})

	assert.Error(t, err)
	assert.EqualError(t, err, assert.AnError.Error())
}

func TestGetUser_OK(t *testing.T) {
	res := &domain.User{ID: "1"}

	usm := setupUserService(t)
	usm.repo.On("GetUser", mock.IsType(nil), mock.AnythingOfType("string")).Return(res, nil)

	user, err := usm.service.GetUser(context.Context(nil), "1")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, res.ID, user.ID)
}

func TestGetUser_Error(t *testing.T) {
	usm := setupUserService(t)
	usm.repo.On("GetUser", mock.IsType(nil), mock.AnythingOfType("string")).Return(nil, assert.AnError)

	user, err := usm.service.GetUser(context.Context(nil), "")

	assert.Error(t, err)
	assert.EqualError(t, err, assert.AnError.Error())
	assert.Nil(t, user)
}
