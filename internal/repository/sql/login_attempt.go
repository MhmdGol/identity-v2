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

func (lar *LoginAttemptRepo) Create(ctx context.Context, id model.ID) error {
	attempt := sqlmodel.LoginAttempt{
		UserID:      int64(id),
		Attempts:    0,
		LastAttempt: time.Now().UTC(),
		BanExpiry:   time.Now().UTC(),
	}

	_, err := lar.db.NewInsert().Model(&attempt).Exec(ctx)
	return err
}

func (lar *LoginAttemptRepo) ByID(ctx context.Context, id model.ID) (model.LoginAttempt, error) {
	var attempt sqlmodel.LoginAttempt
	err := lar.db.NewSelect().Model(&attempt).Where("user_id = ?", id).Scan(ctx)
	if err != nil {
		return model.LoginAttempt{}, err
	}

	return model.LoginAttempt{
		ID:          model.ID(attempt.UserID),
		Attempts:    attempt.Attempts,
		LastAttempt: attempt.LastAttempt,
		BanExpiry:   attempt.BanExpiry,
	}, nil
}

func (lar *LoginAttemptRepo) Update(ctx context.Context, a model.LoginAttempt) error {
	var attempt sqlmodel.LoginAttempt
	err := lar.db.NewSelect().Model(&attempt).Where("user_id = ?", a.ID).Scan(ctx)
	if err != nil {
		return err
	}

	attempt.Attempts = a.Attempts
	attempt.LastAttempt = a.LastAttempt
	attempt.BanExpiry = a.BanExpiry

	_, err = lar.db.NewUpdate().Model(&attempt).Where("id = ?", attempt.ID).Exec(ctx)
	return err
}
