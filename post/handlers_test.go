package post_test

import (
	"fmt"
	"forum/auth"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

	Describe("Update post", func() {
		It("should update successfully", func() {
			requestJSON := `{
				"content": "測試修改文章。"
			}`
			req := httptest.NewRequest(http.MethodPut, "/topics/golang/30", strings.NewReader(requestJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(echo.HeaderAuthorization, `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQHh4eC5jb20iLCJleHAiOjE1NDMyMDgwODMsImlkIjoyLCJ1c2VybmFtZSI6ImFkbWluIn0.TumjZ-tJHoVBTyVpW5nxTTi3fZOuV1yhFaL1aFF846M`)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("category", "id")
			c.SetParamValues("golang", "30")

			jwtMiddleware := middleware.JWT([]byte(auth.JwtSecret))
			err := jwtMiddleware(handler.UpdatePost)(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()
			fmt.Println("Update post recBody", recBody)

			Expect(recBody).To(ContainSubstring("post"))
		})
	})
})
