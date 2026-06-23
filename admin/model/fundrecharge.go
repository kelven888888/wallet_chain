package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type FundRecharge struct {
	Id uint `json:"id" form:"id" `

	CreateTime time.Time  `json:"-" `
	UpdateTime *time.Time `json:"-" `
	Remarks    string
	Username   string
	Amount     decimal.Decimal `form:"amount"`
	Address    string
	PathType   string `form:"path_type"`
	//WalletPath string `form:"wallet_path"`
	Type         int `comment:"1虚拟货币2银行卡" form:"type"`
	Cms          float64
	LevelCode    string
	AgCode       string
	Hash         string `comment:"hash"`
	Status       int
	Pic          string `form:"pic"`
	ExchangeRate decimal.Decimal

	//ModelTime
}

func (*FundRecharge) TableName() string {
	return "account_funds_recharge_log"
}
