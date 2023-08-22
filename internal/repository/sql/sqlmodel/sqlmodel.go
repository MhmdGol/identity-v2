package sqlmodel

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID             int64     `bun:"id,pk"`
	UUN            string    `bun:"uun"`
	Username       string    `bun:"username"`
	HashedPassword string    `bun:"hashed_password"`
	Email          string    `bun:"email"`
	Created_at     time.Time `bun:"created_at"`
	TOTPIsActive   bool      `bun:"totp_is_active"`
	TOTPSecret     string    `bun:"totp_secret"`
	Role           int32     `bun:"role_id"`
	Status         int32     `bun:"status_id"`
}

type Role struct {
	bun.BaseModel `bun:"table:roles"`

	ID   int32  `bun:"id,pk,autoincrement"`
	Name string `bun:"name"`
}

type Status struct {
	bun.BaseModel `bun:"table:statuses"`

	ID   int32  `bun:"id,pk,autoincrement"`
	Name string `bun:"name"`
}

type Track struct {
	bun.BaseModel `bun:"table:tracks"`

	ID     int32 `bun:"id,pk,autoincrement"`
	UserID int64 `bun:"user_id"`
	Action int32 `bun:"action"`
}

type Action struct {
	bun.BaseModel `bun:"table:actions"`

	ID   int32  `bun:"id,pk,autoincrement"`
	Name string `bun:"name"`
}

type Session struct {
	bun.BaseModel `bun:"table:sessions"`

	ID     int32     `bun:"id,pk,autoincrement"`
	UserID int64     `bun:"user_id"`
	Exp    time.Time `bun:"exp"`
}
