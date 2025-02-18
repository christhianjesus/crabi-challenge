package infrastructure

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const ValidatorCtxKey = "validator"

// Helper function to set user_id context variable
func SetUserID(c echo.Context) {
	// by default token is stored under `user` key
	claims := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)
	c.Set("user_id", claims["user_id"])
}

func SetValidator(next echo.HandlerFunc) echo.HandlerFunc {
	validate := validator.New()

	// register function to get tag name from json tags.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return func(c echo.Context) error {
		c.Set(ValidatorCtxKey, validate)

		return next(c)
	}
}
