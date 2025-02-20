package infrastructure

import (
	"net/http"

	"github.com/christhianjesus/crabi-challenge/internal/application"
	"github.com/christhianjesus/crabi-challenge/internal/domain"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type authHandler struct {
	srv application.AuthService
	sec string
}

func NewAuthHandler(srv application.AuthService, secret string) *authHandler {
	return &authHandler{srv, secret}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *authHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	request := new(LoginRequest)

	if err := c.Bind(request); err != nil {
		return err
	}

	if validate, ok := c.Get(ValidatorCtxKey).(*validator.Validate); ok {
		if err := validate.Struct(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	userID, err := h.srv.Login(ctx, request.Email, request.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
	})

	tokenString, err := token.SignedString([]byte(h.sec))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, &LoginResponse{Token: tokenString})
}

func (h *authHandler) Signin(c echo.Context) error {
	ctx := c.Request().Context()
	request := new(domain.User)

	if err := c.Bind(request); err != nil {
		return err
	}

	if validate, ok := c.Get(ValidatorCtxKey).(*validator.Validate); ok {
		if err := validate.Struct(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	err := h.srv.Signin(ctx, request)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}
