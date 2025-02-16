package main

import (
	"context"
	"log"
	"time"

	"github.com/christhianjesus/crabi-challenge/internal/infrastructure"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	c := GetContext()
	mongoOptions := options.Client().ApplyURI(c.GetMongoURL())

	mongoClient, err := mongo.Connect(mongoOptions)
	if err != nil {
		log.Fatal(err)
	}

	mongoDB := mongoClient.Database("default")

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// TODO
	_ = infrastructure.NewMongoUserRepository(mongoDB)

	e := echo.New()

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
