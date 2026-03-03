package redis

import (
	"context"
	"fmt"
	"leaderBoard/internal/store"

	"github.com/redis/go-redis/v9"
)

type Store interface {
	SubmitScore(name string, score int)
	GetLeaderboardPage() []string
	GetUserScore(name string) int
	RemoveUser(name string)
}

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

func (r *RedisStore) SubmitScore(ctx context.Context, boardID, userID string, score float64) (int64, float64, error) {
	key := fmt.Sprintf("lb:%s", boardID)

	pipe := r.client.TxPipeline()

	pipe.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: userID,
	})

	rankCmd := pipe.ZRevRank(ctx, key, userID)
	scoreCmd := pipe.ZScore(ctx, key, userID)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, 0, err
	}

	return rankCmd.Val() + 1, scoreCmd.Val(), nil
}

func (r *RedisStore) GetLeaderboardPage(
	ctx context.Context,
	boardID string,
	offset, limit int64,
) ([]store.Entry, int64, error) {

	key := fmt.Sprintf("lb:%s", boardID)

	total, err := r.client.ZCard(ctx, key).Result()
	if err != nil {
		return nil, 0, err
	}

	results, err := r.client.
		ZRevRangeWithScores(ctx, key, offset, offset+limit-1).
		Result()
	if err != nil {
		return nil, 0, err
	}

	entries := make([]store.Entry, 0, len(results))

	for _, z := range results {
		userID, ok := z.Member.(string)
		if !ok {
			continue
		}

		entries = append(entries, store.Entry{
			UserID: userID,
			Score:  z.Score,
		})
	}

	return entries, total, nil
}

func (r *RedisStore) GetUserRank(ctx context.Context, boardID, userID string) (int64, float64, error) {
	key := fmt.Sprintf("lb:%s", boardID)

	rank, err := r.client.ZRevRank(ctx, key, userID).Result()
	if err != nil {
		return 0, 0, err
	}

	score, err := r.client.ZScore(ctx, key, userID).Result()
	if err != nil {
		return 0, 0, err
	}

	return rank + 1, score, nil
}

func (r *RedisStore) RemoveUser(ctx context.Context, boardID, userID string) error {
	key := fmt.Sprintf("lb:%s", boardID)
	removed, err := r.client.ZRem(ctx, key, userID).Result()
	if err != nil {
		return err
	}
	if removed == 0 {
		return store.ErrUserNotFound
	}
	return nil
}
