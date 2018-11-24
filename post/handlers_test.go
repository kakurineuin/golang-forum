package post_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Post Handlers", func() {
	Describe("Find topics statistics", func() {
		It("should find successfully", func() {
			req := httptest.NewRequest(http.MethodGet, "/topics/statistics", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := handler.FindTopicsStatistics(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()
			fmt.Println("topics statistics", recBody)

			Expect(recBody).To(ContainSubstring("golang"))
			Expect(recBody).To(ContainSubstring("nodeJS"))
		})
	})

	Describe("Find topics", func() {
		It("should find successfully", func() {
			req := httptest.NewRequest(http.MethodGet, "/topics/golang?offset=0&limit=10", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("category")
			c.SetParamValues("golang")
			err := handler.FindTopics(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()
			fmt.Println("Find topics recBody", recBody)

			Expect(recBody).To(ContainSubstring("topics"))
			Expect(recBody).To(ContainSubstring("totalCount"))
		})
	})

	Describe("Find topic", func() {
		It("should find successfully", func() {
			req := httptest.NewRequest(http.MethodGet, "/topics/golang/30?offset=0&limit=10", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("category", "id")
			c.SetParamValues("golang", "30")
			err := handler.FindTopic(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()
			fmt.Println("Find topic recBody", recBody)

			Expect(recBody).To(ContainSubstring("posts"))
			Expect(recBody).To(ContainSubstring("totalCount"))
		})
	})

	Describe("Create post", func() {
		It("should create successfully", func() {
			requestJSON := `{
				"userProfileID": 2,
				"replyPostID": null,
				"topic": "測試新增文章",
				"content": "測試新增文章。"
			}`
			req := httptest.NewRequest(http.MethodPost, "/topics/golang", strings.NewReader(requestJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("category")
			c.SetParamValues("golang")
			err := handler.CreatePost(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusCreated))

			recBody := rec.Body.String()
			fmt.Println("Create post recBody", recBody)

			Expect(recBody).To(ContainSubstring("post"))
		})
	})

	Describe("Update post", func() {
		It("should update successfully", func() {
			requestJSON := `{
				"content": "測試修改文章。"
			}`
			req := httptest.NewRequest(http.MethodPut, "/topics/golang/30", strings.NewReader(requestJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("category", "id")
			c.SetParamValues("golang", "30")
			c.Set("user", createToken())
			err := handler.UpdatePost(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()
			fmt.Println("Update post recBody", recBody)

			Expect(recBody).To(ContainSubstring("post"))
		})
	})
})

func createToken() *jwt.Token {
	token := jwt.New(jwt.SigningMethodHS256)
	exp := time.Now().Add(time.Hour * 72).Unix()

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = float64(2)
	claims["username"] = "admin"
	claims["email"] = "admin@xxx.com"
	claims["exp"] = exp
	return token
}
