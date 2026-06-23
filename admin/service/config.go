package service

import (
	"errors"
	"wallet_chain.com/admin/model"
	"wallet_chain.com/global"
)

type Config struct {
}

func (this *Config) GetAll() (result map[string]string, err error) {
	var models []model.Config

	query := global.SHOP_DB

	err = query.Order(" id ASC").Find(&models).Error
	if err != nil {
		return result, err
	}
	results := make(map[string]string)
	for _, v := range models {
		results[v.Key] = v.Value
	}

	return results, nil

}

func (this *Config) Save(key string, val string) error {
	var models model.Config
	models.Value = val
	return global.SHOP_DB.Model(model.Config{}).Where(" `key`= ?", key).Updates(&models).Error

}
func (this *Config) GetKeyValue(key string) (string, error) {
	var models model.Config

	global.SHOP_DB.Model(model.Config{}).Where(" `key`= ?", key).Find(&models)
	if models.Id == 0 {
		return "0", errors.New("没有记录")
	}
	return models.Value, nil

}
