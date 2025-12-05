package model

import (
	"time"
)

type User struct {
	ID        int64
	Username  string
	Email     string
	CreatedAt time.Time
	updatedAt time.Time
}
