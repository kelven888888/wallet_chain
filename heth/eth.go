package heth

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/params"
	"math"
	"math/big"
	"net/http"
	"strings"
	"time"
	admodel "wallet_chain.com/admin/model"
	"wallet_chain.com/app"
	"wallet_chain.com/ethclient"
	"wallet_chain.com/global"
	"wallet_chain.com/model"
	"wallet_chain.com/utils"
	"wallet_chain.com/xenv"

	"github.com/moremorefun/mcommon"
	"github.com/parnurzeal/gorequest"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/gin-gonic/gin"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/crypto"
)

func genAddressAndAesKey() (string, string, error) {
	// 生成私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyStr := hexutil.Encode(privateKeyBytes)
	// 加密密钥
	privateKeyStrEn, err := mcommon.AesEncrypt(privateKeyStr, xenv.Cfg.AESKey)
	if err != nil {
		return "", "", err
	}
	// 获取地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", errors.New("can't change public key")
	}
	// 地址全部储存为小写方便处理
	address := AddressBytesToStr(crypto.PubkeyToAddress(*publicKeyECDSA))
	return address, privateKeyStrEn, nil
}

// CreateHotAddress 创建自用地址
func CreateHotAddresseth(num int64) ([]string, error) {

	var rowaccount []*admodel.Account
	var addresses []string
	// 遍历差值次数
	for i := int64(0); i < num; i++ {
		address, privateKeyStrEn, err := genAddressAndAesKey()
		if err != nil {
			return nil, err
		}
		// 存入待添加队列

		rowaccount = append(rowaccount, &admodel.Account{
			Chain:       CoinSymbol,
			Address:     address,
			PrivateKey:  privateKeyStrEn,
			Status:      -1,
			AccountType: 2,
		})
	}
	// 一次性将生成的地址存入数据库
	_, err := model.SQLCreateManyTAddressKey(
		context.Background(),
		xenv.DbCon,
		rowaccount,
		true,
	)
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

// CheckAddressFree 检测是否有充足的备用地址
func CheckAddressFree() {
	lockKey := "EthCheckAddressFree"
	app.LockWrap(lockKey, func() {
		// 获取配置 允许的最小剩余地址数
		minFreeCount, err := app.SQLGetTAppConfigIntValueByK(
			context.Background(),
			xenv.DbCon,
			"min_free_address",
		)

		// 获取当前剩余可用地址数
		freeCount, err := app.SQLGetTAddressKeyFreeCount(
			context.Background(),
			xenv.DbCon,
			CoinSymbol,
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		fmt.Println(freeCount, minFreeCount)
		// 如果数据库中剩余可用地址小于最小允许可用地址
		if freeCount < minFreeCount {

			var rowaccount []*admodel.Account
			// 遍历差值次数
			for i := int64(0); i < minFreeCount-freeCount; i++ {
				address, privateKeyStrEn, err := genAddressAndAesKey()
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}

				rowaccount = append(rowaccount, &admodel.Account{
					Chain:       CoinSymbol,
					Address:     address,
					PrivateKey:  privateKeyStrEn,
					Status:      0,
					AccountType: 1,
				})
			}
			// 一次性将生成的地址存入数据库
			_, err = model.SQLCreateManyTAddressKey(
				context.Background(),
				xenv.DbCon,
				rowaccount,
				true,
			)
			if err != nil {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
		}
	})
}

// CheckBlockSeek 检测到账
func CheckBlockSeek() {
	lockKey := "EthCheckBlockSeek"
	app.LockWrap(lockKey, func() {
		// 获取配置 延迟确认数
		confirmValue, err := app.SQLGetTAppConfigIntValueByK(
			context.Background(),
			xenv.DbCon,
			"eth_block_confirm_num",
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 获取状态 当前处理完成的最新的block number
		seekValue, err := app.SQLGetTAppStatusIntValueByK(
			context.Background(),
			xenv.DbCon,
			"seek_num",
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// rpc 获取当前最新区块数
		rpcBlockNum, err := ethclient.RpcBlockNumber(context.Background())
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		startI := seekValue + 1
		endI := rpcBlockNum - confirmValue + 1
		if startI < endI {
			// 手续费钱包列表
			feeAddressValue, err := app.SQLethGetHotADDRESSS(
				CoinSymbol,
			)
			if err != nil {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}

			var feeAddresses []string
			for _, address := range feeAddressValue {
				if address.Address == "" {
					continue
				}
				feeAddresses = append(feeAddresses, address.Address)
			}
			// 遍历获取需要查询的block信息
			for i := startI; i < endI; i++ {
				// rpc获取block信息
				//mcommon.Log.Debugf("eth check block: %d", i)
				rpcBlock, err := ethclient.RpcBlockByNum(context.Background(), i)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				// 接收地址列表
				var toAddresses []string
				// map[接收地址] => []交易信息
				toAddressTxMap := make(map[string][]*types.Transaction)
				// 遍历block中的tx
				for _, rpcTx := range rpcBlock.Transactions() {
					// 转账数额大于0 and 不是创建合约交易
					if rpcTx.Value().Cmp(big.NewInt(0)) > 0 && rpcTx.To() != nil {
						//msg, err := rpcTx.AsMessage(types.NewLondonSigner(rpcTx.ChainId()), nil)
						//if err != nil {
						//	mcommon.Log.Errorf("AsMessage err: [%T] %s", err, err.Error())
						//	return
						//}
						signer := types.MakeSigner(params.MainnetChainConfig, rpcBlock.Number(), rpcBlock.Time())

						from, err := signer.Sender(rpcTx)
						if err != nil {
							global.SHOP_LOG.Error("恢复发送者地址失败:" + err.Error())
						}
						if mcommon.IsStringInSlice(feeAddresses, AddressBytesToStr(from)) {
							// 如果打币地址在手续费热钱包地址则不处理
							continue
						}
						toAddress := AddressBytesToStr(*(rpcTx.To()))
						toAddressTxMap[toAddress] = append(toAddressTxMap[toAddress], rpcTx)
						if !mcommon.IsStringInSlice(toAddresses, toAddress) {
							toAddresses = append(toAddresses, toAddress)
						}
					}
				}
				// 从db中查询这些地址是否是冲币地址中的地址
				dbAddressRows, err := app.SQLSelectTAddressKeyColByAddress(

					toAddresses,
				)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				// 待插入数据
				var dbTxRows []*model.DBTTx
				// map[接收地址] => 产品id
				addressProductMap := make(map[string]int64)
				for _, dbAddressRow := range dbAddressRows {
					addressProductMap[dbAddressRow.Address] = dbAddressRow.Status
				}
				// 时间
				now := time.Now().Unix()
				// 遍历数据库中有交易的地址
				for _, dbAddressRow := range dbAddressRows {
					if dbAddressRow.Status < 0 {
						continue
					}
					// 获取地址对应的交易列表
					txes := toAddressTxMap[dbAddressRow.Address]
					for _, tx := range txes {

						signer := types.MakeSigner(params.MainnetChainConfig, rpcBlock.Number(), rpcBlock.Time())

						from, err := signer.Sender(tx)
						if err != nil {
							global.SHOP_LOG.Error("恢复发送者地址失败:" + err.Error())
						}
						fromAddress := AddressBytesToStr(from)
						toAddress := AddressBytesToStr(*(tx.To()))
						balanceReal, err := WeiBigIntToEthStr(tx.Value())
						if err != nil {
							mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
							return
						}
						dbTxRows = append(dbTxRows, &model.DBTTx{
							ProductID:    addressProductMap[toAddress],
							TxID:         tx.Hash().String(),
							FromAddress:  fromAddress,
							ToAddress:    toAddress,
							BalanceReal:  balanceReal,
							CreateTime:   now,
							HandleStatus: app.TxStatusInit,
							HandleMsg:    "",
							HandleTime:   now,
							OrgStatus:    app.TxOrgStatusInit,
							OrgMsg:       "",
							OrgTime:      now,
						})
					}
				}
				// 插入交易数据
				_, err = model.SQLCreateManyTTx(
					context.Background(),
					xenv.DbCon,
					dbTxRows,
					true,
				)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				// 更新检查到的最新区块数
				_, err = app.SQLUpdateTAppStatusIntByKGreater(
					context.Background(),
					xenv.DbCon,
					&model.DBTAppStatusInt{
						K: "eth_block",
						V: i,
					},
					&model.DBTAppStatusInt{
						K: "eth_block_top",
						V: rpcBlockNum,
					},
				)
				if err != nil {
					mcommon.Log.Errorf("SQLUpdateTAppStatusIntByK err: [%T] %s", err, err.Error())
					return
				}
			}
		}
	})
}

// CheckRawTxSend 发送交易
func CheckRawTxSend() {
	lockKey := "EthCheckRawTxSend"
	app.LockWrap(lockKey, func() {
		// 获取待发送的数据
		sendRows, err := app.SQLSelectTSendColByStatus(app.SendStatusInit)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 首先单独处理提币，提取提币通知要使用的数据
		var withdrawIDs []int64
		for _, sendRow := range sendRows {
			switch sendRow.RelatedType {
			case app.SendRelationTypeWithdraw:
				if !mcommon.IsIntInSlice(withdrawIDs, sendRow.RelatedId) {
					withdrawIDs = append(withdrawIDs, sendRow.RelatedId)
				}
			}
		}
		withdrawMap, err := app.SQLGetWithdrawMap(
			withdrawIDs,
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 产品
		var productIDs []int64
		for _, withdrawRow := range withdrawMap {
			if !mcommon.IsIntInSlice(productIDs, withdrawRow.ProductId) {
				productIDs = append(productIDs, withdrawRow.ProductId)
			}
		}
		productMap, err := app.SQLGetProductMap(
			context.Background(),
			xenv.DbCon,
			[]string{
				model.DBColTProductID,
				model.DBColTProductAppName,
				model.DBColTProductCbURL,
				model.DBColTProductAppSk,
			},
			productIDs,
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 执行发送
		var sendIDs []int64
		var txIDs []int64
		var erc20TxIDs []int64
		var erc20TxFeeIDs []int64
		withdrawIDs = []int64{}
		// 通知数据
		var notifyRows []*model.DBTProductNotify
		now := time.Now().Unix()
		var sendTxHashes []string
		onSendOk := func(sendRow *admodel.TSend) error {
			// 将发送成功和占位数据计入数组
			if !mcommon.IsIntInSlice(sendIDs, sendRow.Id) {
				sendIDs = append(sendIDs, sendRow.Id)
			}
			switch sendRow.RelatedType {
			case app.SendRelationTypeTx:
				if !mcommon.IsIntInSlice(txIDs, sendRow.RelatedId) {
					txIDs = append(txIDs, sendRow.RelatedId)
				}
			case app.SendRelationTypeWithdraw:
				if !mcommon.IsIntInSlice(withdrawIDs, sendRow.RelatedId) {
					withdrawIDs = append(withdrawIDs, sendRow.RelatedId)
				}
			case app.SendRelationTypeTxErc20:
				if !mcommon.IsIntInSlice(erc20TxIDs, sendRow.RelatedId) {
					erc20TxIDs = append(erc20TxIDs, sendRow.RelatedId)
				}
			case app.SendRelationTypeTxErc20Fee:
				if !mcommon.IsIntInSlice(erc20TxFeeIDs, sendRow.RelatedId) {
					erc20TxFeeIDs = append(erc20TxFeeIDs, sendRow.RelatedId)
				}
			}
			// 如果是提币，创建通知信息
			if sendRow.RelatedType == app.SendRelationTypeWithdraw {
				withdrawRow, ok := withdrawMap[sendRow.RelatedId]
				if !ok {
					mcommon.Log.Errorf("withdrawMap no: %d", sendRow.RelatedId)
					return nil
				}
				productRow, ok := productMap[withdrawRow.ProductId]
				if !ok {
					mcommon.Log.Errorf("productMap no: %d", withdrawRow.ProductId)
					return nil
				}
				nonce := mcommon.GetUUIDStr()
				reqObj := gin.H{
					"tx_hash":     sendRow.TxId,
					"balance":     withdrawRow.BalanceReal,
					"app_name":    productRow.AppName,
					"out_serial":  withdrawRow.OutSerial,
					"address":     withdrawRow.ToAddress,
					"symbol":      withdrawRow.Symbol,
					"notify_type": app.NotifyTypeWithdrawSend,
				}
				reqObj["sign"] = mcommon.WechatGetSign(productRow.AppSk, reqObj)
				req, err := json.Marshal(reqObj)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return err
				}
				notifyRows = append(notifyRows, &model.DBTProductNotify{
					Nonce:        nonce,
					ProductID:    withdrawRow.ProductId,
					ItemType:     app.SendRelationTypeWithdraw,
					ItemID:       withdrawRow.Id,
					NotifyType:   app.NotifyTypeWithdrawSend,
					TokenSymbol:  withdrawRow.Symbol,
					URL:          productRow.CbURL,
					Msg:          string(req),
					HandleStatus: app.NotifyStatusInit,
					HandleMsg:    "",
					CreateTime:   now,
					UpdateTime:   now,
				})
			}
			return nil
		}
		for _, sendRow := range sendRows {
			// 发送数据中需要排除占位数据
			if sendRow.Hex != "" {
				rawTxBytes, err := hex.DecodeString(sendRow.Hex)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					continue
				}
				tx := &types.Transaction{}
				err = tx.UnmarshalBinary(rawTxBytes)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					continue
				}
				err = ethclient.RpcSendTransaction(
					context.Background(),
					tx,
				)
				if err != nil {
					if !strings.Contains(err.Error(), "known transaction") {
						mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
						continue
					}
				}
				sendTxHashes = append(sendTxHashes, sendRow.TxId)

				err = onSendOk(sendRow)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
			} else if mcommon.IsStringInSlice(sendTxHashes, sendRow.TxId) {
				err = onSendOk(sendRow)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
			}
		}
		// 插入通知
		_, err = model.SQLCreateManyTProductNotify(
			context.Background(),
			xenv.DbCon,
			notifyRows,
			true,
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新提币状态
		_, err = app.SQLUpdateTWithdrawStatusByIDs(
			context.Background(),
			xenv.DbCon,
			withdrawIDs,
			&model.DBTWithdraw{
				HandleStatus: app.WithdrawStatusSend,
				HandleMsg:    "send",
				HandleTime:   now,
			},
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新eth零钱整理状态
		_, err = app.SQLUpdateTTxOrgStatusByIDs(
			context.Background(),
			xenv.DbCon,
			txIDs,
			model.DBTTx{
				OrgStatus: app.TxOrgStatusSend,
				OrgMsg:    "send",
				OrgTime:   now,
			},
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新erc20零钱整理状态
		_, err = app.SQLUpdateTTxErc20OrgStatusByIDs(

			global.SHOP_DB,
			erc20TxIDs,
			model.DBTTxErc20{
				OrgStatus: app.TxOrgStatusSend,
				OrgMsg:    "send",
				OrgTime:   now,
			},
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新erc20手续费状态
		_, err = app.SQLUpdateTTxErc20OrgStatusByIDs(
			global.SHOP_DB,
			erc20TxFeeIDs,
			model.DBTTxErc20{
				OrgStatus: app.TxOrgStatusFeeSend,
				OrgMsg:    "fee send",
				OrgTime:   now,
			},
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新发送状态
		_, err = app.SQLUpdateTSendStatusByIDs(
			context.Background(),
			xenv.DbCon,
			sendIDs,
			model.DBTSend{
				HandleStatus: app.SendStatusSend,
				HandleMsg:    "send",
				HandleTime:   now,
			},
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
	})
}

// CheckRawTxConfirm 确认tx是否打包完成
func CheckRawTxConfirm() {
	lockKey := "EthCheckRawTxConfirm"
	app.LockWrap(lockKey, func() {
		sendRows, err := app.SQLSelectTSendColByStatus(

			app.SendStatusSend,
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		var withdrawIDs []int64
		for _, sendRow := range sendRows {
			if sendRow.RelatedType == app.SendRelationTypeWithdraw {
				// 提币
				if !mcommon.IsIntInSlice(withdrawIDs, sendRow.RelatedId) {
					withdrawIDs = append(withdrawIDs, sendRow.RelatedId)
				}
			}
		}
		withdrawMap, err := app.SQLGetWithdrawMap(
			withdrawIDs,
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		var productIDs []int64
		for _, withdrawRow := range withdrawMap {
			if !mcommon.IsIntInSlice(productIDs, withdrawRow.ProductId) {
				productIDs = append(productIDs, withdrawRow.ProductId)
			}
		}
		productMap, err := app.SQLGetProductMap(
			context.Background(),
			xenv.DbCon,
			[]string{
				model.DBColTProductID,
				model.DBColTProductAppName,
				model.DBColTProductCbURL,
				model.DBColTProductAppSk,
			},
			productIDs,
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}

		now := time.Now().Unix()
		var notifyRows []*model.DBTProductNotify
		var sendIDs []int64
		var txIDs []int64
		var erc20TxIDs []int64
		var erc20TxFeeIDs []int64
		withdrawIDs = []int64{}
		var sendHashes []string
		for _, sendRow := range sendRows {
			if !mcommon.IsStringInSlice(sendHashes, sendRow.TxId) {
				rpcTx, err := ethclient.RpcTransactionByHash(
					context.Background(),
					sendRow.TxId,
				)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					continue
				}
				if rpcTx == nil {
					continue
				}
				sendHashes = append(sendHashes, sendRow.TxId)
			}
			if sendRow.RelatedType == app.SendRelationTypeWithdraw {
				// 提币
				withdrawRow, ok := withdrawMap[sendRow.RelatedId]
				if !ok {
					mcommon.Log.Errorf("no withdrawMap: %d", sendRow.RelatedId)
					return
				}
				productRow, ok := productMap[withdrawRow.ProductId]
				if !ok {
					mcommon.Log.Errorf("no productMap: %d", withdrawRow.ProductId)
					return
				}
				nonce := mcommon.GetUUIDStr()
				reqObj := gin.H{
					"tx_hash":     sendRow.TxId,
					"balance":     withdrawRow.BalanceReal,
					"app_name":    productRow.AppName,
					"out_serial":  withdrawRow.OutSerial,
					"address":     withdrawRow.ToAddress,
					"symbol":      withdrawRow.Symbol,
					"notify_type": app.NotifyTypeWithdrawConfirm,
				}
				reqObj["sign"] = mcommon.WechatGetSign(productRow.AppSk, reqObj)
				req, err := json.Marshal(reqObj)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				notifyRows = append(notifyRows, &model.DBTProductNotify{
					Nonce:        nonce,
					ProductID:    withdrawRow.ProductId,
					ItemType:     app.SendRelationTypeWithdraw,
					ItemID:       withdrawRow.Id,
					NotifyType:   app.NotifyTypeWithdrawConfirm,
					TokenSymbol:  withdrawRow.Symbol,
					URL:          productRow.CbURL,
					Msg:          string(req),
					HandleStatus: app.NotifyStatusInit,
					HandleMsg:    "",
					CreateTime:   now,
					UpdateTime:   now,
				})

			}
			// 将发送成功和占位数据计入数组
			if !mcommon.IsIntInSlice(sendIDs, sendRow.Id) {
				sendIDs = append(sendIDs, sendRow.Id)
			}
			switch sendRow.RelatedType {
			case app.SendRelationTypeTx:
				if !mcommon.IsIntInSlice(txIDs, sendRow.RelatedId) {
					txIDs = append(txIDs, sendRow.RelatedId)
				}
			case app.SendRelationTypeWithdraw:
				if !mcommon.IsIntInSlice(withdrawIDs, sendRow.RelatedId) {
					withdrawIDs = append(withdrawIDs, sendRow.RelatedId)
				}
			case app.SendRelationTypeTxErc20:
				if !mcommon.IsIntInSlice(erc20TxIDs, sendRow.RelatedId) {
					erc20TxIDs = append(erc20TxIDs, sendRow.RelatedId)
				}
			case app.SendRelationTypeTxErc20Fee:
				if !mcommon.IsIntInSlice(erc20TxFeeIDs, sendRow.RelatedId) {
					erc20TxFeeIDs = append(erc20TxFeeIDs, sendRow.RelatedId)
				}
			}
		}
		// 添加通知信息
		_, err = model.SQLCreateManyTProductNotify(
			context.Background(),
			xenv.DbCon,
			notifyRows,
			true,
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新提币状态
		_, err = app.SQLUpdateTWithdrawStatusByIDs(
			context.Background(),
			xenv.DbCon,
			withdrawIDs,
			&model.DBTWithdraw{
				HandleStatus: app.WithdrawStatusConfirm,
				HandleMsg:    "confirmed",
				HandleTime:   now,
			},
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新eth零钱整理状态
		_, err = app.SQLUpdateTTxOrgStatusByIDs(
			context.Background(),
			xenv.DbCon,
			txIDs,
			model.DBTTx{
				OrgStatus: app.TxOrgStatusConfirm,
				OrgMsg:    "confirmed",
				OrgTime:   now,
			},
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新erc20零钱整理状态
		_, err = app.SQLUpdateTTxErc20OrgStatusByIDs(
			global.SHOP_DB,
			erc20TxIDs,
			model.DBTTxErc20{
				OrgStatus: app.TxOrgStatusConfirm,
				OrgMsg:    "confirmed",
				OrgTime:   now,
			},
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新erc20零钱整理eth手续费状态
		_, err = app.SQLUpdateTTxErc20OrgStatusByIDs(
			global.SHOP_DB,
			erc20TxFeeIDs,
			model.DBTTxErc20{
				OrgStatus: app.TxOrgStatusFeeConfirm,
				OrgMsg:    "eth fee confirmed",
				OrgTime:   now,
			},
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 更新发送状态
		_, err = app.SQLUpdateTSendStatusByIDs(
			context.Background(),
			xenv.DbCon,
			sendIDs,
			model.DBTSend{
				HandleStatus: app.SendStatusConfirm,
				HandleMsg:    "confirmed",
				HandleTime:   now,
			},
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
	})
}

// CheckWithdraw 检测提现
func CheckWithdraw() {
	lockKey := "EthCheckWithdraw"
	app.LockWrap(lockKey, func() {
		// 获取需要处理的提币数据
		withdrawRows, err := app.SQLSelectTWithdrawColByStatus(
			context.Background(),
			xenv.DbCon,
			[]string{
				model.DBColTWithdrawID,
			},
			app.WithdrawStatusInit,
			[]string{CoinSymbol},
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if len(withdrawRows) == 0 {
			// 没有要处理的提币
			return
		}
		// 获取热钱包地址
		hotAddressValue, err := app.SQLGetTAppConfigStrValueByK(
			context.Background(),
			xenv.DbCon,
			"hot_wallet_address",
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		_, err = StrToAddressBytes(hotAddressValue)
		if err != nil {
			mcommon.Log.Errorf("eth hot address err: [%T] %s", err, err.Error())
			return
		}
		// 获取私钥
		privateKey, err := GetPkOfAddress(
			context.Background(),
			xenv.DbCon,
			hotAddressValue,
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 获取热钱包余额
		hotAddressBalance, err := ethclient.RpcBalanceAt(
			context.Background(),
			hotAddressValue,
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		pendingBalanceRealStr, err := app.SQLGetTSendPendingBalanceReal(
			context.Background(),
			xenv.DbCon,
			hotAddressValue,
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		pendingBalance, err := EthStrToWeiBigInit(pendingBalanceRealStr)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		hotAddressBalance.Sub(hotAddressBalance, pendingBalance)
		// 获取gap price
		gasPriceValue, err := app.SQLGetTAppStatusIntValueByK(
			context.Background(),
			xenv.DbCon,
			"to_user_gas_price",
		)
		if err != nil {
			mcommon.Log.Warnf("err: [%T] %s", err, err.Error())
			return
		}
		gasPrice := gasPriceValue
		tipPriceValue, err := app.SQLGetTAppStatusIntValueByK(
			context.Background(),
			xenv.DbCon,
			"to_tip_gas_price",
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		gasLimit := int64(21000)
		feeValue := gasLimit * gasPrice
		chainID, err := ethclient.RpcNetworkID(context.Background())
		if err != nil {
			mcommon.Log.Warnf("err: [%T] %s", err, err.Error())
			return
		}
		for _, withdrawRow := range withdrawRows {
			err = handleWithdraw(withdrawRow.ID, chainID, hotAddressValue, privateKey, hotAddressBalance, gasLimit, gasPrice, tipPriceValue, feeValue)
			if err != nil {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				continue
			}
		}
	})
}

func handleWithdraw(withdrawID int64, chainID int64, hotAddress string, privateKey *ecdsa.PrivateKey, hotAddressBalance *big.Int, gasLimit, gasPrice, tipPrice, feeValue int64) error {
	isComment := false
	dbTx, err := xenv.DbCon.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer func() {
		if !isComment {
			_ = dbTx.Rollback()
		}
	}()
	// 处理业务
	withdrawRow, err := app.SQLGetTWithdrawColForUpdate(
		context.Background(),
		dbTx,
		[]string{
			model.DBColTWithdrawID,
			model.DBColTWithdrawBalanceReal,
			model.DBColTWithdrawToAddress,
		},
		withdrawID,
		app.WithdrawStatusInit,
	)
	if err != nil {
		return err
	}
	if withdrawRow == nil {
		return nil
	}
	balanceBigInt, err := EthStrToWeiBigInit(withdrawRow.BalanceReal)
	if err != nil {
		return err
	}
	hotAddressBalance.Sub(hotAddressBalance, balanceBigInt)
	hotAddressBalance.Sub(hotAddressBalance, big.NewInt(feeValue))
	if hotAddressBalance.Cmp(new(big.Int)) < 0 {
		mcommon.Log.Errorf("hot balance limit")
		hotAddressBalance.Add(hotAddressBalance, balanceBigInt)
		hotAddressBalance.Add(hotAddressBalance, big.NewInt(feeValue))
		return nil
	}
	// nonce
	nonce, err := GetNonce(
		hotAddress,
	)
	if err != nil {
		return err
	}
	// 创建交易
	var data []byte
	_, err = StrToAddressBytes(withdrawRow.ToAddress)
	if err != nil {
		return err
	}
	signedTx, err := NewSignTransaction(
		nonce,
		withdrawRow.ToAddress,
		balanceBigInt,
		gasLimit,
		gasPrice,
		tipPrice,
		data,
		chainID,
		privateKey,
	)
	if err != nil {
		return err
	}
	rawTxBytes, err := signedTx.MarshalBinary()
	if err != nil {
		mcommon.Log.Warnf("MarshalBinary err: [%T] %s", err, err.Error())
		return err
	}
	rawTxHex := hex.EncodeToString(rawTxBytes)
	txHash := strings.ToLower(signedTx.Hash().Hex())
	now := time.Now().Unix()
	_, err = app.SQLUpdateTWithdrawGenTx(
		context.Background(),
		dbTx,
		&model.DBTWithdraw{
			ID:           withdrawID,
			TxHash:       txHash,
			HandleStatus: app.WithdrawStatusHex,
			HandleMsg:    "hex",
			HandleTime:   now,
		},
	)
	if err != nil {
		return err
	}
	_, err = model.SQLCreateTSend(
		context.Background(),
		dbTx,
		&model.DBTSend{
			RelatedType:  app.SendRelationTypeWithdraw,
			RelatedID:    withdrawID,
			TxID:         txHash,
			FromAddress:  hotAddress,
			ToAddress:    withdrawRow.ToAddress,
			BalanceReal:  withdrawRow.BalanceReal,
			Gas:          gasLimit,
			GasPrice:     gasPrice,
			Nonce:        nonce,
			Hex:          rawTxHex,
			HandleStatus: app.SendStatusInit,
			HandleMsg:    "",
			HandleTime:   now,
		},
		false,
	)
	if err != nil {
		return err
	}
	// 处理完成
	err = dbTx.Commit()
	if err != nil {
		return err
	}
	isComment = true
	return nil
}

// CheckTxNotify 创建eth冲币通知
func CheckTxNotify() {
	lockKey := "EthCheckTxNotify"
	app.LockWrap(lockKey, func() {
		txRows, err := app.SQLSelectTTxColByStatus(
			context.Background(),
			xenv.DbCon,
			[]string{
				model.DBColTTxID,
				model.DBColTTxProductID,
				model.DBColTTxTxID,
				model.DBColTTxToAddress,
				model.DBColTTxBalanceReal,
			},
			app.TxStatusInit,
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		var productIDs []int64
		for _, txRow := range txRows {
			if !mcommon.IsIntInSlice(productIDs, txRow.ProductID) {
				productIDs = append(productIDs, txRow.ProductID)
			}
		}
		productMap, err := app.SQLGetProductMap(
			context.Background(),
			xenv.DbCon,
			[]string{
				model.DBColTProductID,
				model.DBColTProductAppName,
				model.DBColTProductCbURL,
				model.DBColTProductAppSk,
			},
			productIDs,
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}

		var notifyTxIDs []int64
		var notifyRows []*model.DBTProductNotify
		now := time.Now().Unix()
		for _, txRow := range txRows {
			productRow, ok := productMap[txRow.ProductID]
			if !ok {
				mcommon.Log.Warnf("no productMap: %d", txRow.ProductID)
				notifyTxIDs = append(notifyTxIDs, txRow.ID)
				continue
			}
			nonce := mcommon.GetUUIDStr()
			reqObj := gin.H{
				"tx_hash":     txRow.TxID,
				"app_name":    productRow.AppName,
				"address":     txRow.ToAddress,
				"balance":     txRow.BalanceReal,
				"symbol":      CoinSymbol,
				"notify_type": app.NotifyTypeTx,
			}
			reqObj["sign"] = mcommon.WechatGetSign(productRow.AppSk, reqObj)
			req, err := json.Marshal(reqObj)
			if err != nil {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				continue
			}
			notifyRows = append(notifyRows, &model.DBTProductNotify{
				Nonce:        nonce,
				ProductID:    txRow.ProductID,
				ItemType:     app.SendRelationTypeTx,
				ItemID:       txRow.ID,
				NotifyType:   app.NotifyTypeTx,
				TokenSymbol:  CoinSymbol,
				URL:          productRow.CbURL,
				Msg:          string(req),
				HandleStatus: app.NotifyStatusInit,
				HandleMsg:    "",
				CreateTime:   now,
				UpdateTime:   now,
			})
			notifyTxIDs = append(notifyTxIDs, txRow.ID)
		}
		_, err = model.SQLCreateManyTProductNotify(
			context.Background(),
			xenv.DbCon,
			notifyRows,
			true,
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		_, err = app.SQLUpdateTTxStatusByIDs(
			context.Background(),
			xenv.DbCon,
			notifyTxIDs,
			model.DBTTx{
				HandleStatus: app.TxStatusNotify,
				HandleMsg:    "notify",
				HandleTime:   now,
			},
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
	})
}

// CheckErc20BlockSeek 检测erc20到账
func CheckErc20BlockSeek() {
	lockKey := "Erc20CheckBlockSeek"
	app.LockWrap(lockKey, func() {
		fmt.Println("Erc20CheckBlockSeek")
		// 获取配置 延迟确认数
		confirmValue, err := app.SQLGetTAppConfigIntValueByK(
			context.Background(),
			xenv.DbCon,
			"eth_block_confirm_num",
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 获取状态 当前处理完成的最新的block number
		seekValue, err := app.SQLGetTAppStatusIntValueByK(
			context.Background(),
			xenv.DbCon,
			"eth_block",
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// rpc 获取当前最新区块数
		rpcBlockNum, err := ethclient.RpcBlockNumber(context.Background())
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}

		fmt.Println(seekValue, rpcBlockNum)
		startI := seekValue + 1
		endI := rpcBlockNum - confirmValue + 1
		if startI < endI {
			// 读取abi
			type LogTransfer struct {
				From   string
				To     string
				Tokens *big.Int
			}
			contractAbi, err := abi.JSON(strings.NewReader(ethclient.EthABI))
			if err != nil {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			// 获取所有token
			var configTokenRowAddresses []string
			configTokenRowMap := make(map[string]*admodel.TAppConfigToken)
			configTokenRows, err := app.SQLSelectTAppConfigTokenColAll(
				CoinSymbol,
			)
			if err != nil {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			for _, contractRow := range configTokenRows {
				configTokenRowAddresses = append(configTokenRowAddresses, contractRow.TokenAddress)
				configTokenRowMap[contractRow.TokenAddress] = contractRow
			}
			// 遍历获取需要查询的block信息
			//25403124 usdt 0x28c6c06298d514db089934071355e5743bf21d60
			//25403124 usdc 0xee7ae85f2fe2239e27d9c1e23fffe168d63b4055
			var startI, endI = int64(25403124), int64(25403125)
			for i := startI; i < endI; i++ {
				//for i := 25396074; i < 25396074; i++ {
				//mcommon.Log.Debugf("erc20 check block: %d", i)
				if len(configTokenRowAddresses) > 0 {
					// rpc获取block信息
					logs, err := ethclient.RpcFilterLogs(
						context.Background(),
						i,
						i,
						configTokenRowAddresses,
						contractAbi.Events["Transfer"],
					)
					if err != nil {
						mcommon.Log.Warnf("err: [%T] %s", err, err.Error())
						time.Sleep(3 * time.Second)
						continue
					}
					// 接收地址列表
					var toAddresses []string
					// map[接收地址] => []交易信息

					toAddressLogMap := make(map[string][]types.Log)

					for _, log := range logs {
						toAddress := CommonHashToAddrssStringLower(log.Topics[2])

						if log.Removed {
							continue
						}

						if !mcommon.IsStringInSlice(toAddresses, toAddress) {
							toAddresses = append(toAddresses, toAddress)
						}

						toAddressLogMap[toAddress] = append(toAddressLogMap[toAddress], log)
					}

					// 从db中查询这些地址是否是冲币地址中的地址
					fmt.Println(toAddresses)
					dbAddressRows, err := app.SQLSelectTAddressKeyColByAddress(
						toAddresses,
					)
					if err != nil {
						mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
						return
					}
					// map[接收地址] => 产品id
					addressProductMap := make(map[string]int64)
					for _, dbAddressRow := range dbAddressRows {
						addressProductMap[dbAddressRow.Address] = dbAddressRow.Status
					}
					// 时间
					now := time.Now().Unix()
					// 待添加数组
					var txErc20Rowss []*admodel.Transactions
					// 遍历数据库中有交易的地址
					for _, dbAddressRow := range dbAddressRows {
						//去掉热钱包
						if dbAddressRow.Status < 0 {
							continue
						}
						// 获取地址对应的交易列表
						logs, ok := toAddressLogMap[dbAddressRow.Address]
						if !ok {
							mcommon.Log.Errorf("toAddressLogMap no: %s", dbAddressRow.Address)
							continue
						}
						for _, log := range logs {
							//fmt.Println(log.Address)
							var transferEvent LogTransfer
							err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", log.Data)
							if err != nil {
								mcommon.Log.Warnf("err: [%T] %s", err, err.Error())
								return
							}
							transferEvent.From = CommonHashToAddrssStringLower(log.Topics[1])
							transferEvent.To = CommonHashToAddrssStringLower(log.Topics[2])
							contractAddress := strings.ToLower(log.Address.Hex())
							//fmt.Println(contractAddress)
							configTokenRow, ok := configTokenRowMap[contractAddress]
							if !ok {
								mcommon.Log.Errorf("no configTokenRowMap of: %s", contractAddress)
								return
							}
							rpcTxReceipt, err := ethclient.RpcTransactionReceipt(
								context.Background(),
								log.TxHash.Hex(),
							)
							if err != nil {
								mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
								return
							}
							if rpcTxReceipt.Status <= 0 {
								continue
							}
							rpcTx, err := ethclient.RpcTransactionByHash(
								context.Background(),
								log.TxHash.Hex(),
							)
							if err != nil {
								mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
								continue
							}
							if strings.ToLower(rpcTx.To().Hex()) != contractAddress {
								fmt.Println(strings.ToLower(rpcTx.To().Hex()), contractAddress)
								mcommon.Log.Errorf("err: [%T] %s", err, "合约地址和tx的to地址不匹配")
								continue
							}
							// 检测input
							input, err := contractAbi.Pack(
								"transfer",
								CommonHashToAddrss(log.Topics[2]),
								transferEvent.Tokens,
							)
							if err != nil {
								mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
								return
							}
							if hexutil.Encode(input) != hexutil.Encode(rpcTx.Data()) {
								mcommon.Log.Errorf("err: [%T] %s", err, "input 不匹配")
								// input 不匹配
								continue
							}
							balanceReal, err := TokenWeiBigIntToEthStr(transferEvent.Tokens, configTokenRow.TokenDecimals)
							if err != nil {
								mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
								continue
							}
							// 放入待插入数组
							exist, err := app.SQLgetEXitTRAN(
								log.TxHash.Hex(),
							)
							if exist {
								continue
							}
							txErc20Rowss = append(txErc20Rowss, &admodel.Transactions{

								TxID:        log.TxHash.Hex(),
								FromAddress: transferEvent.From,
								Address:     transferEvent.To,
								Amount:      balanceReal,
								Status:      app.TxStatusInit,
								Timestamp:   now,
								BlockHeight: i,
								Fee:         "0",
								Contract:    configTokenRow.TokenAddress,
								TokenSymbol: configTokenRow.TokenSymbol,
								Type:        "receive",
								Chain:       CoinSymbol,
								TokenID:     configTokenRow.Id,
							})

						}
					}

					_, err = model.SQLCreateManyTTxErc20new(
						context.Background(),
						xenv.DbCon,
						txErc20Rowss,
						true,
					)
					if err != nil {
						mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
						return
					}
				}
				// 更新检查到的最新区块数
				_, err = app.SQLUpdateTAppStatusIntByKGreater(
					context.Background(),
					xenv.DbCon,
					&model.DBTAppStatusInt{
						K: "eth_block",
						V: i,
					},
					&model.DBTAppStatusInt{
						K: "eth_block_top",
						V: rpcBlockNum,
					},
				)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				time.Sleep(1 * time.Second)
			}
		}
	})
}

// CheckErc20TxNotify 创建erc20冲币通知
func CheckErc20TxNotify() {
	lockKey := "Erc20CheckTxNotify"
	app.LockWrap(lockKey, func() {
		txRows, err := app.SQLSelectTTxErc20ColByStatus(app.TxStatusInit)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}

		product, err := app.SQLGetTProduct()
		if product.Id == 0 {
			global.SHOP_LOG.Error("产品没有配置")
			return
		}
		secret := product.AppSk
		recharnotifyurl := product.CbUrl
		for _, value := range txRows {

			reqObj := gin.H{
				"data": gin.H{
					"tx_hash": value.TxID,
					//"type":         "recharge",
					"address":      value.Address,
					"amount":       value.Amount,
					"from_address": value.FromAddress,
					"symbol":       value.TokenSymbol,
				},
				"type": "recharge",
				"uuid": utils.UUID(),
			}

			reqObj["sign"] = mcommon.WechatGetSign(secret, reqObj)
			req, err := json.Marshal(reqObj)
			if err != nil {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			retryCount := 0

			//{"data":{"address":"0x28c6c06298d514db089934071355e5743bf21d60","amount":"4470.437556","from_address":"0xc9c49074a35296b59c93223440245f3030028bcf","symbol":"erc20_usdt","tx_hash":"0x7d80ac3cc0b3617932735b8d1034c512847670c7355c10
			//	f862eafa45a12418d4"},"sign":"F65E3459F61366FF4EBB9B5D96E647A5","type":"recharge","uuid":"9903585b-6dda-414b-947c-2444c2f536f1"}

		GotoHttpRetry:

			_, body, errs := gorequest.New().
				Post(recharnotifyurl).Set("Content-Type", "application/json").
				Send(string(req)).
				EndBytes()
			err = app.UpdateTransationhandeltimeadd(value.Id)
			if err != nil {
				global.SHOP_LOG.Error(err.Error())
				continue
			}
			if errs != nil {
				mcommon.Log.Errorf("err: [%T] %s", errs, errs)
				retryCount++
				if retryCount < 3 {
					time.Sleep(1 * time.Second)
					goto GotoHttpRetry
				}
				continue
			}
			fmt.Println(string(body))
			time.Sleep(1 * time.Second)

		}

	})
}

// CheckErc20TxOrg erc20零钱整理
func CheckErc20TxOrg() {
	lockKey := "Erc20CheckTxOrg"
	app.LockWrap(lockKey, func() {
		// 计算转账token所需的手续费
		erc20GasUseValue, err := app.SQLGetTAppConfigIntValueByK(
			context.Background(),
			xenv.DbCon,
			"erc20_gas_use",
		)
		if err != nil {
			mcommon.Log.Warnf("err: [%T] %s", err, err.Error())
			return
		}
		gasPriceValue, err := app.SQLGetTAppStatusIntValueByK(
			context.Background(),
			xenv.DbCon,
			"to_cold_gas_price",
		)
		if err != nil {
			mcommon.Log.Warnf("err: [%T] %s", err, err.Error())
			return
		}
		tipPriceValue, err := app.SQLGetTAppStatusIntValueByK(
			context.Background(),
			xenv.DbCon,
			"to_tip_gas_price",
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		erc20Fee := big.NewInt(erc20GasUseValue * gasPriceValue)
		ethGasUse := int64(21000)
		ethFee := big.NewInt(ethGasUse * gasPriceValue)
		// chainID
		chainID, err := ethclient.RpcNetworkID(context.Background())
		if err != nil {
			mcommon.Log.Warnf("err: [%T] %s", err, err.Error())
			return
		}

		// 开始事物
		isComment := false
		dbTx, err := xenv.DbCon.BeginTxx(context.Background(), nil)
		if err != nil {
			mcommon.Log.Warnf("err: [%T] %s", err, err.Error())
			return
		}
		defer func() {
			if !isComment {
				_ = dbTx.Rollback()
			}
		}()
		// 查询需要处理的交易
		txRows, err := app.SQLSelectTTxErc20ColByOrgForUpdate(
			[]int64{app.TxOrgStatusInit, app.TxOrgStatusFeeConfirm},
			[]string{"erc20_usdt", "erc20_usdc"},
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if len(txRows) <= 0 {
			return
		}
		// 整理信息
		type StOrgInfo struct {
			TxIDs        []int64
			ToAddress    string
			TokenID      int64
			TokenBalance *big.Int
		}

		var tokenIDs []int64
		for _, txRow := range txRows {
			if !mcommon.IsIntInSlice(tokenIDs, txRow.TokenID) {
				tokenIDs = append(tokenIDs, txRow.TokenID)
			}
		}
		tokenMap, err := app.SQLGetAppConfigTokenMap()
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}

		txMap := make(map[int64]*admodel.Transactions)
		// 地址eth余额
		addressEthBalanceMap := make(map[string]*big.Int)
		// 整理信息map
		orgMap := make(map[string]*StOrgInfo)
		// 整理地址
		var toAddresses []string
		for _, txRow := range txRows {
			tokenRow, ok := tokenMap[txRow.TokenID]
			if !ok {
				mcommon.Log.Errorf("no token of: %d", txRow.TokenID)
				return
			}
			// 转换为map
			txMap[txRow.Id] = txRow
			// 读取eth余额
			_, ok = addressEthBalanceMap[txRow.Address]
			if !ok {
				balance, err := ethclient.RpcBalanceAt(
					context.Background(),
					txRow.Address,
				)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				addressEthBalanceMap[txRow.Address] = balance
			}
			// 整理信息
			orgKey := fmt.Sprintf("%s-%d", txRow.Address, txRow.TokenID)
			orgInfo, ok := orgMap[orgKey]
			if !ok {
				orgInfo = &StOrgInfo{
					TokenID:      txRow.TokenID,
					ToAddress:    txRow.Address,
					TokenBalance: new(big.Int),
				}
				orgMap[orgKey] = orgInfo
			}
			orgInfo.TxIDs = append(orgInfo.TxIDs, txRow.Id)
			txBalance, err := TokenEthStrToWeiBigInit(txRow.Amount, tokenRow.TokenDecimals)
			if err != nil {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			orgInfo.TokenBalance.Add(orgInfo.TokenBalance, txBalance)
			// 待查询id
			if !mcommon.IsStringInSlice(toAddresses, txRow.Address) {
				toAddresses = append(toAddresses, txRow.Address)
			}
		}
		// 整理地址key
		addressPKMap, err := GetPKMapOfAddresses(toAddresses)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		// 需要手续费的整理信息
		now := time.Now().Unix()
		needEthFeeMap := make(map[string]*StOrgInfo)
		tran := global.SHOP_DB.Begin()
		for k, orgInfo := range orgMap {
			// 检测是否达到整理金额
			tokenRow, ok := tokenMap[orgInfo.TokenID]
			if !ok {
				mcommon.Log.Errorf("no tokenMap: %d", orgInfo.TokenID)
				continue
			}
			orgMinBalance, err := TokenEthStrToWeiBigInit(tokenRow.OrgMinBalance, tokenRow.TokenDecimals)
			if err != nil {
				mcommon.Log.Warnf("err: [%T] %s", err, err.Error())
				continue
			}
			if orgInfo.TokenBalance.Cmp(orgMinBalance) < 0 {
				mcommon.Log.Errorf("token balance < org min balance")
				continue
			}
			// 计算eth费用
			toAddress := orgInfo.ToAddress
			addressEthBalanceMap[toAddress] = addressEthBalanceMap[toAddress].Sub(addressEthBalanceMap[toAddress], erc20Fee)
			if addressEthBalanceMap[toAddress].Cmp(new(big.Int)) < 0 {
				// eth手续费不足
				// 处理添加手续费
				needEthFeeMap[k] = orgInfo
				continue
			}
			// 处理token转账
			privateKey, ok := addressPKMap[toAddress]
			if !ok {
				mcommon.Log.Errorf("addressMap no: %s", toAddress)
				continue
			}
			// 获取nonce值
			nonce, err := GetNonce(toAddress)
			if err != nil {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				continue
			}
			// 生成交易
			contractAbi, err := abi.JSON(strings.NewReader(ethclient.EthABI))
			if err != nil {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			input, err := contractAbi.Pack(
				"transfer",
				common.HexToAddress(tokenRow.ColdAddress),
				orgInfo.TokenBalance,
			)
			if err != nil {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			signedTx, err := NewSignTransaction(
				nonce,
				tokenRow.TokenAddress,
				big.NewInt(0),
				erc20GasUseValue,
				gasPriceValue,
				tipPriceValue,
				input,
				chainID,
				privateKey,
			)
			if err != nil {
				mcommon.Log.Warnf("err: [%T] %s", err, err.Error())
				continue
			}
			rawTxBytes, err := signedTx.MarshalBinary()
			if err != nil {
				mcommon.Log.Warnf("MarshalBinary err: [%T] %s", err, err.Error())
				continue
			}
			rawTxHex := hex.EncodeToString(rawTxBytes)
			txHash := strings.ToLower(signedTx.Hash().Hex())
			// 创建存入数据
			balanceReal, err := TokenWeiBigIntToEthStr(orgInfo.TokenBalance, tokenRow.TokenDecimals)
			if err != nil {
				mcommon.Log.Warnf("err: [%T] %s", err, err.Error())
				continue
			}
			// 待插入数据
			var sendRows []*admodel.TSend
			for rowIndex, txID := range orgInfo.TxIDs {
				if rowIndex == 0 {
					sendRows = append(sendRows, &admodel.TSend{
						RelatedType:  app.SendRelationTypeTxErc20,
						RelatedId:    txID,
						TokenId:      orgInfo.TokenID,
						TxId:         txHash,
						FromAddress:  toAddress,
						ToAddress:    tokenRow.ColdAddress,
						BalanceReal:  balanceReal,
						Gas:          erc20GasUseValue,
						GasPrice:     gasPriceValue,
						Nonce:        nonce,
						Hex:          rawTxHex,
						CreateTime:   now,
						HandleStatus: app.SendStatusInit,
						HandleMsg:    "",
						HandleTime:   now,
					})
				} else {
					sendRows = append(sendRows, &admodel.TSend{
						RelatedType:  app.SendRelationTypeTxErc20,
						RelatedId:    txID,
						TokenId:      orgInfo.TokenID,
						TxId:         txHash,
						FromAddress:  toAddress,
						ToAddress:    tokenRow.ColdAddress,
						BalanceReal:  "",
						Gas:          0,
						GasPrice:     0,
						Nonce:        -1,
						Hex:          "",
						CreateTime:   now,
						HandleStatus: app.SendStatusInit,
						HandleMsg:    "",
						HandleTime:   now,
					})
				}
			}
			// 插入发送队列
			_, err = model.SQLCreateManyTSend(tran, sendRows)
			if err != nil {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			// 更新整理状态
			_, err = app.SQLUpdateTTxErc20OrgStatusByIDs(
				global.SHOP_DB,
				orgInfo.TxIDs,
				model.DBTTxErc20{
					OrgStatus: app.TxOrgStatusHex,
					OrgMsg:    "hex",
					OrgTime:   now,
				},
			)
			if err != nil {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
		}
		// 生成eth转账
		if len(needEthFeeMap) > 0 {
			// 获取热钱包地址
			feeAddressValue, err := app.SQLethGetHotADDRESS(
				CoinSymbol,
			)
			if err != nil {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			_, err = StrToAddressBytes(feeAddressValue)
			if err != nil {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			// 获取私钥
			privateKey, err := GetPkOfAddress(
				context.Background(),
				dbTx,
				feeAddressValue,
			)
			if err != nil {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			feeAddressBalance, err := ethclient.RpcBalanceAt(
				context.Background(),
				feeAddressValue,
			)
			if err != nil {
				mcommon.Log.Errorf("RpcBalanceAt err: [%T] %s", err, err.Error())
				return
			}
			pendingBalanceReal, err := app.SQLGetTSendPendingBalanceReal(
				context.Background(),
				dbTx,
				feeAddressValue,
			)
			if err != nil {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			pendingBalance, err := EthStrToWeiBigInit(pendingBalanceReal)
			if err != nil {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			feeAddressBalance.Sub(feeAddressBalance, pendingBalance)
			// 生成手续费交易

			for _, orgInfo := range needEthFeeMap {
				feeAddressBalance.Sub(feeAddressBalance, ethFee)
				feeAddressBalance.Sub(feeAddressBalance, erc20Fee)
				if feeAddressBalance.Cmp(new(big.Int)) < 0 {
					mcommon.Log.Errorf("eth fee balance limit")
					return
				}
				// nonce
				nonce, err := GetNonce(
					feeAddressValue,
				)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				// 创建交易
				var data []byte
				signedTx, err := NewSignTransaction(
					nonce,
					orgInfo.ToAddress,
					erc20Fee,
					ethGasUse,
					gasPriceValue,
					tipPriceValue,
					data,
					chainID,
					privateKey,
				)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				rawTxBytes, err := signedTx.MarshalBinary()
				if err != nil {
					mcommon.Log.Warnf("MarshalBinary err: [%T] %s", err, err.Error())
					return
				}
				rawTxHex := hex.EncodeToString(rawTxBytes)
				txHash := strings.ToLower(signedTx.Hash().Hex())
				now := time.Now().Unix()
				balanceReal, err := WeiBigIntToEthStr(erc20Fee)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}

				// 待插入数据
				var sendRows []*admodel.TSend
				for rowIndex, txID := range orgInfo.TxIDs {
					if rowIndex == 0 {
						sendRows = append(sendRows, &admodel.TSend{
							RelatedType:  app.SendRelationTypeTxErc20Fee,
							RelatedId:    txID,
							TokenId:      0,
							TxId:         txHash,
							FromAddress:  feeAddressValue,
							ToAddress:    orgInfo.ToAddress,
							BalanceReal:  balanceReal,
							Gas:          ethGasUse,
							GasPrice:     gasPriceValue,
							Nonce:        nonce,
							Hex:          rawTxHex,
							CreateTime:   now,
							HandleStatus: app.SendStatusInit,
							HandleMsg:    "",
							HandleTime:   now,
						})
					} else {
						sendRows = append(sendRows, &admodel.TSend{
							RelatedType:  app.SendRelationTypeTxErc20Fee,
							RelatedId:    txID,
							TokenId:      0,
							TxId:         txHash,
							FromAddress:  feeAddressValue,
							ToAddress:    orgInfo.ToAddress,
							BalanceReal:  "",
							Gas:          0,
							GasPrice:     0,
							Nonce:        -1,
							Hex:          "",
							CreateTime:   now,
							HandleStatus: app.SendStatusInit,
							HandleMsg:    "",
							HandleTime:   now,
						})
					}
				}
				// 插入发送数据
				_, err = model.SQLCreateManyTSend(
					tran,
					sendRows,
				)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				// 更新整理状态
				_, err = app.SQLUpdateTTxErc20OrgStatusByIDs(
					tran,
					orgInfo.TxIDs,
					model.DBTTxErc20{
						OrgStatus: app.TxOrgStatusFeeHex,
						OrgMsg:    "fee hex",
						OrgTime:   now,
					},
				)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
			}
		}

		tran.Commit()
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		isComment = true
	})
}

// CheckErc20Withdraw erc20提币
func CheckErc20Withdraw() {
	lockKey := "Erc20CheckWithdraw"
	app.LockWrap(lockKey, func() {
		var tokenSymbols []string
		tokenMap := make(map[string]*admodel.TAppConfigToken)
		addressKeyMap := make(map[string]*ecdsa.PrivateKey)
		addressEthBalanceMap := make(map[string]*big.Int)
		addressTokenBalanceMap := make(map[string]*big.Int)
		tokenRows, err := app.SQLSelectTAppConfigTokenColAll(
			CoinSymbol,
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		for _, tokenRow := range tokenRows {
			tokenMap[tokenRow.TokenSymbol] = tokenRow
			if !mcommon.IsStringInSlice(tokenSymbols, tokenRow.TokenSymbol) {
				tokenSymbols = append(tokenSymbols, tokenRow.TokenSymbol)
			}
		}
		withdrawRows, err := app.SQLSelectTWithdrawColByStatus(
			context.Background(),
			xenv.DbCon,
			[]string{
				model.DBColTWithdrawID,
				model.DBColTWithdrawSymbol,
			},
			app.WithdrawStatusInit,
			tokenSymbols,
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if len(withdrawRows) == 0 {
			return
		}
		for _, tokenRow := range tokenRows {
			// 获取私钥
			_, err = StrToAddressBytes(tokenRow.HotAddress)
			if err != nil {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
			hotAddress := tokenRow.HotAddress
			_, ok := addressKeyMap[hotAddress]
			if !ok {
				// 获取私钥
				keyRow, err := app.SQLGetTAddressKeyColByAddress(
					context.Background(),
					xenv.DbCon,
					[]string{
						model.DBColTAddressKeyPwd,
					},
					hotAddress,
				)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				if keyRow == nil {
					mcommon.Log.Errorf("no key of: %s", hotAddress)
					return
				}
				key, err := mcommon.AesDecrypt(keyRow.PrivateKey, xenv.Cfg.AESKey)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				if len(key) == 0 {
					mcommon.Log.Errorf("error key of: %s", hotAddress)
					return
				}
				key = strings.TrimPrefix(key, "0x")
				privateKey, err := crypto.HexToECDSA(key)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				addressKeyMap[hotAddress] = privateKey
			}
			_, ok = addressEthBalanceMap[hotAddress]
			if !ok {
				hotAddressBalance, err := ethclient.RpcBalanceAt(
					context.Background(),
					hotAddress,
				)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				pendingBalanceReal, err := app.SQLGetTSendPendingBalanceReal(
					context.Background(),
					xenv.DbCon,
					hotAddress,
				)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				pendingBalance, err := EthStrToWeiBigInit(pendingBalanceReal)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				hotAddressBalance.Sub(hotAddressBalance, pendingBalance)
				addressEthBalanceMap[hotAddress] = hotAddressBalance
			}
			tokenBalanceKey := fmt.Sprintf("%s-%s", tokenRow.HotAddress, tokenRow.TokenSymbol)
			_, ok = addressTokenBalanceMap[tokenBalanceKey]
			if !ok {
				tokenBalance, err := ethclient.RpcTokenBalance(
					context.Background(),
					tokenRow.TokenAddress,
					tokenRow.HotAddress,
				)
				if err != nil {
					mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
					return
				}
				addressTokenBalanceMap[tokenBalanceKey] = tokenBalance
			}
		}
		// 获取gap price
		gasPriceValue, err := app.SQLGetTAppStatusIntValueByK(
			context.Background(),
			xenv.DbCon,
			"to_user_gas_price",
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		gasPrice := gasPriceValue
		tipPriceValue, err := app.SQLGetTAppStatusIntValueByK(
			context.Background(),
			xenv.DbCon,
			"to_tip_gas_price",
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		erc20GasUseValue, err := app.SQLGetTAppConfigIntValueByK(
			context.Background(),
			xenv.DbCon,
			"erc20_gas_use",
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		gasLimit := erc20GasUseValue
		// eth fee
		feeValue := big.NewInt(gasLimit * gasPrice)
		chainID, err := ethclient.RpcNetworkID(context.Background())
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		for _, withdrawRow := range withdrawRows {
			err = handleErc20Withdraw(withdrawRow.ID, chainID, &tokenMap, &addressKeyMap, &addressEthBalanceMap, &addressTokenBalanceMap, gasLimit, gasPrice, tipPriceValue, feeValue)
			if err != nil {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				continue
			}
		}
	})
}

func handleErc20Withdraw(withdrawID int64, chainID int64, tokenMap *map[string]*admodel.TAppConfigToken, addressKeyMap *map[string]*ecdsa.PrivateKey, addressEthBalanceMap *map[string]*big.Int, addressTokenBalanceMap *map[string]*big.Int, gasLimit, gasPrice, tipPrice int64, feeValue *big.Int) error {
	isComment := false
	dbTx, err := xenv.DbCon.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer func() {
		if !isComment {
			_ = dbTx.Rollback()
		}
	}()
	// 处理业务
	withdrawRow, err := app.SQLGetTWithdrawColForUpdate(
		context.Background(),
		dbTx,
		[]string{
			model.DBColTWithdrawID,
			model.DBColTWithdrawBalanceReal,
			model.DBColTWithdrawToAddress,
			model.DBColTWithdrawSymbol,
		},
		withdrawID,
		app.WithdrawStatusInit,
	)
	if err != nil {
		return err
	}
	if withdrawRow == nil {
		return nil
	}
	tokenRow, ok := (*tokenMap)[withdrawRow.Symbol]
	if !ok {
		mcommon.Log.Errorf("no tokenMap: %s", withdrawRow.Symbol)
		return nil
	}
	hotAddress := tokenRow.HotAddress
	key, ok := (*addressKeyMap)[hotAddress]
	if !ok {
		mcommon.Log.Errorf("no addressKeyMap: %s", hotAddress)
		return nil
	}
	(*addressEthBalanceMap)[hotAddress] = (*addressEthBalanceMap)[hotAddress].Sub(
		(*addressEthBalanceMap)[hotAddress],
		feeValue,
	)
	if (*addressEthBalanceMap)[hotAddress].Cmp(new(big.Int)) < 0 {
		mcommon.Log.Errorf("%s eth limit", hotAddress)
		return nil
	}
	tokenBalanceKey := fmt.Sprintf("%s-%s", tokenRow.HotAddress, tokenRow.TokenSymbol)
	tokenBalance, err := TokenEthStrToWeiBigInit(withdrawRow.BalanceReal, tokenRow.TokenDecimals)
	if err != nil {
		return err
	}
	(*addressTokenBalanceMap)[tokenBalanceKey] = (*addressTokenBalanceMap)[tokenBalanceKey].Sub(
		(*addressTokenBalanceMap)[tokenBalanceKey],
		tokenBalance,
	)
	if (*addressTokenBalanceMap)[tokenBalanceKey].Cmp(new(big.Int)) < 0 {
		mcommon.Log.Errorf("%s token limit", tokenBalanceKey)
		return nil
	}
	// 获取nonce值
	nonce, err := GetNonce(hotAddress)
	if err != nil {
		return err
	}
	// 生成交易
	contractAbi, err := abi.JSON(strings.NewReader(ethclient.EthABI))
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return err
	}
	input, err := contractAbi.Pack(
		"transfer",
		common.HexToAddress(withdrawRow.ToAddress),
		tokenBalance,
	)
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return err
	}
	signedTx, err := NewSignTransaction(
		nonce,
		tokenRow.TokenAddress,
		big.NewInt(0),
		gasLimit,
		gasPrice,
		tipPrice,
		input,
		chainID,
		key,
	)
	if err != nil {
		return err
	}
	rawTxBytes, err := signedTx.MarshalBinary()
	if err != nil {
		mcommon.Log.Warnf("MarshalBinary err: [%T] %s", err, err.Error())
		return err
	}
	rawTxHex := hex.EncodeToString(rawTxBytes)
	txHash := strings.ToLower(signedTx.Hash().Hex())
	now := time.Now().Unix()
	_, err = app.SQLUpdateTWithdrawGenTx(
		context.Background(),
		dbTx,
		&model.DBTWithdraw{
			ID:           withdrawID,
			TxHash:       txHash,
			HandleStatus: app.WithdrawStatusHex,
			HandleMsg:    "hex",
			HandleTime:   now,
		},
	)
	if err != nil {
		return err
	}
	_, err = model.SQLCreateTSend(
		context.Background(),
		dbTx,
		&model.DBTSend{
			RelatedType:  app.SendRelationTypeWithdraw,
			RelatedID:    withdrawID,
			TxID:         txHash,
			FromAddress:  hotAddress,
			ToAddress:    withdrawRow.ToAddress,
			BalanceReal:  withdrawRow.BalanceReal,
			Gas:          gasLimit,
			GasPrice:     gasPrice,
			Nonce:        nonce,
			Hex:          rawTxHex,
			HandleStatus: app.SendStatusInit,
			HandleMsg:    "",
			HandleTime:   now,
		},
		false,
	)
	if err != nil {
		return err
	}
	// 处理完成
	err = dbTx.Commit()
	if err != nil {
		return err
	}
	isComment = true
	return nil
}

// CheckGasPrice 检测gas price
func CheckGasPrice() {
	lockKey := "EthCheckGasPrice"
	app.LockWrap(lockKey, func() {
		// 获取最高单价
		maxValue, err := app.SQLGetTAppStatusIntValueByK(
			context.Background(),
			xenv.DbCon,
			"max_gas_price_eth",
		)
		if err != nil {
			if !strings.Contains(err.Error(), "no app status int of") {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
		}
		if maxValue <= 0 {
			maxValue = 80000000000
			// 创建
			_, err := model.SQLCreateTAppStatusInt(
				context.Background(),
				xenv.DbCon,
				&model.DBTAppStatusInt{
					K: "max_gas_price_eth",
					V: maxValue,
				},
				true,
			)
			if err != nil {
				mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
				return
			}
		}
		type StRespGasPrice struct {
			Status  string `json:"status"`
			Message string `json:"message"`
			Result  struct {
				LastBlock       int64   `json:"LastBlock,string"`
				SafeGasPrice    float64 `json:"SafeGasPrice,string"`
				ProposeGasPrice float64 `json:"ProposeGasPrice,string"`
				FastGasPrice    float64 `json:"FastGasPrice,string"`
				SuggestBaseFee  float64 `json:"suggestBaseFee,string"`
				GasUsedRatio    string  `json:"gasUsedRatio"`
			} `json:"result"`
		}
		gresp, body, errs := gorequest.New().
			Proxy(xenv.Cfg.Proxy).
			Get(fmt.Sprintf("https://api.etherscan.io/v2/api?chainid=1&module=gastracker&action=gasoracle&apikey=%s", global.SHOP_CONFIG.System.EtherScanapi)).
			Timeout(time.Second * 120).
			End()
		if errs != nil {
			mcommon.Log.Errorf("err: [%T] %s", errs[0], errs[0].Error())
			return
		}
		if gresp.StatusCode != http.StatusOK {
			// 状态错误
			mcommon.Log.Errorf("req status error: %d", gresp.StatusCode)
			return
		}
		var resp StRespGasPrice
		err = json.Unmarshal([]byte(body), &resp)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if resp.Status != "1" {
			// 状态错误
			mcommon.Log.Errorf("req status error: %s", resp.Status)
			return
		}
		toUserGasPrice := int64(2 * resp.Result.FastGasPrice * math.Pow10(9))
		toColdGasPrice := int64(1.2 * resp.Result.FastGasPrice * math.Pow10(9))
		tipFee := int64(math.Ceil(resp.Result.FastGasPrice-resp.Result.SuggestBaseFee) * math.Pow10(9))
		if tipFee < 0 {
			tipFee = 1 * int64(math.Pow10(9))
		}
		if toUserGasPrice > maxValue {
			toUserGasPrice = maxValue
		}
		if toColdGasPrice > maxValue {
			toColdGasPrice = maxValue
		}
		if tipFee > toColdGasPrice {
			tipFee = toColdGasPrice
		}
		_, err = app.SQLUpdateTAppStatusIntByK(
			context.Background(),
			xenv.DbCon,
			&model.DBTAppStatusInt{
				K: "to_user_gas_price",
				V: toUserGasPrice,
			},
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		_, err = app.SQLUpdateTAppStatusIntByK(
			context.Background(),
			xenv.DbCon,
			&model.DBTAppStatusInt{
				K: "to_cold_gas_price",
				V: toColdGasPrice,
			},
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		_, err = app.SQLUpdateTAppStatusIntByK(
			context.Background(),
			xenv.DbCon,
			&model.DBTAppStatusInt{
				K: "to_tip_gas_price",
				V: tipFee,
			},
		)
		if err != nil {
			mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
	})

}
func Gettokenbanance() {
	tokenBalance, err := ethclient.RpcTokenBalance(
		context.Background(),
		"0xdac17f958d2ee523a2206206994597c13d831ec7",
		"0xee7ae85f2fe2239e27d9c1e23fffe168d63b4055",
	)
	if err != nil {
		mcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	fmt.Println(tokenBalance)

}
