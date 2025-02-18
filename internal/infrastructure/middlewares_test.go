package infrastructure

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSetUserID_OK(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/v1/any", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSJ9.TcXz_IwlmxO5nPd3m0Yo67WyYptabkqZW4R9HNwPmKE")

	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SuccessHandler: SetUserID,
		SigningKey:     []byte("secret"),
	})

	err := jwtMiddleware(func(c echo.Context) error {
		return nil
	})(ctx)

	assert.NoError(t, err)
	assert.IsType(t, "", ctx.Get("user_id"))
	assert.Equal(t, "1", ctx.Get("user_id").(string))
}

func TestSetValidator_OK(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/any", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	ctxMiddleware := SetValidator

	err := ctxMiddleware(func(c echo.Context) error {
		return nil
	})(ctx)

	assert.NoError(t, err)
	assert.IsType(t, new(validator.Validate), ctx.Get(ValidatorCtxKey))
	assert.NotEmpty(t, ctx.Get(ValidatorCtxKey).(*validator.Validate))
}

func TestSetValidator_IgnoreNameTag(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/any", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	ctxMiddleware := SetValidator

	err := ctxMiddleware(func(c echo.Context) error {
		validate := c.Get(ValidatorCtxKey).(*validator.Validate)
		err := validate.Struct(struct {
			DefaultName string `json:"-" validate:"required"`
		}{})

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	})(ctx)

	assert.Error(t, err)
	assert.EqualError(t, err, "code=400, message=Key: 'DefaultName' Error:Field validation for 'DefaultName' failed on the 'required' tag")
}
