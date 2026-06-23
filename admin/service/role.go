package service

import (
	"errors"
	"fmt"
	"strings"
	"wallet_chain.com/admin/model"
	"wallet_chain.com/admin/model/common/request"
	"wallet_chain.com/global"
)

type Role struct {
}

// 获取菜单分类
func (this *Role) GetMenus(roleIds []string, types int) ([]*model.Role, error) {
	var menus []*model.Role
	if types == 2 {
		global.SHOP_DB.Where("is_menu= ? and id in ? and deleted_at is  null ", 1, roleIds).Order(" sort asc").Find(&menus)
	} else {
		global.SHOP_DB.Where("is_menu=? and deleted_at is  null", 1).Order(" sort asc").Find(&menus)
	}
	//fmt.Printf("%#v", menus)
	//if err != nil {
	//	return menus, errors.New(err.Error())
	//}
	tree := this.getTree(menus)

	return tree, nil
}
func (this *Role) Delete(id uint32) error {
	global.SHOP_DB.Where("pid=?", id).Delete(&model.Role{})
	return global.SHOP_DB.Where("id=?", id).Delete(&model.Role{}).Error
}
func (this *Role) Save(role *model.Role) error {

	return global.SHOP_DB.Save(&role).Error
}

// 获取单个权限
func (this *Role) GetRole(id uint32) (*model.Role, error) {

	var role *model.Role
	if id < 0 {
		return role, errors.New("id parameter")
	}

	err := global.SHOP_DB.Where("id=?", id).First(&role).Error
	if err != nil {
		return role, err
	}

	return role, nil
}

// 获取权限分类
func (this *Role) GetRoles(isTree ...bool) []*model.Role {
	var menus []*model.Role
	global.SHOP_DB.Find(&menus)

	if len(isTree) > 0 {
		menus = this.getTree(menus)
	}

	return menus
}

// 生成树分类
func (this *Role) getTree(menus []*model.Role) []*model.Role {
	tmpMap := make(map[uint32]*model.Role)
	var tree []*model.Role
	for _, menu := range menus {
		tmpMap[menu.Id] = menu
	}

	for _, menu := range menus {
		if _, ok := tmpMap[menu.Pid]; ok {
			tmpMap[menu.Pid].Child = append(tmpMap[menu.Pid].Child, tmpMap[menu.Id])
		} else {
			tree = append(tree, tmpMap[menu.Id])
		}
	}
	return tree
}
func (this *Role) GetAll(req request.PageInfo) ([]*model.Role, int64) {
	var role []*model.Role
	query := global.SHOP_DB.Model(model.Role{})
	if req.Keyword != "" {
		query = query.Where("name like?", "%"+req.Keyword+"%")
	}
	var count int64
	query.Count(&count)
	err := query.Limit(req.Limit).Offset((req.Page - 1) * req.Limit).Order("id ASC").Find(&role).Error

	if err != nil {
		return nil, 0
	}

	return role, count
}

type Roles struct {
	ID     int
	Module string
	Action string
}

func Getmenu() {
	//var group = model.Group{}
	var role []Roles
	err := global.SHOP_DB.Model(&model.Role{}).Select("ID,Module,Action").Find(&role).Error
	if err != nil {
		panic(err)
	}

	for _, v := range role {
		key := strings.ToLower(v.Module + v.Action)
		//fmt.Println(key)
		global.BlackCache.SetDefault(key, v.ID)
	}

}
func (this *Role) AuthCheck(path string, group_id uint) bool {
	group := new(model.Group)
	err := global.SHOP_DB.Model(model.Group{}).Where("id=?", group_id).Select("role_ids").Find(&group).Error
	if err != nil {
		fmt.Println(err.Error())
	}

	var role []model.Role
	var MenuIds []string

	MenuIds = strings.Split(group.RoleIds, ",")

	err = global.SHOP_DB.Model(&model.Role{}).Where("id in (?)", MenuIds).Select("ID,Module,Action").Find(&role).Error
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, v := range role {

		roles := strings.ToLower(v.Module + v.Action)

		if roles == path {
			return true
		}
	}
	return false

}
