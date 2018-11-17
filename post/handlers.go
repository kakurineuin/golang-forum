package post

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"net/http"
)

// Handler 處理請求的 handler。
type Handler struct {
	DB *gorm.DB
}

func (h Handler) FindPostsStatistics(c echo.Context) (err error) {

	// 查詢 golang 文章統計資料。
	var golangPostStatistic PostStatistic

	h.DB.Raw(`select
		(select count(*) from post_golang where reply_post_id is null) as topic_count,
		(select count(*) from post_golang where reply_post_id is not null) as reply_count,
		u.account as last_post_account,
		p.created_at as last_post_time
	from post_golang p 
		inner join user_profile u 
			on p.user_profile_id = u.id
	order by p.id desc
	limit 1`).Scan(&golangPostStatistic)

	// 查詢 Node.js 文章統計資料。
	var nodeJSPostStatistic PostStatistic

	h.DB.Raw(`select
		(select count(*) from post_nodejs where reply_post_id is null) as topic_count,
		(select count(*) from post_nodejs where reply_post_id is not null) as reply_count,
		u.account as last_post_account,
		p.created_at as last_post_time
	from post_nodejs p 
		inner join user_profile u 
			on p.user_profile_id = u.id
	order by p.id desc
	limit 1`).Scan(&nodeJSPostStatistic)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"golang": golangPostStatistic,
		"nodeJS": nodeJSPostStatistic,
	})
}

// FindPosts 查詢文章。
func (h Handler) FindPosts(c echo.Context) (err error) {
	category := c.Param("category")
	offset := c.QueryParam("offset")
	limit := c.QueryParam("limit")
	c.Logger().Infof("category: %v, offset: %v, limit: %v", category, offset, limit)

	table := ""

	switch category {
	case "golang":
		table = "post_golang"
	case "nodejs":
		table = "post_nodejs"
	default:
		return errors.New("category is error")
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
		"posts":      findPostsResults,
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

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "新增文章成功。",
		"post":    post,
	})
}
