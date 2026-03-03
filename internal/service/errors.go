package service

import "errors"

var (
	ErrInvalidBoard     = errors.New("invalid board id")
	ErrInvalidUser      = errors.New("invalid user id")
	ErrInvalidScore     = errors.New("invalid score")
	ErrUserNotFound     = errors.New("user not found")
	ErrStoreUnavailable = errors.New("store unavailable")
)
