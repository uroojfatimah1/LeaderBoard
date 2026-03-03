package service

import (
	"context"
	"errors"
	"leaderBoard/internal/store"
)

type Service struct {
	store store.LeaderboardStore
}

func (s *Service) GetLeaderboard(ctx context.Context, boardID string, offset, limit int64) ([]store.Entry, int64, error) {
	if limit > 100 {
		limit = 100
	}

	if boardID == "" {
		return nil, 0, ErrInvalidBoard
	}

	if offset < 0 {
		offset = 0
	}

	if limit <= 0 {
		limit = 50
	}

	return s.store.GetLeaderboardPage(ctx, boardID, offset, limit)
}

func NewService(store store.LeaderboardStore) *Service {
	return &Service{store: store}
}

func (s *Service) SubmitScore(
	ctx context.Context,
	boardID, userID string,
	score float64,
) (int64, float64, error) {

	if boardID == "" {
		return 0, 0, ErrInvalidBoard
	}

	if userID == "" {
		return 0, 0, ErrInvalidUser
	}

	if score < 0 {
		return 0, 0, ErrInvalidScore
	}

	rank, updatedScore, err := s.store.SubmitScore(ctx, boardID, userID, score)
	if err != nil {
		return 0, 0, ErrStoreUnavailable
	}

	return rank, updatedScore, nil
}

func (s *Service) GetUserRank(
	ctx context.Context,
	boardID, userID string,
) (int64, float64, error) {

	if boardID == "" {
		return 0, 0, ErrInvalidBoard
	}

	if userID == "" {
		return 0, 0, ErrInvalidUser
	}

	rank, score, err := s.store.GetUserRank(ctx, boardID, userID)
	if err != nil {
		return 0, 0, ErrUserNotFound
	}

	return rank, score, nil
}

func (s *Service) RemoveUser(
	ctx context.Context,
	boardID, userID string,
) error {

	if boardID == "" {
		return ErrInvalidBoard
	}

	if userID == "" {
		return ErrInvalidUser
	}

	err := s.store.RemoveUser(ctx, boardID, userID)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return ErrStoreUnavailable
	}

	return nil
}
