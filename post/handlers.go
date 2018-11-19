package post

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

// Handler 處理請求的 handler。
type Handler struct {
	DB *gorm.DB
}

// FindTopicsStatistics 查詢主題統計資料。
func (h Handler) FindTopicsStatistics(c echo.Context) (err error) {

	// 查詢 golang 文章統計資料。
	var golangStatistics Statistics

	h.DB.Raw(`select
		(select count(*) from post_golang where reply_post_id is null) as topic_count,
		(select count(*) from post_golang where reply_post_id is not null) as reply_count,
		u.account as last_post_account,
		p.created_at as last_post_time
	from post_golang p
		inner join user_profile u
			on p.user_profile_id = u.id
	order by p.id desc
	limit 1`).Scan(&golangStatistics)

	// 查詢 Node.js 文章統計資料。
	var nodeJSStatistics Statistics

	h.DB.Raw(`select
		(select count(*) from post_nodejs where reply_post_id is null) as topic_count,
		(select count(*) from post_nodejs where reply_post_id is not null) as reply_count,
		u.account as last_post_account,
		p.created_at as last_post_time
	from post_nodejs p
		inner join user_profile u
			on p.user_profile_id = u.id
	order by p.id desc
	limit 1`).Scan(&nodeJSStatistics)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"golang": golangStatistics,
		"nodeJS": nodeJSStatistics,
	})
}

// FindTopics 查詢主題列表。
func (h Handler) FindTopics(c echo.Context) (err error) {
	category := c.Param("category")
	offset := c.QueryParam("offset")
	limit := c.QueryParam("limit")
	c.Logger().Infof("category: %v, offset: %v, limit: %v", category, offset, limit)

	table, err := getTable(category)
	if err != nil {
		return
	}

	sql := fmt.Sprintf(`select
		p.id,
		p.topic,
		IFNULL(last_reply.reply_count, 0) as reply_count,
		p.created_at,
		u.account,
		last_reply.created_at as last_reply_created_at,
		last_reply.account as last_reply_account
	from
		%v p
		inner join user_profile u
			on p.user_profile_id = u.id
		left join view_post_golang_each_topic_last_reply last_reply
			on p.id = last_reply.reply_post_id
	where p.reply_post_id is null
	order by p.id desc
	limit ?, ?`, table)
	rows, err := h.DB.Raw(sql, offset, limit).Rows()
	defer rows.Close()

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return
	}

	topics := make([]Topic, 0)

	for rows.Next() {
		var topic Topic
		h.DB.ScanRows(rows, &topic)
		topics = append(topics, topic)
	}

	// 查詢總筆數。
	totalCount := 0
	sql = fmt.Sprintf(`select
		count(*)
	from
		%v p
		inner join user_profile u
			on p.user_profile_id = u.id
		left join view_post_golang_each_topic_last_reply last_reply
			on p.id = last_reply.reply_post_id
	where p.reply_post_id is null`, table)
	row := h.DB.Raw(sql).Row()
	row.Scan(&totalCount)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"topics":     topics,
		"totalCount": totalCount,
	})
}

// CreatePost 新增文章。
func (h Handler) CreatePost(c echo.Context) (err error) {
	post := new(Post)

	if err = c.Bind(post); err != nil {
		return
	}

	if err = c.Validate(post); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	if err = h.DB.Table("post_" + c.Param("category")).Create(post).Error; err != nil {
		return
	}

	message := ""

	if post.ReplyPostID == nil {
		message = "新增主題成功。"
	} else {
		message = "回覆成功。"
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": message,
		"post":    post,
	})
}

// FindTopic 查詢某個主題的討論文章。
func (h Handler) FindTopic(c echo.Context) (err error) {
	category := c.Param("category")
	id := c.Param("id")
	offset := c.QueryParam("offset")
	limit := c.QueryParam("limit")
	c.Logger().Infof("category: %v, id: %v, offset: %v, limit: %v", category, id, offset, limit)

	table, err := getTable(category)

	if err != nil {
		return
	}

	sql := fmt.Sprintf(`select
		p.id, p.topic, p.content, p.created_at, p.updated_at, u.account, u.role
	from %v p
		inner join user_profile u
			on p.user_profile_id = u.id
	where p.id = ? and p.reply_post_id is null
	union all
	select
		p.id, p.topic, p.content, p.created_at, p.updated_at, u.account, u.role
	from %v p
		inner join user_profile u
			on p.user_profile_id = u.id
	where p.reply_post_id = ?
	order by id
	limit ?, ?`, table, table)
	rows, err := h.DB.Raw(sql, id, id, offset, limit).Rows()
	defer rows.Close()

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return
	}

	findPostsResults := make([]FindPostsResult, 0)

	for rows.Next() {
		var findPostsResult FindPostsResult
		h.DB.ScanRows(rows, &findPostsResult)
		findPostsResults = append(findPostsResults, findPostsResult)
	}

	// 查詢總筆數。
	totalCount := 0
	sql = fmt.Sprintf(`select
		count(*)
	from (
		select
			p.id
		from %v p
			inner join user_profile u
				on p.user_profile_id = u.id
		where p.id = ? and p.reply_post_id is null
		union all
		select
			p.id
		from %v p
			inner join user_profile u
				on p.user_profile_id = u.id
		where p.reply_post_id = ?) t
		`, table, table)
	row := h.DB.Raw(sql, id, id).Row()
	row.Scan(&totalCount)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"posts":      findPostsResults,
		"totalCount": totalCount,
	})
}

func getTable(category string) (string, error) {
	switch category {
	case "golang":
		return "post_golang", nil
	case "nodejs":
		return "post_nodejs", nil
	default:
		return "", errors.New("category is error")
	}
}
