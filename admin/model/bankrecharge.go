package model

import "time"

type BankRecharge struct {
	Id          uint `json:"id" form:"id" `
	CreateTime  time.Time
	UpdateTime  *time.Time
	Remarks     string `json:"remarks" form:"remarks"  `
	Name        string `json:"name" form:"name"  `
	Bank        string `json:"bank" form:"bank"  `
	Addr        string `json:"addr" form:"addr"  `
	Aba         string `json:"aba" form:"aba"  `
	Swift       string `json:"swift" form:"swift"  `
	Account     string `json:"account" form:"account"  `
	Ba          string `json:"ba" form:"ba"  `
	Use         *int   `json:"use" form:"use"  `
	AccountName string `json:"accountname" gorm:"column:accountname;"  `
	Ssn         string `json:"ssn" form:"ssn"  `
	Code        string `json:"code" form:"code"  `

	//ModelTime
}

func (*BankRecharge) TableName() string {
	return "account_bank_recharge"
}
