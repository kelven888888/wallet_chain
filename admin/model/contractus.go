package model

import "time"

type ContractUs struct {
	Id         int64      `json:"id" gorm:"column:id;primaryKey;not null;autoIncrement;"`
	CreateTime *time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime"`
	UpdateTime *time.Time `json:"update_time" gorm:"column:update_time;autoUpdateTime"`
	Remarks    string     `json:"remarks" gorm:"column:remarks;"`
	Business   string     `json:"business" gorm:"column:business;"`
	Reason     string     `json:"reason" gorm:"column:reason;"`
	FirstName  string     `json:"first_name" gorm:"column:first_name;"`
	LastName   string     `json:"last_name" gorm:"column:last_name;"`
	Email      string     `json:"email" gorm:"column:email;"`
	Phone      string     `json:"phone" gorm:"column:phone;"`
	Company    string     `json:"company" gorm:"column:company;"`
	Message    string     `json:"message" gorm:"column:message;"`
}

func (ContractUs) TableName() string {
	return "pc_contact_us"
}
