package post_test

import (
	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"strings"
)

var _ = Describe("Post Handlers", func() {
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
			err := handler.FindPosts(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()

			Expect(recBody).To(ContainSubstring("posts"))
		})
	})
})
