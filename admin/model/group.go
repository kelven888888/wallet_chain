package model

import "strings"

type Group struct {
	Id         uint     `form:"id" json:"id"`
	Name       string   `form:"name" comment:"群组名称"`
	Desc       string   `form:"desc" comment:"群组说明"`
	RoleIds    string   `form:"role_ids" comment:"群组权限ID，多个,分割"`
	RoleIdssli []string `gorm:"-" form:"role_ids[]" json:"role_ids[]"`

	ModelTime
}

func (*Group) TableName() string {
	return "nov_group"
}

// 获取权限id列表
func (m Group) GetRoleIds() []string {
	return strings.Split(m.RoleIds, ",")
}

// 获取权限id列表
func (m Group) GetRoleString() string {
	return strings.Join(m.RoleIdssli, ",")
}
