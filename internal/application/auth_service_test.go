package application

import (
	"context"
	"testing"

	"github.com/christhianjesus/crabi-challenge/internal/domain"
	"github.com/christhianjesus/crabi-challenge/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type authServiceMock struct {
	repoMock *mocks.AuthRepository
	srvMock  *mocks.UserService
	service  AuthService
}

func setupAuthService(t *testing.T) *authServiceMock {
	mockAuthRepository := mocks.NewAuthRepository(t)
	mockUserService := mocks.NewUserService(t)

	return &authServiceMock{
		repoMock: mockAuthRepository,
		srvMock:  mockUserService,
		service:  NewAuthService(mockAuthRepository, mockUserService),
	}
}

func TestSignin_OK(t *testing.T) {
	password := "123"
	checkPassword := func(user *domain.User) bool {
		return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil
	}

	asm := setupAuthService(t)
	asm.srvMock.On("CreateUser", mock.IsType(nil), mock.MatchedBy(checkPassword)).Return(nil)

	err := asm.service.Signin(context.Context(nil), &domain.User{Password: password})

	assert.NoError(t, err)
}

func TestSignin_GenerateFromPasswordError(t *testing.T) {
	password := make([]byte, 100)
	asm := setupAuthService(t)

	err := asm.service.Signin(context.Context(nil), &domain.User{Password: string(password)})

	assert.Error(t, err)
	assert.EqualError(t, err, bcrypt.ErrPasswordTooLong.Error())
}

func TestSignin_CreateUserError(t *testing.T) {
	asm := setupAuthService(t)
	asm.srvMock.On("CreateUser", mock.IsType(nil), mock.AnythingOfType("*domain.User")).Return(assert.AnError)

	err := asm.service.Signin(context.Context(nil), &domain.User{})

	assert.Error(t, err)
	assert.EqualError(t, err, assert.AnError.Error())
}

func TestLogin_OK(t *testing.T) {
	res := "1"
	password := "123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	asm := setupAuthService(t)
	asm.repoMock.On("GetIdAndHash", mock.IsType(nil), mock.AnythingOfType("string")).Return(res, string(hash), nil)

	userID, err := asm.service.Login(context.Context(nil), "an@email.com", password)

	assert.NoError(t, err)
	assert.NotNil(t, userID)
	assert.Equal(t, res, userID)
}

func TestLogin_GetIdAndHashError(t *testing.T) {
	asm := setupAuthService(t)
	asm.repoMock.On("GetIdAndHash", mock.IsType(nil), mock.AnythingOfType("string")).Return("", "", assert.AnError)

	userID, err := asm.service.Login(context.Context(nil), "an@email.com", "")

	assert.Error(t, err)
	assert.EqualError(t, err, assert.AnError.Error())
	assert.Empty(t, userID)
}

func TestLogin_CompareHashAndPasswordError(t *testing.T) {
	emptyHash := ""

	asm := setupAuthService(t)
	asm.repoMock.On("GetIdAndHash", mock.IsType(nil), mock.AnythingOfType("string")).Return("1", emptyHash, nil)

	userID, err := asm.service.Login(context.Context(nil), "an@email.com", "123")

	assert.Error(t, err)
	assert.EqualError(t, err, bcrypt.ErrHashTooShort.Error())
	assert.Empty(t, userID)
}
