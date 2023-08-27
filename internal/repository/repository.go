package repository

import (
	"context"
	"identity-v2/internal/model"
)

type UserRepo interface {
	Create(context.Context, model.UserInfo) error
	ByEmail(context.Context, string) (model.UserInfo, error)
	Exists(context.Context, string) (bool, error)
	Update(context.Context, model.UserInfo) error
}

type SessionRepo interface {
	Add(context.Context, model.Session) error
	Remove(context.Context, model.ID) error
	ByID(context.Context, model.ID) (model.Session, error)
}

type TrackRepo interface {
	Create(context.Context, model.TrackInfo) error
}

type LoginAttemptRepo interface {
	ByID(context.Context, model.ID) (model.LoginAttempt, error)
	Create(context.Context, model.ID) error
	Update(context.Context, model.LoginAttempt) error
}
