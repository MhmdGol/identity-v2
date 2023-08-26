package service

import (
	"context"
	"identity-v2/internal/model"
	"identity-v2/internal/repository"
	"identity-v2/internal/service"
	rate "identity-v2/internal/service/rate_strategy"
	"time"
)

type LoginAttemptService struct {
	loginAttemptRepo repository.LoginAttemptRepo
	rateControl      rate.RateControl
}

var _ service.LoginAttemptService = (*LoginAttemptService)(nil)

func NewLoginAttempService() *LoginAttemptService {
	return &LoginAttemptService{}
}

func (las *LoginAttemptService) CheckAttempt(ctx context.Context, id model.ID) (model.AttemptValid, error) {
	a, err := las.loginAttemptRepo.ByID(ctx, id)
	if err != nil {
		return false, err
	}
	if time.Now().UTC().After(a.LastAttempt.Add(time.Hour * 3)) {
		a2 := model.LoginAttempt{
			ID:          a.ID,
			Attempts:    0,
			LastAttempt: time.Time{},
			BanExpiry:   time.Time{},
		}
		las.loginAttemptRepo.Update(ctx, a2)
		return true, nil
	}
	if time.Now().UTC().After(a.BanExpiry) {
		return true, nil
	}

	return false, nil
}

func (las *LoginAttemptService) FailedAttempt(ctx context.Context, id model.ID) error {
	a, err := las.loginAttemptRepo.ByID(ctx, id)
	if err != nil {
		return err
	}

	attempts := a.Attempts + 1

	if attempts == 3 {
		las.rateControl = rate.NewThirdFailure(las.loginAttemptRepo)
	} else if attempts == 4 {
		las.rateControl = rate.NewFourthFailure(las.loginAttemptRepo)
	} else if attempts > 4 {
		las.rateControl = rate.NewMoreThanFourFailure(las.loginAttemptRepo)
	}

	las.rateControl.Ban(ctx, id, attempts)

	return nil
}

func (las *LoginAttemptService) ResetAttempt(ctx context.Context, id model.ID) error {
	return las.loginAttemptRepo.Update(ctx, model.LoginAttempt{
		ID:          id,
		Attempts:    0,
		LastAttempt: time.Time{},
		BanExpiry:   time.Time{},
	})
}
