package admin

import (
	"time"
)

// User 使用者。
type User struct {
	ID        *int       `json:"id"`
	Username  *string    `json:"username"`
	Email     *string    `json:"email"`
	Role      *string    `json:"role"`
	CreatedAt *time.Time `json:"createdAt"`
}