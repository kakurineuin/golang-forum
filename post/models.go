package post

import (
	"time"
)

// Post 文章。
type Post struct {
	ID            *int       `json:"id"`
	UserProfileID *int       `json:"userProfileID"`
	ReplyPostID   *int       `json:"replyPostID"`
	Title         *string    `json:"title" valid:"required~標題必填。,stringlength(1|30)~標題長度須在1至30之間。"`
	Content       *string    `json:"content" valid:"required~內文必填。,stringlength(1|500)~內文長度須在1至500之間。"`
	CreatedAt     *time.Time `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt"`
}
