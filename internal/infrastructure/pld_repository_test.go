package infrastructure

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/christhianjesus/crabi-challenge/internal/domain"
	"github.com/christhianjesus/crabi-challenge/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type pldRepositoryMock struct {
	client *mocks.HTTPClient
	repo   domain.PLDRepository
}

func setupPLDRepository(t *testing.T) *pldRepositoryMock {
	mockHTTPClient := mocks.NewHTTPClient(t)

	return &pldRepositoryMock{
		client: mockHTTPClient,
		repo:   NewPLDRepository(mockHTTPClient, ""),
	}
}

func TestIsValidUser_MarshalError(t *testing.T) {
	prm := setupPLDRepository(t)

	isValidUser, err := prm.repo.IsValidUser(context.Context(nil), nil)

	assert.Empty(t, isValidUser)
	assert.Error(t, err)
	assert.EqualError(t, err, "json: error calling MarshalJSON for type *infrastructure.pldRequest: Nil user")
}

func TestIsValidUser_NewRequestError(t *testing.T) {
	prm := setupPLDRepository(t)

	isValidUser, err := prm.repo.IsValidUser(context.Context(nil), &domain.User{})

	assert.Empty(t, isValidUser)
	assert.Error(t, err)
	assert.EqualError(t, err, "net/http: nil Context")
}

func TestIsValidUser_DoError(t *testing.T) {
	prm := setupPLDRepository(t)
	prm.client.On("Do", mock.Anything).
		Return(nil, assert.AnError)

	isValidUser, err := prm.repo.IsValidUser(context.TODO(), &domain.User{})

	assert.Empty(t, isValidUser)
	assert.Error(t, err)
	assert.EqualError(t, err, assert.AnError.Error())
}

type FailRead struct{}

func (*FailRead) Read(p []byte) (n int, err error) {
	return 0, assert.AnError
}

func TestIsValidUser_ReadAllError(t *testing.T) {
	prm := setupPLDRepository(t)
	prm.client.On("Do", mock.Anything).
		Return(&http.Response{
			StatusCode: 201,
			Body:       io.NopCloser(&FailRead{}),
		}, nil)

	isValidUser, err := prm.repo.IsValidUser(context.TODO(), &domain.User{})

	assert.Empty(t, isValidUser)
	assert.Error(t, err)
	assert.EqualError(t, err, assert.AnError.Error())
}

func TestIsValidUser_UnmarshalError(t *testing.T) {
	prm := setupPLDRepository(t)
	prm.client.On("Do", mock.Anything).
		Return(&http.Response{
			StatusCode: 201,
			Body:       io.NopCloser(strings.NewReader(``)),
		}, nil)

	isValidUser, err := prm.repo.IsValidUser(context.TODO(), &domain.User{})

	assert.Empty(t, isValidUser)
	assert.Error(t, err)
	assert.EqualError(t, err, "unexpected end of JSON input")
}

func TestIsValidUser_OK(t *testing.T) {
	prm := setupPLDRepository(t)
	prm.client.On("Do", mock.Anything).
		Return(&http.Response{
			StatusCode: 201,
			Body:       io.NopCloser(strings.NewReader(`{"is_in_blacklist": true}`)),
		}, nil)

	isValidUser, err := prm.repo.IsValidUser(context.TODO(), &domain.User{})

	assert.NoError(t, err)
	assert.Equal(t, false, isValidUser)
}
