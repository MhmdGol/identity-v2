package sql

import (
	"context"
	"fmt"
	"identity-v2/internal/model"
	"identity-v2/internal/repository"
	"identity-v2/internal/repository/sql/sqlmodel"

	"github.com/uptrace/bun"
)

type UserRepo struct {
	db *bun.DB
}

var _ repository.UserRepo = (*UserRepo)(nil)

func NewUserRepo(db *bun.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (ur *UserRepo) Create(ctx context.Context, u model.UserInfo) error {
	var role sqlmodel.Role
	err := ur.db.NewSelect().Model(&role).Where("name = ?", u.Role).Scan(ctx)
	if err != nil {
		fmt.Println(1, err)
		return err
	}

	var status sqlmodel.Status
	err = ur.db.NewSelect().Model(&status).Where("name = ?", u.Status).Scan(ctx)
	if err != nil {
		fmt.Println(2, err)
		return err
	}

	newUser := sqlmodel.User{
		ID:             int64(u.ID),
		UUN:            u.UUN,
		Username:       u.Username,
		HashedPassword: u.HashedPassword,
		Email:          u.Email,
		Created_at:     u.Created_at,
		TOTPIsActive:   u.TOTPIsActive,
		TOTPSecret:     "",
		Role:           role.ID,
		Status:         status.ID,
	}

	_, err = ur.db.NewInsert().Model(&newUser).Exec(ctx)
	fmt.Println(3, err)

	return err
}
