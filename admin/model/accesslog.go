package model

import (
	"time"
)

type MAccesslog struct {
	Id          int
	Path        string
	Ip          string
	CreateAt    time.Time
	Method      string
	Address     string
	Country     string
	City        string
	Subdivision string
	Username    string

	//ModelTime
}

func (*MAccesslog) TableName() string {
	return "access_log"
}
