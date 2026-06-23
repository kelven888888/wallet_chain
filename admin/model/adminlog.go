package model

type AdminLog struct {
	Id      uint
	Uid     uint   `comment:"操作人UID"`
	Name    string `comment:"操作人名称"`
	Ip      string `comment:"操作人IP"`
	Type    uint   `comment:"日记类型"`
	Content string `comment:"日记内容"`
	ModelTime
}

func (AdminLog) TableName() string {
	return "nov_admin_log"
}
