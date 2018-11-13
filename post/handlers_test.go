package post_test

import (
	"forum/config"
	"forum/db/gorm"
	"forum/post"
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

var _ = Describe("Post Handlers", func() {
	config.Init("../config")

	gorm.InitDB(
		config.Viper.GetString("database.user"),
		config.Viper.GetString("database.password"),
		config.Viper.GetString("database.dbname"),
	)

	postHandler := post.Handler{DB: gorm.DB}

	e := echo.New()
	validator := validator.InitValidator()
	e.Validator = &validator
	e.Logger.SetLevel(log.INFO)

	// Middleware
	e.Use(middleware.Logger())

	Describe("Find posts", func() {
		It("should find suceesfully", func() {
			requestJSON := `{
				"category": "golang",
				"offset": 0,
				"limit": 10
			}`
			req := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(requestJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := postHandler.FindPosts(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()

			Expect(recBody).To(ContainSubstring("posts"))
		})
	})
})
