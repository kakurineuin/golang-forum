package forum_test

import (
	"encoding/json"
	"github.com/kakurineuin/golang-forum/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Forum Handler", func() {
	var userProfileId int // 使用者 Id。

	BeforeEach(func() {
		// 新增一名使用者。
		username := "test001"
		email := "test001@xxx.com"
		password := "$2a$10$041tGlbd86T90uNSGbvkw.tSExCrlKmy37QoUGl23mfW7YGJjUVjO"
		role := "user"
		user1 := model.UserProfile{
			Username: &username,
			Email:    &email,
			Password: &password,
			Role:     &role,
		}

		if err := dao.DB.Create(&user1).Error; err != nil {
			panic(err)
		}

		userProfileId = *user1.Id

		// 新增文章。
		for _, table := range []string{"post_golang", "post_nodejs"} {
			topic := "測試主題001"
			content := "內容..."
			post1 := model.Post{
				UserProfileId: &userProfileId,
				Topic:         &topic,
				Content:       &content,
			}

			if err := dao.DB.Table(table).Create(&post1).Error; err != nil {
				panic(err)
			}

			reply1 := model.Post{
				UserProfileId: &userProfileId,
				ReplyPostId:   post1.Id,
				Topic:         &topic,
				Content:       &content,
			}

			if err := dao.DB.Table(table).Create(&reply1).Error; err != nil {
				panic(err)
			}
		}
	})

	AfterEach(func() {
		dao.DB.Delete(model.UserProfile{})

		for _, table := range []string{"post_golang", "post_nodejs"} {
			dao.DB.Table(table).Unscoped().Delete(model.Post{})
		}
	})

	Describe("Find forum statistics", func() {
		It("should find successfully", func() {
			req := httptest.NewRequest(http.MethodGet, "/api/forum/statistics", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := forumHandler.FindForumStatistics(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()
			var result struct {
				model.ForumStatistics `json:"forumStatistics"`
			}
			err = json.Unmarshal([]byte(recBody), &result)

			Expect(err).To(BeNil())
			Expect(result).To(MatchAllFields(Fields{
				"ForumStatistics": MatchAllFields(Fields{
					"TopicCount": BeNumerically("==", 2),
					"ReplyCount": BeNumerically("==", 2),
					"UserCount":  BeNumerically("==", 1),
				}),
			}))
		})
	})
})
