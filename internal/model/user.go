package model

import "time"

type User struct {
	UserID           int       `json:"user_id"`
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	TelegramUserID   int64     `json:"telegram_user_id,omitempty"`
	TelegramUsername string    `json:"telegram_username,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

func (u User) ID() int {
	return u.UserID
}
