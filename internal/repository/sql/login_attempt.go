package sql

import (
	"context"
	"identity-v2/internal/model"
	"identity-v2/internal/repository"
	"identity-v2/internal/repository/sql/sqlmodel"
	"time"

	"github.com/uptrace/bun"
)

type LoginAttemptRepo struct {
	db *bun.DB
}

var _ repository.LoginAttemptRepo = (*LoginAttemptRepo)(nil)

func NewLoginAttemptRepo(db *bun.DB) *LoginAttemptRepo {
	return &LoginAttemptRepo{
		db: db,
	}
}

func (lar *LoginAttemptRepo) CreateAttempt(ctx context.Context, id model.ID) error {
	attempt := sqlmodel.LoginAttempt{
		UserID:      int64(id),
		Attempts:    0,
		LastAttempt: time.Time{},
		BanExpiry:   time.Time{},
	}

	_, err := lar.db.NewInsert().Model(&attempt).Exec(ctx)

	return err
}

func (lar *LoginAttemptRepo) IncrAttempt(ctx context.Context, id model.ID) error {
	var attempt sqlmodel.LoginAttempt
	err := lar.db.NewSelect().Model(&attempt).Where("user_id = ?", id).Scan(ctx)
	if err != nil {
		return err
	}

	attempt := sqlmodel.LoginAttempt{
		UserID:      int64(a),
		Attempts:    0,
		LastAttempt: time.Time{},
		BanExpiry:   time.Time{},
	}
}
