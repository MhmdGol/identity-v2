package repository

import "context"

type UserRepo interface {
	Create(context.Context)
}
