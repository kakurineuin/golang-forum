package admin_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/kakurineuin/golang-forum/admin"
	"github.com/kakurineuin/golang-forum/auth"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Admin Handler", func() {
	BeforeEach(func() {

		// 新增 5 名使用者。
		for i := 0; i < 5; i++ {
			index := strconv.Itoa(i)
			username := "test00" + index
			email := "test00" + index + "@xxx.com"
			password := "$2a$10$041tGlbd86T90uNSGbvkw.tSExCrlKmy37QoUGl23mfW7YGJjUVjO"
			role := "user"
			newUser := auth.UserProfile{
				Username: &username,
				Email:    &email,
				Password: &password,
				Role:     &role,
			}

			if err := dao.DB.Create(&newUser).Error; err != nil {
				panic(err)
			}
		}
	})

	AfterEach(func() {
		dao.DB.Delete(auth.UserProfile{})
	})

	Describe("Find users", func() {
		It("should find successfully", func() {
			req := httptest.NewRequest(http.MethodGet, "/users?offset=0&limit=10", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := handler.FindUsers(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()
			var result struct {
				Users      []admin.User `json:"users"`
				TotalCount int          `json:"totalCount"`
			}
			err = json.Unmarshal([]byte(recBody), &result)

			Expect(err).To(BeNil())
			Expect(result).To(MatchAllFields(Fields{
				"Users":      Not(BeEmpty()),
				"TotalCount": BeNumerically("==", 5),
			}))
		})
	})
})
