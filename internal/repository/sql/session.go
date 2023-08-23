package sql

import (
	"context"
	"identity-v2/internal/model"
	"identity-v2/internal/repository"
	"identity-v2/internal/repository/sql/sqlmodel"

	"github.com/uptrace/bun"
)

type SessionRepo struct {
	db *bun.DB
}

var _ repository.SessionRepo = (*SessionRepo)(nil)

func NewSessionRepo(db *bun.DB) *SessionRepo {
	return &SessionRepo{
		db: db,
	}
}

func (sr *SessionRepo) Add(ctx context.Context, s model.Session) error {
	session := sqlmodel.Session{
		UserID: int64(s.UserID),
		Exp:    s.SessionExp,
	}

	_, err := sr.db.NewInsert().Model(&session).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (sr *SessionRepo) Remove(ctx context.Context, id model.ID) error {
	_, err := sr.db.NewDelete().Model((*sqlmodel.Session)(nil)).Where("user_id = ?", id).Exec(ctx)

	return err
}

func (sr *SessionRepo) ByID(ctx context.Context, id model.ID) (model.Session, error) {
	var session sqlmodel.Session
	_, err := sr.db.NewSelect().Model(&session).Where("user_id = ?", id).Exec(ctx)
	if err != nil {
		return model.Session{}, err
	}

	return model.Session{
		UserID:     model.ID(session.UserID),
		SessionExp: session.Exp,
	}, nil
}
