package app

import (
	"context"
	"github.com/moremorefun/mcommon"
	"time"
	admodel "wallet_chain.com/admin/model"
	"wallet_chain.com/global"
	"wallet_chain.com/model"
	"wallet_chain.com/utils"
)

// GetLock 获取运行锁
func GetLock(ctx context.Context, tx mcommon.DbExeAble, k string) (bool, error) {
	genLock := func() error {
		_, err := SQLCreateTAppLockUpdate(
			ctx,
			tx,
			&model.DBTAppLock{
				K:          k,
				V:          1,
				CreateTime: time.Now().Unix(),
			},
		)
		if err != nil {
			return err
		}
		return nil
	}

	lockRow, err := SQLGetTAppLockColByK(
		ctx,
		tx,
		[]string{
			model.DBColTAppLockCreateTime,
		},
		k,
	)
	if err != nil {
		return false, err
	}
	if lockRow == nil {
		err = genLock()
		if err != nil {
			return false, err
		}
		return true, nil
	}
	if time.Now().Unix()-lockRow.CreateTime > 60*30 {
		err = genLock()
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// ReleaseLock 释放运行锁
func ReleaseLock(ctx context.Context, tx mcommon.DbExeAble, k string) error {
	_, err := SQLUpdateTAppLockByK(
		ctx,
		tx,
		&model.DBTAppLock{
			K:          k,
			V:          0,
			CreateTime: time.Now().Unix(),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// LockWrap 包装被lock的函数
func LockWrap(name string, f func()) {
	rdb := global.SHOP_REDIS
	ok := utils.AcquireLock(rdb, name, time.Second*600)
	if !ok {
		return
	}
	defer func() {
		utils.ReleaseLock(rdb, name)

	}()
	f()
}

// SQLGetWithdrawMap 获取提币map
func SQLGetWithdrawMap(ids []int64) (map[int64]*admodel.TWithdraw, error) {
	itemMaps := make(map[int64]*admodel.TWithdraw)
	var withdraw []*admodel.TWithdraw
	global.SHOP_DB.Model(admodel.TWithdraw{}).Where("id in ?", ids).Find(&withdraw)
	for _, itemRow := range withdraw {
		itemMaps[itemRow.Id] = itemRow
	}
	return itemMaps, nil
}

// SQLGetProductMap 获取产品map
func SQLGetProductMap(ctx context.Context, tx mcommon.DbExeAble, cols []string, ids []int64) (map[int64]*model.DBTProduct, error) {
	if !mcommon.IsStringInSlice(cols, model.DBColTProductID) {
		cols = append(cols, model.DBColTProductID)
	}
	itemMap := make(map[int64]*model.DBTProduct)
	itemRows, err := model.SQLSelectTProductCol(
		ctx,
		tx,
		cols,
		ids,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	for _, itemRow := range itemRows {
		itemMap[itemRow.ID] = itemRow
	}
	return itemMap, nil
}

// SQLGetAppConfigTokenMap 获取代币map
func SQLGetAppConfigTokenMap() (map[int64]*admodel.TAppConfigToken, error) {
	//if !mcommon.IsStringInSlice(cols, model.DBColTAppConfigTokenID) {
	//	cols = append(cols, model.DBColTAppConfigTokenID)
	//}
	//itemMap := make(map[int64]*model.DBTAppConfigToken)
	//itemRows, err := model.SQLSelectTAppConfigTokenCol(
	//	ctx,
	//	tx,
	//	cols,
	//	ids,
	//	nil,
	//	nil,
	//)
	//if err != nil {
	//	return nil, err
	//}
	//for _, itemRow := range itemRows {
	//	itemMap[itemRow.ID] = itemRow
	//}
	var modeltoken []*admodel.TAppConfigToken
	itemMaps := make(map[int64]*admodel.TAppConfigToken)
	global.SHOP_DB.Model(admodel.TAppConfigToken{}).Find(&modeltoken)
	for _, v := range modeltoken {
		itemMaps[v.Id] = v
	}
	return itemMaps, nil
}

// SQLGetAddressKeyMap 获取地址map
func SQLGetAddressKeyMap(ctx context.Context, tx mcommon.DbExeAble, cols []string, addresses []string) (map[string]*admodel.Account, error) {
	if !mcommon.IsStringInSlice(cols, model.DBColTAddressKeyAddress) {
		cols = append(cols, model.DBColTAddressKeyAddress)
	}
	itemMap := make(map[string]*admodel.Account)
	itemRows, err := SQLSelectTAddressKeyColByAddress(
		addresses,
	)
	if err != nil {
		return nil, err
	}
	for _, itemRow := range itemRows {
		itemMap[itemRow.Address] = itemRow
	}
	return itemMap, nil
}
