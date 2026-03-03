package store

import (
	"context"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

type LeaderboardStore interface {
	SubmitScore(ctx context.Context, boardID, userID string, score float64) (int64, float64, error)
	GetLeaderboardPage(ctx context.Context, boardID string, offset, limit int64) ([]Entry, int64, error)
	GetUserRank(ctx context.Context, boardID, userID string) (int64, float64, error)
	RemoveUser(ctx context.Context, boardID, userID string) error
}

type Entry struct {
	UserID string
	Score  float64
}
