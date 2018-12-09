package post_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"

	"github.com/kakurineuin/golang-forum/auth"

	"github.com/kakurineuin/golang-forum/db/gorm"
	"github.com/kakurineuin/golang-forum/post"
	"github.com/labstack/echo"

	jwt "github.com/dgrijalva/jwt-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Post Handler", func() {
	var userProfileID int     // 使用者 ID。
	var postGolang1 post.Post // post_golang 文章。

	BeforeEach(func() {
		// 新增一名使用者。
		username := "test001"
		email := "test001@xxx.com"
		password := "$2a$10$041tGlbd86T90uNSGbvkw.tSExCrlKmy37QoUGl23mfW7YGJjUVjO"
		role := "user"
		user1 := auth.UserProfile{
			Username: &username,
			Email:    &email,
			Password: &password,
			Role:     &role,
		}

		if err := gorm.DB.Create(&user1).Error; err != nil {
			panic(err)
		}

		userProfileID = *user1.ID

		// 新增文章。
		for _, table := range []string{"post_golang", "post_nodejs"} {
			topic := "測試主題001"
			content := "內容..."
			post1 := post.Post{
				UserProfileID: &userProfileID,
				Topic:         &topic,
				Content:       &content,
			}

			if err := gorm.DB.Table(table).Create(&post1).Error; err != nil {
				panic(err)
			}

			if table == "post_golang" {
				postGolang1 = post1
			}

			reply1 := post.Post{
				UserProfileID: &userProfileID,
				ReplyPostID:   post1.ID,
				Topic:         &topic,
				Content:       &content,
			}

			if err := gorm.DB.Table(table).Create(&reply1).Error; err != nil {
				panic(err)
			}
		}
	})

	AfterEach(func() {
		gorm.DB.Delete(auth.UserProfile{})

		for _, table := range []string{"post_golang", "post_nodejs"} {
			gorm.DB.Table(table).Unscoped().Delete(post.Post{})
		}
	})

	Describe("Find forum statistics", func() {
		It("should find successfully", func() {
			req := httptest.NewRequest(http.MethodGet, "/forum/statistics", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := handler.FindForumStatistics(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()
			var result struct {
				post.ForumStatistics `json:"forumStatistics"`
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

	Describe("Find topics statistics", func() {
		It("should find successfully", func() {
			req := httptest.NewRequest(http.MethodGet, "/topics/statistics", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := handler.FindTopicsStatistics(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()
			var result struct {
				Golang post.Statistics `json:"golang"`
				NodeJS post.Statistics `json:"nodejs"`
			}
			err = json.Unmarshal([]byte(recBody), &result)

			Expect(err).To(BeNil())
			Expect(result).To(MatchAllFields(Fields{
				"Golang": MatchAllFields(Fields{
					"TopicCount":       BeNumerically("==", 1),
					"ReplyCount":       BeNumerically("==", 1),
					"LastPostUsername": PointTo(Not(BeEmpty())),
					"LastPostTime":     Not(BeNil()),
				}),
				"NodeJS": MatchAllFields(Fields{
					"TopicCount":       BeNumerically("==", 1),
					"ReplyCount":       BeNumerically("==", 1),
					"LastPostUsername": PointTo(Not(BeEmpty())),
					"LastPostTime":     Not(BeNil()),
				}),
			}))
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
			var result struct {
				Topics     []post.Topic `json:"topics"`
				TotalCount int          `json:"totalCount"`
			}
			err = json.Unmarshal([]byte(recBody), &result)

			Expect(err).To(BeNil())
			Expect(result).To(MatchAllFields(Fields{
				"Topics":     Not(BeEmpty()),
				"TotalCount": BeNumerically("==", 1),
			}))
		})
	})

	Describe("Find topic", func() {
		It("should find successfully", func() {
			id := strconv.Itoa(*postGolang1.ID)
			req := httptest.NewRequest(http.MethodGet, "/topics/golang/"+id+"?offset=0&limit=10", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("category", "id")
			c.SetParamValues("golang", id)
			err := handler.FindTopic(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()
			var result struct {
				Posts      []post.Post `json:"posts"`
				TotalCount int         `json:"totalCount"`
			}
			err = json.Unmarshal([]byte(recBody), &result)

			Expect(err).To(BeNil())
			Expect(result).To(MatchAllFields(Fields{
				"Posts":      Not(BeEmpty()),
				"TotalCount": BeNumerically("==", 2), // 主題加上回覆共 2 篇。
			}))
		})
	})

	Describe("Create post", func() {
		It("should create successfully", func() {
			requestJSON := `{
				"userProfileID": 1,
				"replyPostID": null,
				"topic": "測試新增文章",
				"content": "測試新增文章"
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
			var result struct {
				Post post.Post `json:"post"`
			}
			err = json.Unmarshal([]byte(recBody), &result)

			Expect(err).To(BeNil())
			Expect(result).To(MatchAllFields(Fields{
				"Post": MatchAllFields(Fields{
					"ID":            PointTo(BeNumerically(">=", 0)),
					"UserProfileID": PointTo(BeNumerically("==", 1)),
					"ReplyPostID":   BeNil(),
					"Topic":         PointTo(Equal("測試新增文章")),
					"Content":       PointTo(Equal("測試新增文章")),
					"CreatedAt":     Not(BeNil()),
					"UpdatedAt":     Not(BeNil()),
					"DeletedAt":     BeNil(),
				}),
			}))
		})
	})

	Describe("Update post", func() {
		It("should update successfully", func() {
			id := strconv.Itoa(*postGolang1.ID)
			requestJSON := `{
				"content": "測試修改文章"
			}`
			req := httptest.NewRequest(http.MethodPut, "/topics/golang/"+id, strings.NewReader(requestJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("category", "id")
			c.SetParamValues("golang", id)
			c.Set("user", createToken(userProfileID))
			err := handler.UpdatePost(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()
			var result struct {
				Post post.Post `json:"post"`
			}
			err = json.Unmarshal([]byte(recBody), &result)

			Expect(err).To(BeNil())
			Expect(result).To(MatchAllFields(Fields{
				"Post": MatchAllFields(Fields{
					"ID":            PointTo(BeNumerically("==", *postGolang1.ID)),
					"UserProfileID": PointTo(BeNumerically("==", *postGolang1.UserProfileID)),
					"ReplyPostID":   BeNil(),
					"Topic":         PointTo(Equal(*postGolang1.Topic)),
					"Content":       PointTo(Equal("測試修改文章")),
					"CreatedAt":     Not(BeNil()),
					"UpdatedAt":     Not(BeNil()),
					"DeletedAt":     BeNil(),
				}),
			}))
		})
	})

	Describe("Delete post", func() {
		It("should delete successfully", func() {
			id := strconv.Itoa(*postGolang1.ID)
			req := httptest.NewRequest(http.MethodDelete, "/topics/golang/"+id, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("category", "id")
			c.SetParamValues("golang", id)
			c.Set("user", createToken(userProfileID))
			err := handler.DeletePost(c)

			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			recBody := rec.Body.String()
			var result struct {
				Post post.Post `json:"post"`
			}
			err = json.Unmarshal([]byte(recBody), &result)

			Expect(err).To(BeNil())
			Expect(result).To(MatchAllFields(Fields{
				"Post": MatchAllFields(Fields{
					"ID":            PointTo(BeNumerically("==", *postGolang1.ID)),
					"UserProfileID": PointTo(BeNumerically("==", *postGolang1.UserProfileID)),
					"ReplyPostID":   BeNil(),
					"Topic":         PointTo(Equal(*postGolang1.Topic)),
					"Content":       PointTo(Equal("此篇文章已被刪除。")),
					"CreatedAt":     Not(BeNil()),
					"UpdatedAt":     Not(BeNil()),
					"DeletedAt":     Not(BeNil()),
				}),
			}))
		})
	})
})

func createToken(userProfileID int) *jwt.Token {
	token := jwt.New(jwt.SigningMethodHS256)
	exp := time.Now().Add(time.Hour * 72).Unix()

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = float64(userProfileID)
	claims["username"] = "admin"
	claims["email"] = "admin@xxx.com"
	claims["exp"] = exp
	return token
}
