package main

import (
	"forum/auth"
	"forum/config"
	"forum/db/gorm"
	"forum/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	config.Init("./config")

	gorm.InitDB(
		config.Viper.GetString("database.user"),
		config.Viper.GetString("database.password"),
		config.Viper.GetString("database.dbname"),
	)

	e := echo.New()
	validator := validator.InitValidator()
	e.Validator = &validator
	e.Logger.SetLevel(log.INFO)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "frontend/build",
		Browse: true,
		HTML5:  true,
	}))

	apiGroup := e.Group("/api")

	// Auth route
	authHandler := auth.Handler{DB: gorm.DB}
	authGroup := apiGroup.Group("/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/logout", authHandler.Logout)

	// Posts group
	postsGroup := apiGroup.Group("/posts")
	postsGroup.Use(middleware.JWT([]byte("secret")))

	e.Logger.Fatal(e.Start(":1323"))
}
