package model

type TSend struct {
	Id           int64  `json:"id" gorm:"column:id;primaryKey;not null;type:int(11)"`
	RelatedType  int8   `json:"related_type" gorm:"column:related_type;not null;type:tinyint(4)"`
	RelatedId    int64  `json:"related_id" gorm:"column:related_id;not null;type:int(11)"`
	TokenId      int64  `json:"token_id" gorm:"column:token_id;not null;type:int(11)"`
	TxId         string `json:"tx_id" gorm:"column:tx_id;not null;type:varchar(128)"`
	FromAddress  string `json:"from_address" gorm:"column:from_address;not null;type:varchar(128)"`
	ToAddress    string `json:"to_address" gorm:"column:to_address;not null;type:varchar(128)"`
	BalanceReal  string `json:"balance_real" gorm:"column:balance_real;not null;type:varchar(128)"`
	Gas          int64  `json:"gas" gorm:"column:gas;not null;type:bigint(20)"`
	GasPrice     int64  `json:"gas_price" gorm:"column:gas_price;not null;type:bigint(20)"`
	Nonce        int64  `json:"nonce" gorm:"column:nonce;not null;type:int(11)"`
	Hex          string `json:"hex" gorm:"column:hex;not null;type:varchar(2048)"`
	CreateTime   int64  `json:"create_time" gorm:"column:create_time;not null;type:bigint(20)"`
	HandleStatus int8   `json:"handle_status" gorm:"column:handle_status;not null;type:tinyint(4)"`
	HandleMsg    string `json:"handle_msg" gorm:"column:handle_msg;not null;type:varchar(1024)"`
	HandleTime   int64  `json:"handle_time" gorm:"column:handle_time;not null;type:bigint(20)"`
	Chain        string `json:"chain" gorm:"column:chain;not null;type:varchar(1024)"`
}

func (*TSend) TableName() string {
	return "t_send"
}
