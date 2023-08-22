package service

import (
	"context"
	"identity-v2/internal/model"
	"identity-v2/internal/repository"
	"identity-v2/internal/service"
	"identity-v2/pkg/bcrypthash"
	"time"

	"github.com/bwmarrin/snowflake"
)

type UserService struct {
	userRepo repository.UserRepo
	snow     *snowflake.Node
}

var _ service.UserService = (*UserService)(nil)

func NewUserService(
	userRepo repository.UserRepo,
	snow *snowflake.Node,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		snow:     snow,
	}
}

func (us *UserService) Create(ctx context.Context, u model.RawUser) error {
	hpass, err := bcrypthash.HashPassword(u.Password)
	if err != nil {
		return err
	}

	return us.userRepo.Create(ctx, model.UserInfo{
		ID:             model.ID(us.snow.Generate().Int64()),
		UUN:            u.UUN,
		Username:       u.Username,
		HashedPassword: hpass,
		Email:          u.Email,
		Created_at:     time.Now(),
		TOTPIsActive:   false,
		Role:           u.Role,
		Status:         u.Status,
	})
}
