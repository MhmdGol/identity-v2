package service

import (
	"context"
	"identity-v2/internal/model"
)

type UserService interface {
	Create(context.Context, model.RawUser) error
}

type AuthService interface {
	Login(context.Context, model.LoginInfo) (model.JwtToken, error)
}

type LoginAttemptService interface {
	CheckAttempt(context.Context, model.ID) (model.AttemptValid, error)
	FailedAttempt(context.Context, model.ID) error
	ResetAttempt(context.Context, model.ID) error
}
