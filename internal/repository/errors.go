package repository

import "errors"

var (
	ErrNotFound  = errors.New("url not found")
	ErrCollision = errors.New("code collision")
)
