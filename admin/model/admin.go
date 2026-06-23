package model

import "time"

type Admin struct {
	Id            uint       `json:"id" form:"id" `
	Account       string     `json:"account" gorm:"account" comment:"用户账号"`
	Mail          string     `json:"mail" comment:"用户邮箱"`
	Name          string     `comment:"用户昵称"`
	Mobile        uint64     `json:"mobile" comment:"用户手机号码"`
	Password      string     `json:"password" comment:"用户密码"`
	GroupId       uint       `comment:"用户所属群组" json:"group_id,string"  gorm:"group_id"`
	Status        *int       `comment:"用户状态" json:"status" gorm:"status" `
	LoginVisit    uint       `comment:"登录次数"`
	LastLoginIp   string     `comment:"最后登录IP"`
	LastLoginedAt *time.Time `comment:"最后登录时间"`
	Groupname     string     `json:"groupname" comment:"群组名" gorm:"->" `
	RoleIds       string     `json:"-" gorm:"-"`
	ModelTime
}

func (Admin) TableName() string {
	return "nov_admin"
}
