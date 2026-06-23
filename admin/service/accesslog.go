package service

import (
	"fmt"
	"wallet_chain.com/admin/model"
	"wallet_chain.com/admin/model/common/request"
	"wallet_chain.com/global"
)

type Saccesslog struct {
}

func (this *Saccesslog) GetAll(pageInfo request.PageInfo) ([]model.MAccesslog, int64) {
	var models []model.MAccesslog

	query := global.SHOP_DB.Model(model.MAccesslog{})
	if pageInfo.Keyword != "" {

		query.Where("ip LIKE ? or username like ?  ", "%"+pageInfo.Keyword+"%", "%"+pageInfo.Keyword+"%")
	}
	if pageInfo.Status != "" {

		query.Where("username =? ", pageInfo.Status)
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
