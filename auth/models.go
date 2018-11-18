package auth

import (
	"time"
)

// UserProfile 使用者資料。
type UserProfile struct {
	ID        *int       `json:"id"`
	Account   *string    `json:"account" valid:"required~帳號必填。,alphanum~帳號只能英文字母或數字。,stringlength(5|20)~帳號長度須在5至20之間。"`
	Email     *string    `json:"email" valid:"required~email必填。,email~email格式不正確。,stringlength(5|30)~email長度須在5至30之間。"`
	Password  *string    `json:"password" valid:"required~密碼必填。,stringlength(5|20)~密碼長度須在5至20之間。"`
	Role      *string    `json:"role"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

// LoginRequest 登入請求。
type LoginRequest struct {
	Email    *string `json:"email" valid:"required~email必填。,email~email格式不正確。,stringlength(5|30)~email長度須在5至30之間。"`
	Password *string `json:"password" valid:"required~密碼必填。,stringlength(5|20)~密碼長度須在5至20之間。"`
}
