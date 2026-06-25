package model

// OtherParam .
type OtherParam struct {
	Key   string `gorm:"'key' unique"`
	Value string `gorm:"'value'"`
}

// TableName 表名
func (fh OtherParam) TableName() string {
	return "param"
}

// Account 账户分配的地址
type Account struct {
	Model
	Address    string `gorm:"'address' unique DEFAULT '' "`         // 唯一索引
	PublicKey  string `gorm:"'public_key'"  json:"publickey"`       // 公钥新版字段 如果有就是新版
	PrivateKey string `gorm:"'private_key'" sql:"comment:'地址私钥'"`   // 地址私钥
	Index      int    `gorm:"'index' DEFAULT 0" sql:"comment:'位置'"` // 唯一
	User       string `gorm:"'user'"`
	Ctime      int64  `gorm:"'ctime'"`                          // 创建时间
	Amount     int64  `gorm:"'amount' index INTEGER DEFAULT 0"` // 主链币种余额
}

func (fh Account) TableName() string {
	return "account"
}

// Balance  代币余额
type Balance struct {
	Model
	Address  string `gorm:"'address' DEFAULT '' "`
	Contract string `gorm:"'contract' index"` // 哪种合约
	Amount   int64  `gorm:"'amount' index INTEGER DEFAULT 0"`
}

func (fh Balance) TableName() string {
	return "balance"
}

// Transactions .
type Transactions struct {
	Model
	TxID        string `gorm:"'tx_id'" json:"txid"`
	BlockHeight int64  `gorm:"'block_height'" json:"blockheight"`
	PublicKey   string `gorm:"'public_key'"  json:"publickey"` // 公钥新版字段 如果有就是新版
	Address     string `gorm:"'address' index" json:""`
	FromAddress string `gorm:"'from_address'" json:"fromaddress"`
	Contract    string `gorm:"'contract' index"` // 哪种合约
	Amount      string `gorm:"'amount'" json:"amount"`
	Fee         string `gorm:"'fee'" json:"fee"` // 保留字段
	Timestamp   int64  `gorm:"'timestamp'"`
	Type        string `gorm:"'type'"` // send recive collect
}

func (fh Transactions) TableName() string {
	return "transactions"
}
