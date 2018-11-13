package post_test

import (
	"context"
	"forum/config"
	"forum/db/gorm"
	"forum/post"
	"forum/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

var e *echo.Echo
var handler post.Handler

func TestPost(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Post Suite")
}

var _ = BeforeSuite(func() {
	config.Init("../config")

	gorm.InitDB(
		config.Viper.GetString("database.user"),
		config.Viper.GetString("database.password"),
		config.Viper.GetString("database.dbname"),
	)

	handler = post.Handler{DB: gorm.DB}

	e = echo.New()
	validator := validator.InitValidator()
	e.Validator = &validator
	e.Logger.SetLevel(log.INFO)

	// Middleware
	e.Use(middleware.Logger())
})

var _ = AfterSuite(func() {
	gorm.DB.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
})
