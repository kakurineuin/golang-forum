package model

// ForumStatistics 論壇統計資料。
type ForumStatistics struct {
	TopicCount int `json:"topicCount"`
	ReplyCount int `json:"replyCount"`
	UserCount  int `json:"userCount"`
}
