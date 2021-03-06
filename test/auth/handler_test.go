package auth_test

import (
	"encoding/json"
	"github.com/kakurineuin/golang-forum/model"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Auth Handler", func() {
	BeforeEach(func() {
		username := "test001"
		email := "test001@xxx.com"
		password := "$2a$10$041tGlbd86T90uNSGbvkw.tSExCrlKmy37QoUGl23mfW7YGJjUVjO"
		role := "user"
		isDisabled := 0
		test001 := model.UserProfile{
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
		dao.DB.Delete(model.UserProfile{})
	})

	Describe("Register", func() {
		It("should fail to register if some parameters are missing", func() {
			requestJSON := "{}"
			req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(requestJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := authHandler.Register(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusBadRequest))
		})

		It("should register successfully", func() {
			requestJSON := `{
				"username": "test002",
				"email": "test002@xxx.com",
				"password": "test123"
			}`
			req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(requestJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := authHandler.Register(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()
			var result struct {
				Token             string `json:"token"`
				Exp               int    `json:"exp"`
				model.UserProfile `json:"userProfile"`
			}
			err = json.Unmarshal([]byte(recBody), &result)

			Expect(err).To(BeNil())
			Expect(result).To(MatchAllFields(Fields{
				"Token": Not(BeEmpty()),
				"Exp":   BeNumerically(">=", 0),
				"UserProfile": MatchAllFields(Fields{
					"Id":         PointTo(BeNumerically(">=", 0)),
					"Username":   PointTo(Equal("test002")),
					"Email":      PointTo(Equal("test002@xxx.com")),
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
		It("should login successfully", func() {
			requestJSON := `{
				"email": "test001@xxx.com",
				"password": "test123"
			}`
			req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(requestJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := authHandler.Login(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()
			var result struct {
				Token             string `json:"token"`
				Exp               int    `json:"exp"`
				model.UserProfile `json:"userProfile"`
			}
			err = json.Unmarshal([]byte(recBody), &result)

			Expect(err).To(BeNil())
			Expect(result).To(MatchAllFields(Fields{
				"Token": Not(BeEmpty()),
				"Exp":   BeNumerically(">=", 0),
				"UserProfile": MatchAllFields(Fields{
					"Id":         PointTo(BeNumerically(">=", 0)),
					"Username":   PointTo(Equal("test001")),
					"Email":      PointTo(Equal("test001@xxx.com")),
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
