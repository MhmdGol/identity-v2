package repository

import (
	"context"
	"identity-v2/internal/model"
)

type UserRepo interface {
	Create(context.Context, model.UserInfo) error
}
