package service

import (
	"context"
	"fmt"
	"identity-v2/internal/model"
	"identity-v2/internal/repository"
	"identity-v2/internal/service"
	"identity-v2/pkg/bcrypthash"
	"identity-v2/pkg/jwt"
	"time"

	"github.com/pquerna/otp/totp"
)

type AuthService struct {
	userRepo        repository.UserRepo
	sessionRepo     repository.SessionRepo
	trackRepo       repository.TrackRepo
	loginAttemptSvc service.LoginAttemptService
	jwt             *jwt.JwtToken
}

var _ service.AuthService = (*AuthService)(nil)

func NewAuthService(
	userRepo repository.UserRepo,
	sessionRepo repository.SessionRepo,
	trackRepo repository.TrackRepo,
	loginAttemptSvc service.LoginAttemptService,
	jwt *jwt.JwtToken,
) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		sessionRepo:     sessionRepo,
		trackRepo:       trackRepo,
		loginAttemptSvc: loginAttemptSvc,
		jwt:             jwt,
	}
}

func (as *AuthService) Login(ctx context.Context, l model.LoginInfo) (model.JwtToken, error) {
	user, err := as.userRepo.ByEmail(ctx, l.Email)
	if err != nil {
		return "", err
	}

	s, _ := as.sessionRepo.ByID(ctx, user.ID)
	if time.Now().UTC().Before(s.SessionExp) {
		return "", fmt.Errorf("logout first")
	} else if time.Now().UTC().After(s.SessionExp) {
		as.sessionRepo.Remove(ctx, user.ID)
	}

	valid, err := as.loginAttemptSvc.CheckAttempt(ctx, user.ID)
	if err != nil {
		return "", err
	}
	if !valid {
		return "", fmt.Errorf("banned")
	}

	err = bcrypthash.ValidatePassword(user.HashedPassword, l.Password)
	if err != nil {
		// failed attempt
		err2 := as.loginAttemptSvc.FailedAttempt(ctx, user.ID)
		if err2 != nil {
			return "", err2
		}

		return "", err
	}

	if user.TOTPIsActive {
		isValid := totp.Validate(l.TOTPCode, user.TOTPSecret)
		if !isValid {
			// failed attempt
			err2 := as.loginAttemptSvc.FailedAttempt(ctx, user.ID)
			if err2 != nil {
				return "", err2
			}
			return "", fmt.Errorf("totp code not valid")
		}
	}

	session := model.Session{
		UserID:     user.ID,
		SessionExp: time.Now().UTC().Add(time.Hour),
	}

	err = as.sessionRepo.Add(ctx, session)
	if err != nil {
		return "", err
	}

	// from now he is logged in. then if anything below fails, user doesnt get the token thus a bad state reached

	token, err := as.jwt.MakeToken(model.TokenClaim{
		ID:    user.ID,
		Email: user.Email,
	})
	if err != nil {
		return "", err
	}

	err = as.trackRepo.Create(ctx, model.TrackInfo{
		ID:        user.ID,
		Action:    "login",
		Timestamp: time.Now().UTC(),
	})
	if err != nil {
		return "", err
	}

	err = as.loginAttemptSvc.ResetAttempt(ctx, user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (as *AuthService) Logout(ctx context.Context, id model.ID) error {
	// what if the session is already expired?

	err := as.sessionRepo.Remove(ctx, id)
	if err != nil {
		return err
	}

	err = as.trackRepo.Create(ctx, model.TrackInfo{
		ID:        id,
		Action:    "logout",
		Timestamp: time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (as *AuthService) CheckSession(ctx context.Context, id model.ID) (bool, error) {
	s, err := as.sessionRepo.ByID(ctx, id)
	if err != nil {
		return false, err
	}

	return time.Now().Before(s.SessionExp), nil
}
