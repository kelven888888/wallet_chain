package app

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
	dbmodel "wallet_chain.com/admin/model"
	"wallet_chain.com/global"
	"wallet_chain.com/model"

	"github.com/gin-gonic/gin"
	"github.com/moremorefun/mcommon"
)

// SQLGetTAppConfigIntValueByK 查询配置
func SQLGetTAppConfigIntValueByK(ctx context.Context, tx mcommon.DbExeAble, k string) (int64, error) {
	var parm dbmodel.OtherParam
	global.SHOP_DB.Model(dbmodel.OtherParam{}).Where("`key`=?", k).Find(&parm)
	if parm.Value == "" {
		return 0, fmt.Errorf("no app config int of: %s", k)
	}
	value, _ := strconv.ParseInt(parm.Value, 10, 64)
	return value, nil
}

// SQLGetTAppConfigStrValueByK 查询配置
func SQLGetTAppConfigStrValueByK(ctx context.Context, tx mcommon.DbExeAble, k string) (string, error) {
	var parm dbmodel.OtherParam
	global.SHOP_DB.Model(dbmodel.OtherParam{}).Where("`key`=?", k).Find(&parm)
	if parm.Value == "" {
		return "", fmt.Errorf("no app config int of: %s", k)
	}

	return parm.Value, nil

}

// SQLGetTAppConfigStrValueByK 查询配置
func SQLGetTProduct() (dbmodel.TProduct, error) {
	var parm dbmodel.TProduct
	global.SHOP_DB.Model(dbmodel.TProduct{}).Find(&parm)
	if parm.Id == 0 {
		return parm, nil
	}

	return parm, nil

}
func SQLethGetHotADDRESSS(k string) ([]dbmodel.Account, error) {
	var account []dbmodel.Account
	global.SHOP_DB.Model(dbmodel.Account{}).Where("`chain`=?", k).Find(&account)

	return account, nil

}
func SQLethGetHotADDRESS(k string) (string, error) {
	var TAppConfigToken dbmodel.TAppConfigToken
	global.SHOP_DB.Model(dbmodel.TAppConfigToken{}).Where("`chain`=?", k).Limit(1).Find(&TAppConfigToken)

	return TAppConfigToken.HotAddress, nil

}
func UpdateTransationhandeltimeadd(id int64) error {
	var transactions dbmodel.Transactions
	err := global.SHOP_DB.Model(dbmodel.Transactions{}).Where("`id`=?", id).Find(&transactions).Updates(dbmodel.Transactions{
		HandelTimes: transactions.HandelTimes + 1,
	}).Error

	if err != nil {

		global.SHOP_LOG.Error(err.Error())
		return err
	}
	return nil

}

// SQLGetTAppStatusIntValueByK
func SQLGetTAppStatusIntValueByK(ctx context.Context, tx mcommon.DbExeAble, k string) (int64, error) {
	var parm dbmodel.OtherParam
	global.SHOP_DB.Model(dbmodel.OtherParam{}).Where("`key`=?", k).Find(&parm)
	if parm.Value == "" {
		return 0, fmt.Errorf("no app config int of: %s", k)
	}
	value, _ := strconv.ParseInt(parm.Value, 10, 64)
	return value, nil
}

// SQLGetTAddressKeyFreeCount 获取剩余可用地址数
func SQLGetTAddressKeyFreeCount(ctx context.Context, tx mcommon.DbExeAble, symbol string) (int64, error) {
	var count int64
	global.SHOP_DB.Model(dbmodel.Account{}).Where("chain=? and status=-1 and account_type=1", symbol).Count(&count)

	return count, nil
}

// SQLSelectTAddressKeyColByAddress 根据ids获取
func SQLSelectTAddressKeyColByAddress(addresses []string) ([]*dbmodel.Account, error) {
	if len(addresses) == 0 {
		return nil, nil
	}
	var rows []*dbmodel.Account
	global.SHOP_DB.Where("address in?", addresses).Find(&rows)
	//global.SHOP_DB.Find(&rows)

	return rows, nil
}
func SQLgetEXitTRAN(txid string) (bool, error) {
	if len(txid) == 0 {
		return false, nil
	}
	var rows dbmodel.Transactions
	global.SHOP_DB.Where("tx_id =?", txid).Find(&rows)
	//global.SHOP_DB.Find(&rows)

	if rows.Id > 0 {
		return true, nil
	}
	return false, nil
}

// SQLUpdateTAppStatusIntByK 更新
func SQLUpdateTAppStatusIntByK(ctx context.Context, tx mcommon.DbExeAble, row *model.DBTAppStatusInt) (int64, error) {
	var parms dbmodel.OtherParam
	err := global.SHOP_DB.Model(dbmodel.OtherParam{}).Where("`key`=?", row.K).Find(&parms).Updates(dbmodel.OtherParam{
		Value: strconv.Itoa(int(row.V)),
	}).Error
	if err != nil {
		return 0, err
	}
	return int64(parms.Id), nil
}

// SQLUpdateTAppStatusIntByKGreater 更新
func SQLUpdateTAppStatusIntByKGreater(ctx context.Context, tx mcommon.DbExeAble, row *model.DBTAppStatusInt, rows *model.DBTAppStatusInt) (int64, error) {
	err := global.SHOP_DB.Model(dbmodel.OtherParam{}).Where("`key`=?", row.K).Updates(dbmodel.OtherParam{
		Value: strconv.Itoa(int(row.V)),
	}).Error
	if err != nil {
		return 0, err
	}
	err = global.SHOP_DB.Model(dbmodel.OtherParam{}).Where("`key`=?", rows.K).Updates(dbmodel.OtherParam{
		Value: strconv.Itoa(int(rows.V)),
	}).Error
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// SQLGetTSendMaxNonce 获取地址的nonce
func SQLGetTSendMaxNonce(address string) (int64, error) {
	var i int64

	err := global.SHOP_DB.Raw("SELECT IFNULL(MAX(nonce), -1) FROM t_send WHERE from_address=? LIMIT 1", address).Scan(&i).Error
	if err != nil {
		global.SHOP_LOG.Error(err.Error())
		return 0, nil
	}

	return i + 1, nil
}

// SQLGetTSendPendingBalanceReal 获取地址的打包数额
func SQLGetTSendPendingBalanceReal(ctx context.Context, tx mcommon.DbExeAble, address string) (string, error) {
	i := "0"

	query := `SELECT 
	IFNULL(SUM(CAST(balance_real as DECIMAL(65,18))), "0")
FROM
	t_send
WHERE
	from_address=?
	AND handle_status<2
LIMIT 1`
	global.SHOP_DB.Raw(query, address).Scan(&i)
	return i, nil
}

// SQLGetTAddressKeyColByAddress 根据address查询
func SQLGetTAddressKeyColByAddress(ctx context.Context, tx mcommon.DbExeAble, cols []string, address string) (*dbmodel.Account, error) {
	var account dbmodel.Account
	global.SHOP_DB.Where("address=?", address).Find(&account)
	if account.Id == 0 {
		return nil, nil
	}
	return &account, nil
}

// SQLUpdateTTxOrgStatusByIDs 更新
func SQLUpdateTTxOrgStatusByIDs(ctx context.Context, tx mcommon.DbExeAble, ids []int64, row model.DBTTx) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_tx
SET
    org_status=:org_status,
    org_msg=:org_msg,
    org_time=:org_time
WHERE
	id IN (:ids)`,
		gin.H{
			"ids":        ids,
			"org_status": row.OrgStatus,
			"org_msg":    row.OrgMsg,
			"org_time":   row.OrgTime,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLUpdateTTxStatusByIDs 更新
func SQLUpdateTTxStatusByIDs(ctx context.Context, tx mcommon.DbExeAble, ids []int64, row model.DBTTx) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_tx
SET
    handle_status=:handle_status,
    handle_msg=:handle_msg,
    handle_time=:handle_time
WHERE
	id IN (:ids)`,
		gin.H{
			"ids":           ids,
			"handle_status": row.HandleStatus,
			"handle_msg":    row.HandleMsg,
			"handle_time":   row.HandleTime,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLUpdateTSendStatusByIDs 更新
func SQLUpdateTSendStatusByIDs(ctx context.Context, tx mcommon.DbExeAble, ids []int64, row model.DBTSend) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_send
SET
    handle_status=:handle_status,
    handle_msg=:handle_msg,
    handle_time=:handle_time
WHERE
	id IN (:ids)`,
		gin.H{
			"ids":           ids,
			"handle_status": row.HandleStatus,
			"handle_msg":    row.HandleMsg,
			"handle_time":   row.HandleTime,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLSelectTSendColByStatus 根据ids获取
func SQLSelectTSendColByStatus(status int64) ([]*dbmodel.TSend, error) {
	//	query := strings.Builder{}
	//	query.WriteString("SELECT\n")
	//	query.WriteString(strings.Join(cols, ",\n"))
	//	query.WriteString(`
	//FROM
	//	t_send
	//WHERE
	//	handle_status=:handle_status
	//ORDER BY id`)
	//
	//	var rows []*model.DBTSend
	//	err := mcommon.DbSelectNamedContent(
	//		ctx,
	//		tx,
	//		&rows,
	//		query.String(),
	//		gin.H{
	//			"handle_status": status,
	//		},
	//	)
	//	if err != nil {
	//		return nil, err
	//	}
	var send []*dbmodel.TSend
	err := global.SHOP_DB.Model(dbmodel.TSend{}).Where("handle_status=?", status).Find(&send).Error
	if err != nil {
		global.SHOP_LOG.Error(err.Error())
		return nil, err
	}

	return send, nil
}

// SQLSelectTWithdrawColByStatus 根据ids获取
func SQLSelectTWithdrawColByStatus(ctx context.Context, tx mcommon.DbExeAble, cols []string, status int64, symbols []string) ([]*model.DBTWithdraw, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_withdraw
WHERE
	handle_status=:handle_status
	AND symbol IN (:symbols)`)

	var rows []*model.DBTWithdraw
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		gin.H{
			"handle_status": status,
			"symbols":       symbols,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLGetTWithdrawColForUpdate 根据id查询
func SQLGetTWithdrawColForUpdate(ctx context.Context, tx mcommon.DbExeAble, cols []string, id int64, status int64) (*model.DBTWithdraw, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_withdraw
WHERE
	id=:id
	AND handle_status=:handle_status
FOR UPDATE`)

	var row model.DBTWithdraw
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		gin.H{
			"id":            id,
			"handle_status": status,
		},
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLUpdateTWithdrawGenTx 更新
func SQLUpdateTWithdrawGenTx(ctx context.Context, tx mcommon.DbExeAble, row *model.DBTWithdraw) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_withdraw
SET
    tx_hash=:tx_hash,
    handle_status=:handle_status,
    handle_msg=:handle_msg,
    handle_time=:handle_time
WHERE
	id=:id`,
		gin.H{
			"id":            row.ID,
			"tx_hash":       row.TxHash,
			"handle_status": row.HandleStatus,
			"handle_msg":    row.HandleMsg,
			"handle_time":   row.HandleTime,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLUpdateTWithdrawStatusByIDs 更新
func SQLUpdateTWithdrawStatusByIDs(ctx context.Context, tx mcommon.DbExeAble, ids []int64, row *model.DBTWithdraw) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_withdraw
SET
    handle_status=:handle_status,
    handle_msg=:handle_msg,
    handle_time=:handle_time
WHERE
	id IN (:ids)`,
		gin.H{
			"ids":           ids,
			"handle_status": row.HandleStatus,
			"handle_msg":    row.HandleMsg,
			"handle_time":   row.HandleTime,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTAppLockColByK 根据id查询
func SQLGetTAppLockColByK(ctx context.Context, tx mcommon.DbExeAble, cols []string, k string) (*model.DBTAppLock, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_lock
WHERE
	k=:k
	AND v=1`)

	var row model.DBTAppLock
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		gin.H{
			"k": k,
		},
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLCreateTAppLockUpdate 创建
func SQLCreateTAppLockUpdate(ctx context.Context, tx mcommon.DbExeAble, row *model.DBTAppLock) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = mcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_app_lock (
    id,
    k,
    v,
    create_time
) VALUES (
    :id,
    :k,
    :v,
    :create_time
) ON DUPLICATE KEY UPDATE 
	v=:v,
	create_time=:create_time`,
			gin.H{
				"id":          row.ID,
				"k":           row.K,
				"v":           row.V,
				"create_time": row.CreateTime,
			},
		)
	} else {
		lastID, err = mcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_app_lock (
    k,
    v,
    create_time
) VALUES (
    :k,
    :v,
    :create_time
) ON DUPLICATE KEY UPDATE 
	v=:v,
	create_time=:create_time`,
			gin.H{
				"k":           row.K,
				"v":           row.V,
				"create_time": row.CreateTime,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLUpdateTAppLockByK 更新
func SQLUpdateTAppLockByK(ctx context.Context, tx mcommon.DbExeAble, row *model.DBTAppLock) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_app_lock
SET
    v=:v,
    create_time=:create_time
WHERE
	k=:k`,
		gin.H{
			"id":          row.ID,
			"k":           row.K,
			"v":           row.V,
			"create_time": row.CreateTime,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLSelectTTxColByStatus 根据ids获取
func SQLSelectTTxColByStatus(ctx context.Context, tx mcommon.DbExeAble, cols []string, status int64) ([]*model.DBTTx, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx
WHERE
	handle_status=:handle_status`)

	var rows []*model.DBTTx
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		gin.H{
			"handle_status": status,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTProductNotifyColByStatusAndTime 根据ids获取
func SQLSelectTProductNotifyColByStatusAndTime(ctx context.Context, tx mcommon.DbExeAble, cols []string, status int64, t int64) ([]*model.DBTProductNotify, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_product_notify
WHERE
	handle_status=:handle_status
	AND update_time<:update_time`)

	var rows []*model.DBTProductNotify
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		gin.H{
			"handle_status": status,
			"update_time":   t,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTProductNotifyStatusByID 更新
func SQLUpdateTProductNotifyStatusByID(ctx context.Context, tx mcommon.DbExeAble, row *model.DBTProductNotify) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_product_notify
SET
    handle_status=:handle_status,
    handle_msg=:handle_msg,
    update_time=:update_time
WHERE
	id=:id`,
		gin.H{
			"id":            row.ID,
			"handle_status": row.HandleStatus,
			"handle_msg":    row.HandleMsg,
			"update_time":   row.UpdateTime,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLSelectTAppConfigTokenColAll 根据ids获取
func SQLSelectTAppConfigTokenColAll(CoinSymbol string) ([]*dbmodel.TAppConfigToken, error) {
	//	query := strings.Builder{}
	//	query.WriteString("SELECT\n")
	//	query.WriteString(strings.Join(cols, ",\n"))
	//	query.WriteString(`
	//FROM
	//	t_app_config_token`)
	//
	//	var rows []*model.DBTAppConfigToken
	//	err := mcommon.DbSelectNamedContent(
	//		ctx,
	//		tx,
	//		&rows,
	//		query.String(),
	//		gin.H{},
	//	)
	var mtoken []*dbmodel.TAppConfigToken
	err := global.SHOP_DB.Model(dbmodel.TAppConfigToken{}).Where("chain=?", CoinSymbol).Find(&mtoken).Error
	if err != nil {
		return nil, err
	}
	return mtoken, nil
}

// SQLSelectTTxErc20ColByStatus 根据ids获取
func SQLSelectTTxErc20ColByStatus(status int64) ([]*dbmodel.Transactions, error) {
	var rows []*dbmodel.Transactions
	global.SHOP_DB.Model(dbmodel.Transactions{}).Where("status=? and handel_times<3", status).Find(&rows)
	return rows, nil
}

// SQLSelectTTxErc20ColByOrgForUpdate 获取未整理交易
func SQLSelectTTxErc20ColByOrgForUpdate(orgStatuses []int64, token_symbol []string) ([]*dbmodel.Transactions, error) {

	var rows []*dbmodel.Transactions
	err := global.SHOP_DB.Model(dbmodel.Transactions{}).Where("org_status in ? and token_symbol in ?", orgStatuses, token_symbol).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTTxErc20OrgStatusByIDs 更新
func SQLUpdateTTxErc20OrgStatusByIDs(db *gorm.DB, ids []int64, row model.DBTTxErc20) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	//	count, err := mcommon.DbExecuteCountNamedContent(
	//		ctx,
	//		tx,
	//		`UPDATE
	//	t_tx_erc20
	//SET
	//    org_status=:org_status,
	//    org_msg=:org_msg,
	//    org_time=:org_time
	//WHERE
	//	id IN (:ids)`,
	//		gin.H{
	//			"ids":        ids,
	//			"org_status": row.OrgStatus,
	//			"org_msg":    row.OrgMsg,
	//			"org_time":   row.OrgTime,
	//		},
	//	)
	//	if err != nil {
	//		return 0, err
	//	}
	err := db.Model(dbmodel.Transactions{}).Where("id in ?", ids).Updates(
		dbmodel.Transactions{
			OrgStatus: row.OrgStatus,
			OrgMsg:    row.OrgMsg,
			OrgTime:   row.OrgTime,
		},
	).Error
	if err != nil {
		db.Rollback()
		global.SHOP_LOG.Error(err.Error())
		return 0, err
	}
	return int64(len(ids)), nil
}
