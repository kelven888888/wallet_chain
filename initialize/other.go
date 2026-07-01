package initialize

import (
	"context"
	"fmt"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"os"
	"strings"
	"wallet_chain.com/admin/model"
	"wallet_chain.com/heth"
	"wallet_chain.com/trx"

	"wallet_chain.com/global"
	"wallet_chain.com/utils"
)

func OtherInit() {
	dr, err := utils.ParseDuration(global.SHOP_CONFIG.JWT.ExpiresTime)
	if err != nil {
		panic(err)
	}
	global.BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(dr),
	)
	allowlanguage := strings.Split(global.SHOP_CONFIG.System.Language_Array, ",")
	for _, v := range allowlanguage {
		jsons, err := utils.ReadFileContent(fmt.Sprintf("./language/%s.json", v))
		if err != nil {
			fmt.Println(err.Error())
		}
		result := strings.Join(jsons, "")                                   // 使用空格作为分隔符
		global.BlackCache.SetDefault(fmt.Sprintf("language_%s", v), result) // 输
	}
	//启动删除所有锁
	type TAppLock struct {
		Id         int32  `json:"id" gorm:"column:id;primaryKey;not null;type:int(11)"`
		K          string `json:"k" gorm:"column:k;not null;type:varchar(64)"`
		V          int8   `json:"v" gorm:"column:v;not null;type:tinyint(2)"`
		CreateTime int64  `json:"create_time" gorm:"column:create_time;not null;type:bigint(20)"`
	}

}
func Walletinit() {

	//启动删除所有锁
	type TAppLock struct {
		Id         int32  `json:"id" gorm:"column:id;primaryKey;not null;type:int(11)"`
		K          string `json:"k" gorm:"column:k;not null;type:varchar(64)"`
		V          int8   `json:"v" gorm:"column:v;not null;type:tinyint(2)"`
		CreateTime int64  `json:"create_time" gorm:"column:create_time;not null;type:bigint(20)"`
	}
	var tapplock []TAppLock
	global.SHOP_DB.Table("t_app_lock").Find(&tapplock)

	for _, v := range tapplock {
		global.SHOP_REDIS.Del(context.Background(), v.K)
	}
	var tokens []model.TAppConfigToken
	global.SHOP_DB.Model(model.TAppConfigToken{}).Find(&tokens)
	for _, v := range tokens {
		var account model.Account
		//查看是否有热钱包
		global.SHOP_DB.Model(model.Account{}).Where("chain=? and account_type=2", v.Chain).Find(&account)
		if account.Id == 0 {
			if v.Chain == "eth" {
				_, err := heth.CreateHotAddresseth(1)
				if err != nil {
					fmt.Println("生成热钱包失败")
					os.Exit(0)
				}
				global.SHOP_DB.Model(model.Account{}).Where("chain=? and account_type=2", v.Chain).Limit(1).Find(&account)
				global.SHOP_DB.Model(model.TAppConfigToken{}).Where("chain=?", v.Chain).Updates(model.TAppConfigToken{
					HotAddress: account.Address,
				})
			}
			if v.Chain == "trx" {
				address, err := trx.CreateHotAddresstrx()
				if err != nil {
					fmt.Println("生成热钱包失败")
					os.Exit(0)
				}

				global.SHOP_DB.Model(model.TAppConfigToken{}).Where("chain=?", v.Chain).Updates(model.TAppConfigToken{
					HotAddress: address.Address,
				})
			}

		}
	}

}
