package model

import (
	"time"
)

type Language struct {
	LangId   int
	Name     string
	Code     string
	Status   string
	CreateAt time.Time
	OrderBy  int
}

func (*Language) TableName() string {
	return "language"
}
