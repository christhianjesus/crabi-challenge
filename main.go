package main

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	c := GetContext()

	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Health
	e.GET("/health", nil)

	// Auth routes
	e.POST("/signin", nil)
	e.POST("/login", nil)

	// versioning endpoints
	v1 := e.Group("/v1", echojwt.JWT(c.GetJwtKey()))

	// App routes
	v1.GET("/user", nil)

	e.Logger.Fatal(e.Start(":" + c.GetHttpPort()))
}
