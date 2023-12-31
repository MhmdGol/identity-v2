package cache

import (
	"context"
	"encoding/json"
	"identity-v2/internal/model"
	"identity-v2/internal/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	userRepo repository.UserRepo
	redis    *redis.Client
}

var _ repository.UserRepo = (*UserCache)(nil)

func NewUserCache(
	userRepo repository.UserRepo,
	redis *redis.Client,
) *UserCache {
	return &UserCache{
		userRepo: userRepo,
		redis:    redis,
	}
}

func (uc *UserCache) Create(ctx context.Context, u model.UserInfo) error {
	// cache by email, because we read by email

	userVal, err := json.Marshal(u)
	if err != nil {
		return err
	}

	userKey := u.Email

	err = uc.redis.Set(ctx, userKey, userVal, time.Hour).Err()
	if err != nil {
		return err
	}

	return uc.userRepo.Create(ctx, u)
}

func (uc *UserCache) ByEmail(ctx context.Context, e string) (model.UserInfo, error) {
	userVal, err := uc.redis.Get(ctx, e).Result()
	if err == redis.Nil {
		user, err2 := uc.userRepo.ByEmail(ctx, e)
		if err2 != nil {
			return model.UserInfo{}, err2
		}

		userVal2, err2 := json.Marshal(user)
		if err2 != nil {
			return model.UserInfo{}, err2
		}

		userKey2 := e

		err2 = uc.redis.Set(ctx, userKey2, userVal2, time.Hour).Err()
		if err2 != nil {
			return model.UserInfo{}, err2
		}

		return user, nil
	}

	var user model.UserInfo
	err = json.Unmarshal([]byte(userVal), &user)
	if err != nil {
		return model.UserInfo{}, err
	}

	return user, nil
}

func (uc *UserCache) Exists(ctx context.Context, e string) (bool, error) {
	// what if it is in cache and not in db? it returns true!

	_, err := uc.redis.Get(ctx, e).Result()
	if err == redis.Nil {
		b, err2 := uc.userRepo.Exists(ctx, e)
		if err2 != nil {
			return false, err2
		}

		return b, nil
	}

	return true, nil
}

func (uc *UserCache) Update(ctx context.Context, user model.UserInfo) error {
	err := uc.userRepo.Update(ctx, user)
	if err != nil {
		return err
	}

	userVal, err := json.Marshal(user)
	if err != nil {
		return err
	}

	userKey := user.Email

	err = uc.redis.Set(ctx, userKey, userVal, time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}
