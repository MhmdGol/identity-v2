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

type SessionCache struct {
	sessionRepo repository.SessionRepo
	redis       *redis.Client
}

var _ repository.SessionRepo = (*SessionCache)(nil)

func NewSessionCache(
	sessionRepo repository.SessionRepo,
	redis *redis.Client,
) *SessionCache {
	return &SessionCache{
		sessionRepo: sessionRepo,
		redis:       redis,
	}
}

func (sc *SessionCache) Add(ctx context.Context, s model.Session) error {
	sessionVal, err := json.Marshal(s)
	if err != nil {
		return err
	}

	sessionKey := strconv.FormatInt(int64(s.UserID), 10)

	err = sc.redis.Set(ctx, sessionKey, sessionVal, time.Hour).Err()
	if err != nil {
		return err
	}

	return sc.sessionRepo.Add(ctx, s)
}

func (sc *SessionCache) Remove(ctx context.Context, id model.ID) error {
	sessionKey := strconv.FormatInt(int64(id), 10)

	err := sc.redis.Del(ctx, sessionKey).Err()
	if err != nil {
		return err
	}

	return sc.sessionRepo.Remove(ctx, id)

}

func (sc *SessionCache) ByID(ctx context.Context, id model.ID) (model.Session, error) {
	sessionKey := strconv.FormatInt(int64(id), 10)

	sessionVal, err := sc.redis.Get(ctx, sessionKey).Result()
	if err == redis.Nil {
		session, err2 := sc.sessionRepo.ByID(ctx, id)
		if err2 != nil {
			return model.Session{}, err2
		}

		sessionVal2, err2 := json.Marshal(session)
		if err2 != nil {
			return model.Session{}, err2
		}

		err2 = sc.redis.Set(ctx, sessionKey, sessionVal2, time.Hour).Err()
		if err2 != nil {
			return model.Session{}, err2
		}

		return session, nil
	}

	var session model.Session
	err = json.Unmarshal([]byte(sessionVal), &session)
	if err != nil {
		return model.Session{}, err
	}

	return session, nil
}
