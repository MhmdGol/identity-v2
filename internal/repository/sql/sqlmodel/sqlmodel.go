package sqlmodel

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID             int64
	UUN            string
	Username       string
	HashedPassword string
	Email          string
	Created_at     time.Time
	TOTPSecret     string
	Role           int32 `bun:"rel:belongs-to"`
	Status         int32 `bun:"rel:belongs-to"`
}

type Role struct {
	bun.BaseModel `bun:"table:roles"`

	ID   int32 `bun:"id,pk,autoincrement"`
	Name string
}

type Status struct {
	bun.BaseModel `bun:"table:statuses"`

	ID   int32 `bun:"id,pk,autoincrement"`
	Name string
}

type Track struct {
	bun.BaseModel `bun:"table:tracks"`

	ID     int32 `bun:"id,pk,autoincrement"`
	UserID int64
	Action *Action
}

type Action struct {
	bun.BaseModel `bun:"table:actions"`

	ID   int32 `bun:"id,pk,autoincrement"`
	Name string
}

type Session struct {
	bun.BaseModel `bun:"table:sessions"`

	ID     int32 `bun:"id,pk,autoincrement"`
	UserID int64
	Exp    time.Time
}
