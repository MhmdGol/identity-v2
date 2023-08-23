package model

import (
	"time"
)

type RawUser struct {
	ID       ID
	UUN      string
	Username string
	Password string
	Email    string
	Role     string
	Status   string
}

type UserInfo struct {
	ID             ID
	UUN            string
	Username       string
	HashedPassword string
	Email          string
	Created_at     time.Time
	TOTPIsActive   bool
	TOTPSecret     string
	Role           string
	Status         string
}

type LoginInfo struct {
	Email    string
	Password string
	TOTPCode string
}
