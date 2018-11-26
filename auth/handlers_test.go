package auth_test

import (
	"github.com/kakurineuin/golang-forum/auth"
	"github.com/kakurineuin/golang-forum/db/gorm"
	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"strings"
)

var _ = Describe("Auth Handlers", func() {
	Describe("Register", func() {
		It("should register successfully", func() {
			requestJSON := `{
				"username": "test001",
				"email": "test001@xxx.com",
				"password": "test001"
			}`
			req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(requestJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := handler.Register(c)

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
			gorm.DB.Where("username = ?", "test001").Delete(auth.UserProfile{})
		})
		It("should login successfully", func() {
			requestJSON := `{
				"email": "test001@xxx.com",
				"password": "test001"
			}`
			req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(requestJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := handler.Login(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()

			Expect(recBody).To(ContainSubstring("token"))
			Expect(recBody).To(ContainSubstring("exp"))
			Expect(recBody).To(ContainSubstring("userProfile"))
		})
	})
})
