package post_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Post Handlers", func() {
	Describe("Find topics statistics", func() {
		It("should find suceesfully", func() {
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
		It("should find suceesfully", func() {
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
		It("should find suceesfully", func() {
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
})
