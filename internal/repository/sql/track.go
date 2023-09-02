package sql

import (
	"context"
	"identity-v2/internal/model"
	"identity-v2/internal/repository"
	"identity-v2/internal/repository/sql/sqlmodel"

	"github.com/uptrace/bun"
)

type TrackRepo struct {
	db *bun.DB
}

var _ repository.TrackRepo = (*TrackRepo)(nil)

func NewTrackRepo(db *bun.DB) *TrackRepo {
	return &TrackRepo{
		db: db,
	}
}

func (tr *TrackRepo) Create(ctx context.Context, t model.TrackInfo) error {
	var action sqlmodel.Action
	err := tr.db.NewSelect().Model(&action).Where("name = ?", t.Action).Scan(ctx)
	if err != nil {
		return err
	}

	newTrack := sqlmodel.Track{
		UserID:     int64(t.ID),
		ActionID:   action.ID,
		ActionTime: t.Timestamp,
	}

	_, err = tr.db.NewInsert().Model(&newTrack).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
