package admin_test

import (
	"encoding/json"
	"github.com/kakurineuin/golang-forum/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"net/http"
	"net/http/httptest"
	"strconv"
)

var _ = Describe("Admin Handler", func() {
	userID := ""

	BeforeEach(func() {

		// 新增 5 名使用者。
		for i := 0; i < 5; i++ {
			number := strconv.Itoa(i + 1)
			username := "test00" + number
			email := "test00" + number + "@xxx.com"
			password := "$2a$10$041tGlbd86T90uNSGbvkw.tSExCrlKmy37QoUGl23mfW7YGJjUVjO"
			role := "user"
			newUser := model.UserProfile{
				Username: &username,
				Email:    &email,
				Password: &password,
				Role:     &role,
			}

			if err := dao.DB.Create(&newUser).Error; err != nil {
				panic(err)
			}

			if i == 4 {
				userID = strconv.Itoa(*newUser.ID)
			}
		}
	})

	AfterEach(func() {
		dao.DB.Delete(model.UserProfile{})
	})

	Describe("Find users", func() {
		It("should find successfully", func() {
			req := httptest.NewRequest(http.MethodGet, "/api/admin/users?offset=0&limit=10", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := adminHandler.FindUsers(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()
			var result struct {
				Users      []model.User `json:"users"`
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

	Describe("Disable users", func() {
		It("should disable user successfully", func() {
			req := httptest.NewRequest(http.MethodPost, "/api/admin/users/disable/"+userID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(userID)
			err := adminHandler.DisableUser(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()
			var result struct {
				model.User `json:"user"`
			}
			err = json.Unmarshal([]byte(recBody), &result)

			Expect(err).To(BeNil())
			Expect(result).To(MatchAllFields(Fields{
				"User": MatchAllFields(Fields{
					"ID":         Not(BeNil()),
					"Username":   Not(BeNil()),
					"Email":      Not(BeNil()),
					"Role":       Not(BeNil()),
					"IsDisabled": PointTo(BeNumerically("==", 1)),
					"CreatedAt":  Not(BeNil()),
				}),
			}))
		})
	})
})
