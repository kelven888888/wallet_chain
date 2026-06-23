package model

import "time"

type CrontabLog struct {
	Id        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Remark    string `comment:"任务名称"`
	Type      int    `comment:"1 ai结算 2 更新netval"`
	TDate     string `comment:"时间"`
}

func (CrontabLog) TableName() string {
	return "quan_crontab_log"
}
