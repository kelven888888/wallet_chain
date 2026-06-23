package model

import "time"

type AccountUserMessage struct {
	Id         int64
	CreateTime time.Time
	UpdateTime time.Time
	Remarks    string
	Username   string
	Title      string
	Content    string
	Status     int
	Group      int
	Type       int
	Read       int8
	Extends    *string `gorm:"type:json"`
	//ModelTime
}

func (*AccountUserMessage) TableName() string {
	return "account_user_message"
}
