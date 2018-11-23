package main

import (
	"forum/auth"
	"forum/config"
	"forum/db/gorm"
	"forum/post"
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

	// Posts route
	postHandler := post.Handler{DB: gorm.DB}
	postsGroup := apiGroup.Group("/topics")
	postsGroup.GET("/statistics", postHandler.FindTopicsStatistics)
	postsGroup.GET("/:category/:id", postHandler.FindTopic)
	postsGroup.GET("/:category", postHandler.FindTopics)
	jwtMiddleware := middleware.JWT([]byte(auth.JwtSecret))
	postsGroup.POST("/:category", postHandler.CreatePost, jwtMiddleware)
	postsGroup.PUT("/:category/:id", postHandler.UpdatePost, jwtMiddleware)

	e.Logger.Fatal(e.Start(":1323"))
}
