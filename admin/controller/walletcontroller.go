package controller

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"wallet_chain.com/admin/model"
	"wallet_chain.com/admin/service"
	"wallet_chain.com/crondtab/wallet_api/util"
	"wallet_chain.com/global"
	"wallet_chain.com/utils"
)

type WalletCtr struct {
	Services service.WalletPath
	BaseController
}
type Rechargresp struct {
	Pid         int64  `json:"pid"  mapstructure:"pid"`
	Cid         int64  `json:"cid"  mapstructure:"cid"`
	ChainID     string `json:"chain_id" mapstructure:"chain_id"`
	TokenID     string `json:"token_id" mapstructure:"token_id"`
	Currency    string `json:"currency" mapstructure:"currency"`
	Amount      string `json:"amount" `
	Address     string `json:"address"`
	Status      string `json:"status"`
	Txid        string `json:"txid"`
	BlockHeight string `json:"block_height" mapstructure:"block_height"`
	BlockTime   string `json:"block_time" mapstructure:"block_time"`
	Nonce       string `json:"nonce"`
	Timestamp   int64  `json:"timestamp"`
	Sign        string `json:"sign"`
}

type Withresp struct {
	Pid          int64  `json:"pid"`
	Cid          int64  `json:"cid"`
	Address      string `json:"address"`
	ChainID      string `json:"chain_id" mapstructure:"chain_id"`
	TokenID      string `json:"token_id" mapstructure:"token_id"`
	Currency     string `json:"currency"`
	Amount       string `json:"amount"`
	ThirdPartyID string `json:"third_party_id" mapstructure:"third_party_id"`
	Remark       string `json:"remark"`
	Status       int64  `json:"status"`
	Txid         string `json:"txid"`
	BlockHeight  string `json:"block_height" mapstructure:"block_height"`
	BlockTime    string `json:"block_time" mapstructure:"block_time"`
	Nonce        string `json:"nonce"`
	Timestamp    int64  `json:"timestamp"`
	Sign         string `json:"sign"`
}

func (this *WalletCtr) PayoutBack(ctx *gin.Context) {
	// TODO: 处理提现回调
	// 1. 验证签名
	// 2. 处理转账
	// 3. 记录转账信息
	// 4. ��除手续费
	// 5. 记录手续费
	// 6. 推送消息
	//验签
	var resp Withresp
	if err := ctx.ShouldBindJSON(&resp); err != nil {
		global.SHOP_LOG.Log(0, err.Error())
		return
	}
	data := map[string]interface{}{}
	for k, v := range utils.StructToMap(resp) {
		data[strings.ToLower(k)] = v
	}
	fmt.Println("PayoutBack", data)
	signs, _ := util.DoSign(data, global.SHOP_CONFIG.Wallet.Appkey)
	fmt.Println(signs)
	sign, err := util.VerifySign(data, global.SHOP_CONFIG.Wallet.Appkey, resp.Sign)
	if err != nil || !sign {
		global.SHOP_LOG.Log(0, fmt.Sprintf("充值回调验签失败%s", err.Error()))
		this.Success(ctx, "充值回调验签失败")
		return
	}
	var withresult model.UsdtWithdrawModel
	err = global.SHOP_DB.Where("wallet_path=? and id=?", resp.Address, resp.ThirdPartyID).Find(&withresult).Error
	if err != nil {
		global.SHOP_LOG.Log(0, err.Error())
		this.Error(ctx, err.Error())
		return
	}
	if withresult.Id == 0 {
		this.Error(ctx, "没有记录")
		return
	} else {
		withresult.Status = 1
		withresult.Hash = resp.Txid
		global.SHOP_DB.Updates(&withresult)
	}
	this.Success(ctx, "success")

}
