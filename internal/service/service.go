package service

import (
	"context"
	"identity-v2/internal/model"
)

type UserService interface {
	Create(context.Context, model.RawUser) error
}
