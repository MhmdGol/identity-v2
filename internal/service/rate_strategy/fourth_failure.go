package rate

import (
	"context"
	"identity-v2/internal/model"
	"identity-v2/internal/repository"
	"time"
)

type FourthFailure struct {
	loginAttemptRepo repository.LoginAttemptRepo
}

func NewFourthFailure(loginAttemptRepo repository.LoginAttemptRepo) *FourthFailure {
	return &FourthFailure{
		loginAttemptRepo: loginAttemptRepo,
	}
}

func (rc *FourthFailure) Ban(ctx context.Context, id model.ID, a int32) error {
	return rc.loginAttemptRepo.Update(ctx, model.LoginAttempt{
		ID:          id,
		Attempts:    a,
		LastAttempt: time.Now().UTC(),
		BanExpiry:   time.Now().UTC().Add(time.Minute * 30),
	})
}
