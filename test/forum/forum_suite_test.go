package forum_test

import (
	"context"
	"github.com/kakurineuin/golang-forum/handler"
	"github.com/kakurineuin/golang-forum/service"
	"testing"
	"time"

	"github.com/kakurineuin/golang-forum/config"
	"github.com/kakurineuin/golang-forum/database"
	"github.com/kakurineuin/golang-forum/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var e *echo.Echo
var dao *database.DAO
var forumHandler handler.ForumHandler

func TestPost(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Post Suite")
}

var _ = BeforeSuite(func() {
	config.Init("../../config", "test")

	dao = database.InitDAO(
		config.Viper.GetString("database.user"),
		config.Viper.GetString("database.password"),
		config.Viper.GetString("database.host"),
		config.Viper.GetString("database.dbname"),
	)

	forumService := service.ForumService{
		DAO: dao,
	}
	forumHandler = handler.ForumHandler{ForumService: &forumService}

	e = echo.New()
	validator := validator.InitValidator()
	e.Validator = &validator
	e.Logger.SetLevel(log.INFO)

	// Middleware
	e.Use(middleware.Logger())
})

var _ = AfterSuite(func() {
	dao.DB.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
})
