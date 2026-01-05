package model

type User struct {
	UserID   int
	Username string
}

func (u User) ID() int {
	return u.UserID
}
