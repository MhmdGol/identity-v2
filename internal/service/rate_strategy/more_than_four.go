package rate

import (
	"context"
	"identity-v2/internal/model"
	"identity-v2/internal/repository"
	"time"
)

type MoreThanFourFailure struct {
	loginAttemptRepo repository.LoginAttemptRepo
}

func NewMoreThanFourFailure(loginAttemptRepo repository.LoginAttemptRepo) *MoreThanFourFailure {
	return &MoreThanFourFailure{
		loginAttemptRepo: loginAttemptRepo,
	}
}

func (rc *MoreThanFourFailure) Ban(ctx context.Context, id model.ID, a int32) error {
	return rc.loginAttemptRepo.Update(ctx, model.LoginAttempt{
		ID:          id,
		Attempts:    a,
		LastAttempt: time.Now().UTC(),
		BanExpiry:   time.Now().UTC().Add(time.Hour),
	})
}
