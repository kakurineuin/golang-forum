package post_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Post Handlers", func() {
	Describe("Find posts", func() {
		It("should find suceesfully", func() {
			req := httptest.NewRequest(http.MethodGet, "/?offset=0&limit=10", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/posts/:category")
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
