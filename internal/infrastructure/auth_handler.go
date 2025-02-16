package infrastructure

import (
	"net/http"

	"github.com/christhianjesus/crabi-challenge/internal/application"
	"github.com/christhianjesus/crabi-challenge/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type authHandler struct {
	srv application.AuthService
	sec []byte
}

func NewAuthHandler(srv application.AuthService, secret []byte) *authHandler {
	return &authHandler{srv, secret}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Helper function to set user_id context variable
func (h *authHandler) SetUserID(c echo.Context) {
	// by default token is stored under `user` key
	claims := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)
	c.Set("user_id", claims["user_id"])
}

func (h *authHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	request := new(LoginRequest)

	if err := c.Bind(request); err != nil {
		return err
	}

	userID, err := h.srv.Login(ctx, request.Email, request.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
	})

	tokenString, err := token.SignedString(h.sec)
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

	err := h.srv.Signin(ctx, request)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}
