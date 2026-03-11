package model

type User struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

func (u User) ID() int {
	return u.UserID
}
