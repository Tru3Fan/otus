package service

import "errors"

var (
	ErrEmptyUsername = errors.New("username cannot be empty")
	ErrEmptyTitle    = errors.New("title cannot be empty")
)
