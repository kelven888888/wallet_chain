package model

import "time"

type Heartbeat struct {
	Id        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Type      int
}

func (Heartbeat) TableName() string {
	return "heartbeat"
}
