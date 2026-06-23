package model

import "time"

type WalletPath struct {
	Id         uint      `json:"id" form:"id" `
	CreateTime time.Time `json:"-"`
	UpdateTime time.Time `json:"-"`
	Remarks    string
	Username   string
	WalletPath string `json:"wallet_path" form:"wallet_path" `
	WalletType string `json:"wallet_type" form:"wallet_type" `
	PathType   string `json:"path_type" form:"path_type" `

	//ModelTime
}

func (*WalletPath) TableName() string {
	return "account_wallet_path"
}
