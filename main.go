package main

import (
	"github.com/kakurineuin/golang-forum/config"
	"github.com/kakurineuin/golang-forum/database"
	forumError "github.com/kakurineuin/golang-forum/error"
	"github.com/kakurineuin/golang-forum/handler"
	"github.com/kakurineuin/golang-forum/logger"
	forumMiddleware "github.com/kakurineuin/golang-forum/middleware"
	"github.com/kakurineuin/golang-forum/service"
	"github.com/kakurineuin/golang-forum/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	config.Init("./config", "development")

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

	jwtSecret := config.Viper.GetString("jwt.secret")
	apiGroup := e.Group("/api")

	// Auth route
	authService := service.AuthService{DAO: dao}
	authHandler := handler.AuthHandler{
		AuthService: &authService,
		JwtSecret:   jwtSecret,
	}
	authGroup := apiGroup.Group("/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)

	// Topics route
	topicService := service.TopicService{DAO: dao}
	topicHandler := handler.TopicHandler{TopicService: &topicService}
	topicsGroup := apiGroup.Group("/topics")
	topicsGroup.GET("/statistics", topicHandler.FindTopicsStatistics)
	topicsGroup.GET("/:category/:id", topicHandler.FindTopic)
	topicsGroup.GET("/:category", topicHandler.FindTopics)
	jwtMiddleware := middleware.JWT([]byte(jwtSecret))
	topicsGroup.POST("/:category", topicHandler.CreatePost, jwtMiddleware)
	topicsGroup.PUT("/:category/:id", topicHandler.UpdatePost, jwtMiddleware)
	topicsGroup.DELETE("/:category/:id", topicHandler.DeletePost, jwtMiddleware)

	// Forum route
	forumService := service.ForumService{DAO: dao}
	forumHandler := handler.ForumHandler{ForumService: &forumService}
	forumGroup := apiGroup.Group("/forum")
	forumGroup.GET("/statistics", forumHandler.FindForumStatistics)

	// Admin
	adminService := service.AdminService{DAO: dao}
	adminHandler := handler.AdminHandler{AdminService: &adminService}
	adminGroup := apiGroup.Group("/admin")
	adminGroup.Use(jwtMiddleware, forumMiddleware.Admin)
	adminGroup.GET("/users", adminHandler.FindUsers)
	adminGroup.POST("/users/disable/:id", adminHandler.DisableUser)

	e.Logger.Fatal(e.Start(":1323"))
}
