package model

type TProduct struct {
	Id          int32  `json:"id" gorm:"column:id;primaryKey;not null;type:int(11)"`
	AppName     string `json:"app_name" gorm:"column:app_name;not null;type:varchar(128)"`
	AppSk       string `json:"app_sk" gorm:"column:app_sk;not null;type:varchar(64)"`
	CbUrl       string `json:"cb_url" gorm:"column:cb_url;not null;type:varchar(512)"`
	WhitelistIp string `json:"whitelist_ip" gorm:"column:whitelist_ip;type:longtext"`
}

func (*TProduct) TableName() string {
	return "t_product"
}
