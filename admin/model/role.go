package model

import (
	"fmt"
	"strings"
)

type Role struct {
	Id        uint32  `form:"id"`
	Pid       uint32  `form:"pid" comment:"父模块ID"`
	Name      string  `form:"name" comment:"模块名称"`
	Icon      string  `form:"icon" comment:"菜单icon"`
	IsMenu    uint8   `form:"is_menu" comment:"是否左侧菜单显示"`
	Desc      string  `form:"desc" comment:"模块说明"`
	Module    string  `form:"module" comment:"模块"`
	Action    string  `form:"action" comment:"方法"`
	Sort      uint8   `form:"sort" comment:"排序"`
	IsDefault uint8   `form:"is_default" comment:"是否默认模块"`
	Child     []*Role `gorm:"-"`
	ModelTime
}

func (*Role) TableName() string {
	return "nov_role"
}

// 拼接URL地址
func (m Role) Url() string {
	if m.Module == "" {
		return ""
	} else {
		return fmt.Sprintf("../../admin/%s/%s", strings.ToLower(m.Module), strings.ToLower(m.Action))
	}
}

// 获取是否菜单
func (m Role) IsMenuName() string {
	if m.IsMenu == 0 {
		return `<span class="layui-btn layui-btn-primary layui-btn-xs">否</span>`
	}

	return `<span class="layui-btn layui-btn-xs">是</span>`
}
