package model

import "time"

type Session struct {
	UserID     ID
	SessionExp time.Time
}
