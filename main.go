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
	authGroup.POST("/logout", authHandler.Logout)

	// Posts route
	postHandler := post.Handler{DB: gorm.DB}
	postsGroup := apiGroup.Group("/posts")
	postsGroup.GET("/statistics", postHandler.FindPostsStatistics)
	postsGroup.GET("/:category/topics/:id", postHandler.FindPostsTopics)
	postsGroup.GET("/:category", postHandler.FindPosts)
	jwtMiddleware := middleware.JWT([]byte(auth.JwtSecret))
	postsGroup.POST("/:category", postHandler.CreatePost, jwtMiddleware)

	e.Logger.Fatal(e.Start(":1323"))
}
