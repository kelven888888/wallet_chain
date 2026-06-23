package service

import (
	"errors"
	"gorm.io/gorm"
	"wallet_chain.com/admin/model"
	"wallet_chain.com/admin/model/common/request"
	"wallet_chain.com/global"
)

type AdminlogServer struct {
}

func (this *AdminlogServer) GetAll(req request.PageInfo) ([]*model.SysOperationRecord, int64) {
	var adminlog []*model.SysOperationRecord
	err := global.SHOP_DB.Limit(req.Limit).Offset((req.Page - 1) * req.Limit).Order("id DESC").Find(&adminlog).Error
	if err != nil {
		return nil, 0
	}
	var count int64
	global.SHOP_DB.Model(model.SysOperationRecord{}).Count(&count)
	return adminlog, count
}

func (this *AdminlogServer) Delete(req request.GetById) error {
	var adminlog model.SysOperationRecord
	err := global.SHOP_DB.Where("id = ?", req.ID).First(&adminlog).Error
	if errors.Is(err, gorm.ErrRecordNotFound) { // api记录不存在
		return err
	}
	err = global.SHOP_DB.Delete(&adminlog).Error
	if err != nil {
		return err
	}

	return nil
}
func (this *AdminlogServer) Deletebatch(req request.IdsReq) error {
	var adminlog []model.SysOperationRecord
	err := global.SHOP_DB.Find(&adminlog, "id in ?", req.Ids).Delete(&adminlog).Error

	if err != nil {
		return err
	}

	return nil
}
