package infrastructure

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/christhianjesus/crabi-challenge/internal/domain"
	"github.com/christhianjesus/crabi-challenge/internal/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type userHandlerMock struct {
	service *mocks.UserService
	handler *userHandler
	rec     *httptest.ResponseRecorder
	ctx     echo.Context
}

func setupGetUserHandler(t *testing.T) *userHandlerMock {
	mockUserService := mocks.NewUserService(t)
	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)
	ctx.Set("user_id", "1")

	return &userHandlerMock{
		service: mockUserService,
		handler: NewUserHandler(mockUserService),
		rec:     rec,
		ctx:     ctx,
	}
}

func TestGetUser_OK(t *testing.T) {
	user := &domain.User{ID: "1", Email: "an@email.com"}
	guh := setupGetUserHandler(t)
	guh.service.On("GetUser", mock.Anything, mock.AnythingOfType("string")).Return(user, nil)

	err := guh.handler.Get(guh.ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, guh.rec.Code)
	assert.JSONEq(t, `{"id":"1", "email":"an@email.com", "first_name":"", "last_name":"", "created_at":"0001-01-01T00:00:00Z", "updated_at":"0001-01-01T00:00:00Z"}`, guh.rec.Body.String())
}

func TestGetUser_GetUserError(t *testing.T) {
	guh := setupGetUserHandler(t)
	guh.service.On("GetUser", mock.Anything, mock.AnythingOfType("string")).Return(nil, assert.AnError)

	err := guh.handler.Get(guh.ctx)
	he := err.(*echo.HTTPError)

	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, he.Code)
	assert.Equal(t, assert.AnError.Error(), he.Message)
}
