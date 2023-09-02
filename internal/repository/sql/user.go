package sql

import (
	"context"
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
	// it can be deleted
	var role sqlmodel.Role
	err := ur.db.NewSelect().Model(&role).Where("name = ?", u.Role).Scan(ctx)
	if err != nil {
		return err
	}

	var status sqlmodel.Status
	err = ur.db.NewSelect().Model(&status).Where("name = ?", u.Status).Scan(ctx)
	if err != nil {
		return err
	}
	// -----------------

	newUser := sqlmodel.User{
		ID:             int64(u.ID),
		UUN:            u.UUN,
		Username:       u.Username,
		HashedPassword: u.HashedPassword,
		Email:          u.Email,
		Created_at:     u.Created_at,
		TOTPIsActive:   u.TOTPIsActive,
		TOTPSecret:     "",
		RoleID:         role.ID,
		StatusID:       status.ID,
	}

	_, err = ur.db.NewInsert().Model(&newUser).Exec(ctx)
	return err
}

func (ur *UserRepo) ByEmail(ctx context.Context, e string) (model.UserInfo, error) {
	var user sqlmodel.User
	err := ur.db.NewSelect().
		Model(&user).
		Relation("Role").
		Relation("Status").
		Where("email = ?", e).
		Scan(ctx)
	if err != nil {
		return model.UserInfo{}, err
	}

	return model.UserInfo{
		ID:             model.ID(user.ID),
		UUN:            user.UUN,
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		Email:          user.Email,
		Created_at:     user.Created_at,
		TOTPIsActive:   user.TOTPIsActive,
		TOTPSecret:     user.TOTPSecret,
		Role:           user.Role.Name,
		Status:         user.Status.Name,
	}, nil
}

func (ur *UserRepo) Exists(ctx context.Context, e string) (bool, error) {
	b, err := ur.db.NewSelect().Model((*sqlmodel.User)(nil)).Where("email = ?", e).Exists(ctx)
	return b, err
}

func (ur *UserRepo) Update(ctx context.Context, user model.UserInfo) error {
	var totpIsActiveInt int32
	if user.TOTPIsActive {
		totpIsActiveInt = 1
	} else {
		totpIsActiveInt = 0
	}

	_, err := ur.db.NewUpdate().Model((*sqlmodel.User)(nil)).
		Set("totp_is_active = ?", totpIsActiveInt).
		Set("totp_secret = ?", user.TOTPSecret).
		Where("email = ?", user.Email).
		Exec(ctx)

	return err
}
