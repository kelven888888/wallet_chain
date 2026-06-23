package crondtab

import (
	"fmt"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"wallet_chain.com/admin/model"
	sdk "wallet_chain.com/crondtab/wallet_api"
	"wallet_chain.com/global"
	"wallet_chain.com/utils"
)

type WalletServer struct {
}

func (this *WalletServer) CrondCreateAdd() {
	wallettype := map[string]string{
		"BTC":    "0",
		"ETH":    "60",
		"TRON":   "195",
		"Solana": "1000",
		"Base":   "66",
	}
	for wallet_type, val := range wallettype {
		for {
			var wall_add []model.WalletAddress
			err := global.SHOP_DB.Where("wallet_type=? and status=0", wallet_type).Find(&wall_add).Error
			if err != nil {
				global.SHOP_LOG.Log(0, err.Error())
				continue

			}
			limitavaaddr, _ := strconv.Atoi(global.SHOP_CONFIG.Wallet.LimitAvaAddr)
			if len(wall_add) < limitavaaddr {
				c := sdk.NewClient(global.SHOP_CONFIG.Wallet.Url, global.SHOP_CONFIG.Wallet.Appkey, global.SHOP_CONFIG.Wallet.Pid)
				addcallbacurl := fmt.Sprintf("%swallet/rechargecallbackabcd123", global.SHOP_CONFIG.Wallet.CallbackURL)

				addr, err := c.AddressCreate(val, addcallbacurl, wallet_type)
				if err != nil {
					errmsg := fmt.Sprintf("钱包服务访问异常err:%v", err.Error())
					global.SHOP_LOG.Log(0, errmsg)
					break

				}
				if addr.Code != "00000" {
					errmsg := fmt.Sprintf("err:%v", addr.Msg)
					global.SHOP_LOG.Log(0, errmsg)
					break
				}
				if len(addr.Data.Address) == 0 {
					global.SHOP_LOG.Log(0, fmt.Sprintf("addr is empty"))
					continue
				}
				var wall_addnew model.WalletAddress

				wall_addnew.CreateTime = time.Now()
				wall_addnew.WalletType = wallet_type
				wall_addnew.Status = 0
				wall_addnew.Address = addr.Data.Address
				fmt.Println(wall_addnew)
				err = global.SHOP_DB.Save(&wall_addnew).Error
				if err != nil {
					global.SHOP_LOG.Log(0, err.Error())
					continue
				}

			} else {
				break
			}
			time.Sleep(time.Second * 3)
		}
	}

}
func (this *WalletServer) CrondWithPassToUdun() {
	var wall_add []model.UsdtWithdrawModel
	err := global.SHOP_DB.Where("status=3").Find(&wall_add).Error
	if err != nil {
		global.SHOP_LOG.Log(0, err.Error())
		return

	}
	currentcymap := map[string]string{
		//"BTC":  "0",
		"ETH":                  "60@60",
		"BTC":                  "0@0",
		"USDT Ethereum(ERC20)": "60@0xdac17f958d2ee523a2206206994597c13d831ec7",
		"USDC Ethereum(ERC20)": "60@0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		"USDT TRX(TRC20)":      "195@TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
		"USDC TRX(TRC20)":      "195@TEkxiTehnzSmSe2XqrBj4w32RUN966rdz8",
		"USDT Solana":          "1000@Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB",
		"USDC Solana":          "1000@EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		"USDC Base":            "66@0x833589fcd6edb6e08f4c7c32d4f71b54bda02913",
	}
	for _, v := range wall_add {

		address := v.WalletPath
		amounts := ""
		addcallbacurl := fmt.Sprintf("%swallet/payoutcalback", global.SHOP_CONFIG.Wallet.CallbackURL)

		wallet_type := v.PathType
		if wallet_type == "BTC" || wallet_type == "ETH" {
			price, err := utils.Get_crypto_current_price(wallet_type)
			if err != nil {
				global.SHOP_LOG.Log(0, err.Error())
				continue
			}
			d := decimal.NewFromFloat(global.SHOP_CONFIG.Wallet.Withdrawfee)
			Ramount := v.Amount.Sub(d)
			if Ramount.LessThan(decimal.NewFromInt(0)) {
				global.SHOP_LOG.Log(0, "提币减去手续费小于零")
				continue
			}
			amount := Ramount.Div(utils.Float64ToDecimal(price))
			amountss, err := strconv.ParseFloat(fmt.Sprintf("%.5f", utils.DecimalToFloat(amount)), 64)
			if err != nil {
				global.SHOP_LOG.Log(0, err.Error())
				continue
			}

			amounts = fmt.Sprint(amountss)
		} else {
			d := decimal.NewFromFloat(float64(global.SHOP_CONFIG.Wallet.Withdrawfee))
			if v.Amount.Sub(d).LessThan(decimal.NewFromInt(0)) {
				global.SHOP_LOG.Log(0, "提币小于零")
				continue
			}
			amounts = fmt.Sprint(v.Amount.Sub(d))
		}

		currentcy, ok := currentcymap[wallet_type]
		if !ok {

			global.SHOP_LOG.Log(0, "不支持币种")
			continue
		}
		if wallet_type == "ETH" || wallet_type == "USDT Ethereum(ERC20)" || wallet_type == "USDC Ethereum(ERC20)" {
			if !utils.IsEthAddress(address) {
				global.SHOP_LOG.Log(0, "地址错误")
				v.Msg = "地址错误"
				v.Status = 4
				err = global.SHOP_DB.Model(&v).Updates(&v).Error
				if err != nil {
					global.SHOP_LOG.Log(0, err.Error())
					continue
				}
				continue
			}
		}
		if wallet_type == "BTC" {
			if !utils.IsBTCAddress(address) {
				global.SHOP_LOG.Log(0, "地址错误")
				v.Msg = "地址错误"
				v.Status = 4
				err = global.SHOP_DB.Model(&v).Updates(&v).Error
				if err != nil {
					global.SHOP_LOG.Log(0, err.Error())
					continue
				}
				continue
			}
		}
		if wallet_type == "USDT TRX(TRC20)" || wallet_type == "USDC TRX(TRC20)" {
			if !utils.IsTronAddr(address) {
				global.SHOP_LOG.Log(0, "地址错误")
				v.Msg = "地址错误"
				v.Status = 4
				err = global.SHOP_DB.Model(&v).Updates(&v).Error
				if err != nil {
					global.SHOP_LOG.Log(0, err.Error())
					continue
				}
				continue
			}
		}
		if wallet_type == "USDT Solana" || wallet_type == "USDC Solana" {
			if !utils.IsSolana(address) {
				global.SHOP_LOG.Log(0, "地址错误")
				v.Msg = "地址错误"
				v.Status = 4
				err = global.SHOP_DB.Model(&v).Updates(&v).Error
				if err != nil {
					global.SHOP_LOG.Log(0, err.Error())
					continue
				}
				continue
			}
		}
		if wallet_type == "USDC Base" {
			if !utils.IsBaseAddress(address) {
				global.SHOP_LOG.Log(0, "地址错误")
				v.Msg = "地址错误"
				v.Status = 4
				err = global.SHOP_DB.Model(&v).Updates(&v).Error
				if err != nil {
					global.SHOP_LOG.Log(0, err.Error())
					continue
				}
				continue
			}
		}
		//["USDT TRX(TRC20)", "USDT Ethereum(ERC20)","USDC Ethereum(ERC20)","USDC TRX(TRC20)","BTC","ETH"]
		c := sdk.NewClient(global.SHOP_CONFIG.Wallet.Url, global.SHOP_CONFIG.Wallet.Appkey, global.SHOP_CONFIG.Wallet.Pid)
		result, err := c.Payout(address, currentcy, amounts, fmt.Sprintf("payout%d", v.Id), addcallbacurl, "payout")
		if err != nil {
			errmsg := fmt.Sprintf("err:%v", err.Error())
			global.SHOP_LOG.Log(0, errmsg)
			global.SHOP_LOG.Log(0, errmsg)
			v.Msg = err.Error()
			v.Status = 4
			err = global.SHOP_DB.Model(&v).Updates(&v).Error
			if err != nil {
				global.SHOP_LOG.Log(0, err.Error())
				continue
			}
			continue
		} else {
			if result.Code != "00000" {
				errmsg := fmt.Sprintf("err:%v", result.Msg)
				global.SHOP_LOG.Log(0, errmsg)
				v.Msg = result.Msg
				v.Status = 4
				err = global.SHOP_DB.Model(&v).Updates(&v).Error
				if err != nil {
					global.SHOP_LOG.Log(0, err.Error())
					continue
				}
				continue
			} else {
				v.Status = 5
				v.Msg = "success"
				err = global.SHOP_DB.Model(&v).Updates(&v).Error
				if err != nil {
					global.SHOP_LOG.Log(0, err.Error())
					continue
				}

			}

		}

	}

}
