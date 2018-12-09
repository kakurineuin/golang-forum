package main

import (
	"github.com/kakurineuin/golang-forum/admin"
	"github.com/kakurineuin/golang-forum/auth"
	"github.com/kakurineuin/golang-forum/config"
	"github.com/kakurineuin/golang-forum/database"
	forumError "github.com/kakurineuin/golang-forum/error"
	"github.com/kakurineuin/golang-forum/logger"
	forumMiddleware "github.com/kakurineuin/golang-forum/middleware"
	"github.com/kakurineuin/golang-forum/post"
	"github.com/kakurineuin/golang-forum/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	config.Init("./config", "config")

	dao := database.InitDAO(
		config.Viper.GetString("database.user"),
		config.Viper.GetString("database.password"),
		config.Viper.GetString("database.dbname"),
	)

	e := echo.New()
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		c.Logger().Error(err)

		if customError, ok := err.(forumError.CustomError); ok {
			c.JSON(customError.HTTPStatusCode, map[string]interface{}{
				"message": customError.Message,
			})
			return
		}

		e.DefaultHTTPErrorHandler(err, c)
	}

	validator := validator.InitValidator()
	e.Validator = &validator

	// Logger
	myLogger := logger.InitLogger()
	e.Logger = myLogger
	e.Logger.SetLevel(log.INFO)

	// Middleware
	e.Use(logger.Middleware(myLogger))
	e.Use(middleware.Recover())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "frontend/build",
		Browse: true,
		HTML5:  true,
	}))

	apiGroup := e.Group("/api")

	// Auth route
	authService := auth.Service{DAO: dao}
	authHandler := auth.Handler{Service: &authService}
	authGroup := apiGroup.Group("/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)

	// Posts route
	postService := post.Service{DAO: dao}
	postHandler := post.Handler{Service: &postService}
	postsGroup := apiGroup.Group("/topics")
	postsGroup.GET("/statistics", postHandler.FindTopicsStatistics)
	postsGroup.GET("/:category/:id", postHandler.FindTopic)
	postsGroup.GET("/:category", postHandler.FindTopics)
	jwtMiddleware := middleware.JWT([]byte(auth.JwtSecret))
	postsGroup.POST("/:category", postHandler.CreatePost, jwtMiddleware)
	postsGroup.PUT("/:category/:id", postHandler.UpdatePost, jwtMiddleware)
	postsGroup.DELETE("/:category/:id", postHandler.DeletePost, jwtMiddleware)

	// 查詢論壇統計資料。
	apiGroup.GET("/forum/statistics", postHandler.FindForumStatistics)

	// Admin
	adminService := admin.Service{DAO: dao}
	adminHandler := admin.Handler{Service: &adminService}
	apiGroup.GET("/users", adminHandler.FindUsers, jwtMiddleware, forumMiddleware.Admin)

	e.Logger.Fatal(e.Start(":1323"))
}
