package auth_test

import (
	"forum/auth"
	"forum/config"
	"forum/db/gorm"
	"forum/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"strings"
)

var _ = Describe("Auth Handlers", func() {
	config.Init("../config")

	gorm.InitDB(
		config.Viper.GetString("database.user"),
		config.Viper.GetString("database.password"),
		config.Viper.GetString("database.dbname"),
	)

	authHandler := auth.Handler{DB: gorm.DB}

	e := echo.New()
	validator := validator.InitValidator()
	e.Validator = &validator
	e.Logger.SetLevel(log.INFO)

	// Middleware
	e.Use(middleware.Logger())

	Describe("Register", func() {
		It("should register suceesfully", func() {
			requestJSON := `{
				"account": "test001",
				"email": "test001@xxx.com",
				"password": "test001"
			}`
			req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(requestJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := authHandler.Register(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()

			Expect(recBody).To(ContainSubstring("token"))
			Expect(recBody).To(ContainSubstring("exp"))
			Expect(recBody).To(ContainSubstring("userProfile"))
		})
	})

	Describe("Login", func() {
		AfterEach(func() {
			gorm.DB.Where("account = ?", "test001").Delete(auth.UserProfile{})
		})
		It("should login suceesfully", func() {
			requestJSON := `{
				"email": "test001@xxx.com",
				"password": "test001"
			}`
			req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(requestJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := authHandler.Login(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()

			Expect(recBody).To(ContainSubstring("token"))
			Expect(recBody).To(ContainSubstring("exp"))
			Expect(recBody).To(ContainSubstring("userProfile"))
		})
	})

})
