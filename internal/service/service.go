package service

import (
	"context"
	"identity-v2/internal/model"
)

type UserService interface {
	Create(context.Context, model.RawUser) error
	ByEmail(context.Context, string) (model.UserInfo, error)
	Exists(context.Context, string) (bool, error)
	SetTOTP(context.Context, string) (string, error)
}

type AuthService interface {
	Login(context.Context, model.LoginInfo) (model.JwtToken, error)
	Logout(context.Context, model.ID) error
	CheckSession(context.Context, model.ID) (bool, error)
}

type LoginAttemptService interface {
	CheckAttempt(context.Context, model.ID) (model.AttemptValid, error)
	FailedAttempt(context.Context, model.ID) error
	ResetAttempt(context.Context, model.ID) error
}
