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
