package main

import (
	"context"
	"log"
	"time"

	"github.com/christhianjesus/crabi-challenge/internal/application"
	"github.com/christhianjesus/crabi-challenge/internal/infrastructure"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	c := GetContext()

	mongoClient, err := mongo.Connect(options.Client().ApplyURI(c.GetMongoURL()))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Repositories
	mongoUserRepository := infrastructure.NewMongoUserRepository(mongoClient.Database("default"))

	// Services
	userService := application.NewUserService(mongoUserRepository)
	authService := application.NewAuthService(mongoUserRepository, userService)

	// Handlers
	userHandler := infrastructure.NewUserHandler(userService)
	authHandler := infrastructure.NewAuthHandler(authService, c.jwtKey)

	// Middlewares
	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SuccessHandler: authHandler.SetUserID,
		SigningKey:     []byte(c.jwtKey),
	})

	e := echo.New()

	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Health
	e.GET("/health", nil)

	// Auth routes
	e.POST("/signin", authHandler.Signin)
	e.POST("/login", authHandler.Login)

	// versioning endpoints
	v1 := e.Group("/v1", jwtMiddleware)

	// App routes
	v1.GET("/user", userHandler.Get)

	e.Logger.Fatal(e.Start(":" + c.GetHttpPort()))
}
