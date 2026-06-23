package service

import (
	"errors"
    	"fmt"
    	"wallet_chain.com/admin/model"
    	"wallet_chain.com/admin/model/common/request"
    	"wallet_chain.com/global"
    	"time"
)

type SBanner struct {
}

func (this *SBanner) GetAll(pageInfo request.PageInfo) ([]model.Banner, int64) {
	var models []model.Banner

    	query := global.SHOP_DB.Model(model.Banner{})

    	if pageInfo.Keyword != "" {

        		query.Where(fmt.Sprintf("%s like '%s'", pageInfo.SearchField, "%%"+pageInfo.Keyword+"%%"))
        	}

    	var count int64 = 0
    	query.Count(&count)
    	err := query.Limit(pageInfo.Limit).Offset((pageInfo.Page - 1) * pageInfo.Limit).Order(" id DESC").Find(&models).Error
    	if err != nil {
    		return nil, 0
    	}

    	return models, count

}
func (this *SBanner) GetByID(id request.GetById) (*model.Banner, error) {
	var models *model.Banner
	err := global.SHOP_DB.First(&models, id).Error
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return models, nil
}
func (this *SBanner) Save(models *model.Banner) error {
	if models.Id > 0 {
		return global.SHOP_DB.Updates(&models).Error
	} else {
	    models.CreateTime=time.Now()
		return global.SHOP_DB.Save(&models).Error
	}

}
func (this *SBanner) Delete(id uint32) error {

	return global.SHOP_DB.Where("id=?", id).Delete(&model.Banner{}).Error
}
func (this *SBanner) Deletebatch(req request.IdsReq) error {
	var models []model.Banner
	err := global.SHOP_DB.Find(&models, "id in ?", req.Ids).Delete(&models).Error

	if err != nil {
		return err
	}

	return nil
}
