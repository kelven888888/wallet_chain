package model

type TWithdraw struct {
	Id           int64  `json:"id" gorm:"column:id;primaryKey;not null;type:int(11)"`
	ProductId    int64  `json:"product_id" gorm:"column:product_id;not null;type:int(11)"`
	OutSerial    string `json:"out_serial" gorm:"column:out_serial;not null;type:varchar(64)"`
	ToAddress    string `json:"to_address" gorm:"column:to_address;not null;type:varchar(128)"`
	Memo         string `json:"memo" gorm:"column:memo;not null;type:varchar(256)"`
	Symbol       string `json:"symbol" gorm:"column:symbol;not null;type:varchar(128)"`
	BalanceReal  string `json:"balance_real" gorm:"column:balance_real;not null;type:varchar(128)"`
	TxHash       string `json:"tx_hash" gorm:"column:tx_hash;not null;type:varchar(128)"`
	CreateTime   int64  `json:"create_time" gorm:"column:create_time;not null;type:bigint(20)"`
	HandleStatus int64  `json:"handle_status" gorm:"column:handle_status;not null;type:int(11)"`
	HandleMsg    string `json:"handle_msg" gorm:"column:handle_msg;not null;type:varchar(128)"`
	HandleTime   int64  `json:"handle_time" gorm:"column:handle_time;not null;type:bigint(20)"`
}

func (*TWithdraw) TableName() string {
	return "t_withdraw"
}
