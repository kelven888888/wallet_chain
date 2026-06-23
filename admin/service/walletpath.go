package service

import (
	"fmt"
	"wallet_chain.com/admin/model"
	"wallet_chain.com/admin/model/common/request"
	"wallet_chain.com/global"
)

type WalletPath struct {
}

func (this *WalletPath) GetAll(pageInfo request.PageInfo) ([]model.WalletPath, int64) {
	var models []model.WalletPath

	query := global.SHOP_DB.Model(model.WalletPath{})
	if pageInfo.Keyword != "" {

		query.Where("username LIKE ? or  wallet_path like ?", "%"+pageInfo.Keyword+"%", "%"+pageInfo.Keyword+"%")
	}
	if pageInfo.Status != "" {

		query.Where("status =? ", pageInfo.Status)
	}

	var count int64 = 0
	query.Count(&count)
	err := query.Limit(pageInfo.Limit).Offset((pageInfo.Page - 1) * pageInfo.Limit).Order(" id desc").Find(&models).Error
	if err != nil {
		fmt.Println(err.Error())
		return nil, 0
	}

	return models, count

}
