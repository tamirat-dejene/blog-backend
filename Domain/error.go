package domain

import "errors"

type DomainError struct {
	Err  error
	Code int
}

var (
	ErrNotFound     = errors.New("not found")
	ErrUserNotFound = errors.New("user not found")
	ErrTokenExpired = errors.New("token expired")
	ErrInvalidInput = errors.New("invalid input")
	ErrUnauthorized = errors.New("User not authenticated or authorized")
	ErrInvalidFile  = errors.New("invalid file format")
)
