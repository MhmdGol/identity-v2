package rate

import (
	"context"
	"identity-v2/internal/model"
	"identity-v2/internal/repository"
	"time"
)

type ThirdFailure struct {
	loginAttemptRepo repository.LoginAttemptRepo
}

func NewThirdFailure(loginAttemptRepo repository.LoginAttemptRepo) *ThirdFailure {
	return &ThirdFailure{
		loginAttemptRepo: loginAttemptRepo,
	}
}

func (rc *ThirdFailure) Ban(ctx context.Context, id model.ID, a int32) error {
	return rc.loginAttemptRepo.Update(ctx, model.LoginAttempt{
		ID:          id,
		Attempts:    a,
		LastAttempt: time.Now().UTC(),
		BanExpiry:   time.Now().UTC().Add(time.Minute * 15),
	})
}
