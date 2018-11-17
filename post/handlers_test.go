package post_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Post Handlers", func() {
	Describe("Find posts statistics", func() {
		It("should find suceesfully", func() {
			req := httptest.NewRequest(http.MethodGet, "/posts/statistics", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := handler.FindPostsStatistics(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()
			fmt.Println("posts statistics", recBody)

			Expect(recBody).To(ContainSubstring("golang"))
			Expect(recBody).To(ContainSubstring("nodeJS"))
		})
	})

	Describe("Find posts", func() {
		It("should find suceesfully", func() {
			req := httptest.NewRequest(http.MethodGet, "/posts/golang?offset=0&limit=10", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("category")
			c.SetParamValues("golang")
			err := handler.FindPosts(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()

			Expect(recBody).To(ContainSubstring("posts"))
		})
	})
})
