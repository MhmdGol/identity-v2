package service

import (
	"context"
	"identity-v2/internal/model"
	"identity-v2/internal/repository"
	"identity-v2/internal/service"
	"identity-v2/pkg/bcrypthash"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/casbin/casbin/v2"
)

type UserService struct {
	userRepo repository.UserRepo
	sf       *snowflake.Node
	e        *casbin.Enforcer
}

var _ service.UserService = (*UserService)(nil)

func NewUserService(
	userRepo repository.UserRepo,
	sf *snowflake.Node,
	e *casbin.Enforcer,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		sf:       sf,
		e:        e,
	}
}

func (us *UserService) Create(ctx context.Context, u model.RawUser) error {
	hpass, err := bcrypthash.HashPassword(u.Password)
	if err != nil {
		return err
	}

	err = us.userRepo.Create(ctx, model.UserInfo{
		ID:             model.ID(us.sf.Generate().Int64()),
		UUN:            u.UUN,
		Username:       u.Username,
		HashedPassword: hpass,
		Email:          u.Email,
		Created_at:     time.Now(),
		TOTPIsActive:   false,
		Role:           u.Role,
		Status:         u.Status,
	})
	if err != nil {
		return err
	}

	us.e.LoadPolicy()
	us.e.AddGroupingPolicy(u.Email, u.Role)
	us.e.SavePolicy()

	return nil
}

func (us *UserService) ByEmail(ctx context.Context, e string) (model.UserInfo, error) {
	return us.userRepo.ByEmail(ctx, e)
}
