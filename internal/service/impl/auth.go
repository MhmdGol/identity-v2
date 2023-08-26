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
	jwt *jwt.JwtToken,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		trackRepo:   trackRepo,
		jwt:         jwt,
	}
}

func (as *AuthService) Login(ctx context.Context, l model.LoginInfo) (model.JwtToken, error) {
	user, err := as.userRepo.ByEmail(ctx, l.Email)
	if err != nil {
		return "", err
	}

	s, err := as.sessionRepo.ByID(ctx, user.ID)
	if err != nil {
		return "", err
	}
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
			return "", err // fmt.Errorf("totp code not valid")
		}
	}

	session := model.Session{
		UserID:     user.ID,
		SessionExp: time.Now().Add(time.Hour),
	}

	err = as.sessionRepo.Add(ctx, session)
	if err != nil {
		return "", err
	}

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
		Timestamp: time.Now(),
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
