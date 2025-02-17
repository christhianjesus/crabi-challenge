package infrastructure

import (
	"net/http"

	"github.com/christhianjesus/crabi-challenge/internal/application"
	"github.com/labstack/echo/v4"
)

type userHandler struct {
	srv application.UserService
}

func NewUserHandler(srv application.UserService) *userHandler {
	return &userHandler{srv}
}

func (h *userHandler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Get("user_id").(string)

	user, err := h.srv.GetUser(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}
