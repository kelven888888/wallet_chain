package model

import "time"

type WalletAddress struct {
	Id         int64
	CreateTime time.Time
	UpdateTime *time.Time
	Address    string
	WalletType string
	Status     int
	Remarks    string
}

func (WalletAddress) TableName() string {
	return "wallet_address"
}
