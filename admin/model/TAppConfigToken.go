package model

type TAppConfigToken struct {
	Id            int64  `json:"id" gorm:"column:id;primaryKey;not null;type:int(11)"`
	TokenAddress  string `json:"token_address" gorm:"column:token_address;not null;type:varchar(128)"`
	TokenDecimals int64  `json:"token_decimals" gorm:"column:token_decimals;not null;type:int(11)"`
	TokenSymbol   string `json:"token_symbol" gorm:"column:token_symbol;not null;type:varchar(128)"`
	ColdAddress   string `json:"cold_address" gorm:"column:cold_address;not null;type:varchar(128)"`
	HotAddress    string `json:"hot_address" gorm:"column:hot_address;not null;type:varchar(128)"`
	OrgMinBalance string `json:"org_min_balance" gorm:"column:org_min_balance;not null;type:varchar(128)"`
	CreateTime    int64  `json:"create_time" gorm:"column:create_time;not null;type:bigint(20)"`
	Chain         string `json:"chain" gorm:"column:chain;not null;type:varchar(255)"`
}

func (*TAppConfigToken) TableName() string {
	return "t_app_config_token"
}
