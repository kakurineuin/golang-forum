package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/kakurineuin/golang-forum/auth"
	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Auth Handler", func() {
	Describe("Register", func() {
		AfterEach(func() {
			dao.DB.Delete(auth.UserProfile{})
		})
		It("should register successfully", func() {
			requestJSON := `{
				"username": "test001",
				"email": "test001@xxx.com",
				"password": "test123"
			}`
			req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(requestJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := handler.Register(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()
			var result struct {
				Token            string `json:"token"`
				Exp              int    `json:"exp"`
				auth.UserProfile `json:"userProfile"`
			}
			err = json.Unmarshal([]byte(recBody), &result)

			Expect(err).To(BeNil())
			Expect(result).To(MatchAllFields(Fields{
				"Token": Not(BeEmpty()),
				"Exp":   BeNumerically(">=", 0),
				"UserProfile": MatchAllFields(Fields{
					"ID":         PointTo(BeNumerically(">=", 0)),
					"Username":   PointTo(Not(BeEmpty())),
					"Email":      PointTo(Not(BeEmpty())),
					"Password":   BeNil(), // 密碼不能傳到前端。
					"Role":       PointTo(Equal("user")),
					"IsDisabled": PointTo(BeNumerically("==", 0)),
					"CreatedAt":  Not(BeNil()),
					"UpdatedAt":  Not(BeNil()),
				}),
			}))
		})
	})

	Describe("Login", func() {
		BeforeEach(func() {
			username := "test001"
			email := "test001@xxx.com"
			password := "$2a$10$041tGlbd86T90uNSGbvkw.tSExCrlKmy37QoUGl23mfW7YGJjUVjO"
			role := "user"
			isDisabled := 0
			test001 := auth.UserProfile{
				Username:   &username,
				Email:      &email,
				Password:   &password,
				Role:       &role,
				IsDisabled: &isDisabled,
			}

			if err := dao.DB.Table("user_profile").Create(&test001).Error; err != nil {
				panic(err)
			}
		})
		AfterEach(func() {
			dao.DB.Delete(auth.UserProfile{})
		})
		It("should login successfully", func() {
			requestJSON := `{
				"email": "test001@xxx.com",
				"password": "test123"
			}`
			req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(requestJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := handler.Login(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()
			var result struct {
				Token            string `json:"token"`
				Exp              int    `json:"exp"`
				auth.UserProfile `json:"userProfile"`
			}
			err = json.Unmarshal([]byte(recBody), &result)

			Expect(err).To(BeNil())
			Expect(result).To(MatchAllFields(Fields{
				"Token": Not(BeEmpty()),
				"Exp":   BeNumerically(">=", 0),
				"UserProfile": MatchAllFields(Fields{
					"ID":         PointTo(BeNumerically(">=", 0)),
					"Username":   PointTo(Not(BeEmpty())),
					"Email":      PointTo(Not(BeEmpty())),
					"Password":   BeNil(), // 密碼不能傳到前端。
					"Role":       PointTo(Equal("user")),
					"IsDisabled": PointTo(BeNumerically("==", 0)),
					"CreatedAt":  Not(BeNil()),
					"UpdatedAt":  Not(BeNil()),
				}),
			}))
		})
	})
})
