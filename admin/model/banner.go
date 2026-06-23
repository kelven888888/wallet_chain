package model

import "time"

type Banners struct {
	Id         int64     `json:"id"  form:"id,omitempty"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;autoUpdateTime"`
	Remarks    string
	Title      string `json:"title"  form:"title"`
	PointUrl   string `json:"point_url"  form:"point_url"`
	Image      string `json:"image"  form:"image"`
	Sort       int    `json:"sort"  form:"sort"`
	Status     *int   `json:"status"  form:"status"`
	Language   string `json:"language" form:"language" gorm:"column:language"`
}

func (Banners) TableName() string {
	return "mobile_banner_image"
}
