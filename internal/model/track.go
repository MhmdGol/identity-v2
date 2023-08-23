package model

import "time"

type TrackInfo struct {
	ID        ID
	Action    string
	Timestamp time.Time
}
