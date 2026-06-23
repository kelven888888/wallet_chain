package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type UsdtWithdrawModel struct {
	Id            uint      `json:"id" form:"id" `
	CreateTime    time.Time `json:"-" `
	UpdateTime    time.Time `json:"-" `
	Remarks       string    `json:"-" `
	Username      string
	WalletPath    string          ` form:"wallet_path" `
	Amount        decimal.Decimal ` form:"amount" `
	Status        int
	PathType      string          ` form:"path_type" `
	Type          int             `comment:"1虚拟货币2银行卡" form:"type"`
	Msg           string          `json:"msg" `
	Hash          string          `json:"hash"`
	Fee           decimal.Decimal `json:"fee"`
	TradePassword string          ` form:"trade_password"  gorm:"-"`
	ExchangeRate  decimal.Decimal

	//ModelTime
}

func (*UsdtWithdrawModel) TableName() string {
	return "account_user_withdraw"
}
