package model

import (
	"time"
)

// User 使用者。
type User struct {
	Id         *int       `json:"id"`
	Username   *string    `json:"username"`
	Email      *string    `json:"email"`
	Role       *string    `json:"role"`
	IsDisabled *int       `json:"isDisabled"`
	CreatedAt  *time.Time `json:"createdAt"`
}
