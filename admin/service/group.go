package service

import (
	"errors"
	"wallet_chain.com/admin/model"
	"wallet_chain.com/admin/model/common/request"
	"wallet_chain.com/global"
)

type GroupServer struct {
}

// 获取当前权限组菜单列表
func (this *GroupServer) GetMenus(groupId uint) ([]*model.Role, error) {
	var group model.Group

	err := global.SHOP_DB.Where("id = ?", groupId).First(&group).Error
	if err != nil {
		return nil, err
	}
	roleIds := []string{}
	Roleserver := Role{}
	if groupId != 1 {
		roleIds = group.GetRoleIds()
		result, err := Roleserver.GetMenus(roleIds, 2)
		if err == nil {
			return result, err
		}
		return result, nil
	} else {
		roleIds = group.GetRoleIds()
		result, err := Roleserver.GetMenus(roleIds, 1)
		if err == nil {
			return result, err
		}

	}
	return []*model.Role{}, nil

}
func (this *GroupServer) GetByID(req *request.GetById) (*model.Group, error) {
	var group *model.Group
	err := global.SHOP_DB.First(&group, req.Uint32()).Error
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return group, nil
}
func (this *GroupServer) Save(group *model.Group) error {
	if group.Id > 0 {
		return global.SHOP_DB.Updates(&group).Error
	} else {
		return global.SHOP_DB.Save(&group).Error
	}

}
func (this *GroupServer) Delete(id uint32) error {
	if id == 1 {
		return nil
	}
	return global.SHOP_DB.Where("id=?", id).Delete(&model.Group{}).Error
}
