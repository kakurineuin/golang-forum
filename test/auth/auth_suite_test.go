package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/kakurineuin/golang-forum/auth"
	"github.com/kakurineuin/golang-forum/config"
	"github.com/kakurineuin/golang-forum/db/gorm"
	"github.com/kakurineuin/golang-forum/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var e *echo.Echo
var handler auth.Handler

func TestAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Suite")
}

var _ = BeforeSuite(func() {
	config.Init("../../config", "config_test")

	gorm.InitDB(
		config.Viper.GetString("database.user"),
		config.Viper.GetString("database.password"),
		config.Viper.GetString("database.dbname"),
	)

	service := auth.Service{DB: gorm.DB}
	handler = auth.Handler{Service: &service}

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