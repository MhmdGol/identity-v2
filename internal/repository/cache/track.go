package cache

import (
	"context"
	"encoding/json"
	"identity-v2/internal/model"
	"identity-v2/internal/repository"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type TrackCache struct {
	trackRepo repository.TrackRepo
	redis     *redis.Client
}

var _ repository.TrackRepo = (*TrackCache)(nil)

func NewTrackCache(
	trackRepo repository.TrackRepo,
	redis *redis.Client,
) *TrackCache {
	return &TrackCache{
		trackRepo: trackRepo,
		redis:     redis,
	}
}

func (tc *TrackCache) Create(ctx context.Context, t model.TrackInfo) error {
	trackVal, err := json.Marshal(t)
	if err != nil {
		return err
	}

	trackKey := strconv.FormatInt(int64(t.ID), 10)

	err = tc.redis.Set(ctx, trackKey, trackVal, time.Minute*5).Err()
	if err != nil {
		return err
	}

	return tc.trackRepo.Create(ctx, t)
}
