package infrastructure

import (
	"crypto"
	"crypto/sha256"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/christhianjesus/crabi-challenge/internal/mocks"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type authHandlerMock struct {
	service *mocks.AuthService
	handler *authHandler
}

func setupAuthHandler(t *testing.T) *authHandlerMock {
	mockAuthService := mocks.NewAuthService(t)

	return &authHandlerMock{
		service: mockAuthService,
		handler: NewAuthHandler(mockAuthService, "secret"),
	}
}

func TestLogin_OK(t *testing.T) {
	body := strings.NewReader(`{"email": "an@email.com", "password": "123"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)

	userID := "1"
	lg := setupAuthHandler(t)
	lg.service.On("Login", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(userID, nil)

	err := lg.handler.Login(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.TcXz_IwlmxO5nPd3m0Yo67WyYptabkqZW4R9HNwPmKE"}`, rec.Body.String())
}

func TestLogin_BindError(t *testing.T) {
	body := strings.NewReader(`{`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)

	lg := setupAuthHandler(t)

	err := lg.handler.Login(ctx)
	he := err.(*echo.HTTPError)

	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, he.Code)
	assert.Equal(t, "unexpected EOF", he.Message)
}

func TestLogin_LoginError(t *testing.T) {
	body := strings.NewReader(`{"email": "an@email.com", "password": "123"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)

	userID := ""
	lg := setupAuthHandler(t)
	lg.service.On("Login", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(userID, assert.AnError)

	err := lg.handler.Login(ctx)
	he := err.(*echo.HTTPError)

	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, he.Code)
	assert.Equal(t, assert.AnError.Error(), he.Message)
}

func TestLogin_SignedStringError(t *testing.T) {
	body := strings.NewReader(`{"email": "an@email.com", "password": "123"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)

	userID := "1"
	lg := setupAuthHandler(t)
	lg.service.On("Login", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(userID, nil)

	crypto.RegisterHash(crypto.SHA256, nil)
	err := lg.handler.Login(ctx)
	crypto.RegisterHash(crypto.SHA256, sha256.New)
	he := err.(*echo.HTTPError)

	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, he.Code)
	assert.Equal(t, jwt.ErrHashUnavailable.Error(), he.Message)
}

func TestSignin_OK(t *testing.T) {
	body := strings.NewReader(`{"email": "an@email.com", "password": "123", "first_name": "first_name", "last_name": "last_name"}`)
	req := httptest.NewRequest(http.MethodPost, "/signin", body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)

	lg := setupAuthHandler(t)
	lg.service.On("Signin", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	err := lg.handler.Signin(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Empty(t, rec.Body)
}

func TestSignin_BindError(t *testing.T) {
	body := strings.NewReader(`{`)
	req := httptest.NewRequest(http.MethodPost, "/signin", body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)

	lg := setupAuthHandler(t)

	err := lg.handler.Signin(ctx)
	he := err.(*echo.HTTPError)

	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, he.Code)
	assert.Equal(t, "unexpected EOF", he.Message)
}

func TestSignin_SigninError(t *testing.T) {
	body := strings.NewReader(`{"email": "an@email.com", "password": "123", "first_name": "first_name", "last_name": "last_name"}`)
	req := httptest.NewRequest(http.MethodPost, "/signin", body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)

	lg := setupAuthHandler(t)
	lg.service.On("Signin", mock.Anything, mock.AnythingOfType("*domain.User")).Return(assert.AnError)

	err := lg.handler.Signin(ctx)
	he := err.(*echo.HTTPError)

	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, he.Code)
	assert.Equal(t, assert.AnError.Error(), he.Message)
}

func TestSetUserID_OK(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/v1/any", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.TcXz_IwlmxO5nPd3m0Yo67WyYptabkqZW4R9HNwPmKE")

	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	authHandler := NewAuthHandler(nil, "secret")
	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SuccessHandler: authHandler.SetUserID,
		SigningKey:     []byte("secret"),
	})

	h := jwtMiddleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	h(ctx)

	assert.Equal(t, "1", ctx.Get("user_id").(string))
}
