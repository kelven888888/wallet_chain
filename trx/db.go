package trx

import (
	"errors"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"wallet_chain.com/global"
	"wallet_chain.com/utils/Paginate"
)

const (
	Send         = "send"          // 提币 即 主地址提币到其他地址
	Receive      = "receive"       // 本平台地址 分配的用户地址收到的
	ReceiveOther = "receive_other" // 主地址收到的外来地址转帐
	Collect      = "collect"       // 本平台地址归集到主地址
	CollectOwn   = "collect_own"   // 站内转账
	CollectSend  = "collect_send"  // 本平台地址提币到站外 异常的
)

// OtherParam .
type OtherParam struct {
	Key   string `xorm:"'key' unique"`
	Value string `xorm:"'value'"`
}

// TableName 表名
func (fh OtherParam) TableName() string {
	return "param"
}

// Account 账户分配的地址
type Account struct {
	ID         int64  `xorm:"'id' pk autoincr"`
	Address    string `xorm:"'address' unique DEFAULT '' "`         // 唯一索引
	PublicKey  string `xorm:"'public_key'"  json:"publickey"`       // 公钥新版字段 如果有就是新版
	PrivateKey string `xorm:"'private_key'" sql:"comment:'地址私钥'"`   // 地址私钥
	Index      int    `xorm:"'index' DEFAULT 0" sql:"comment:'位置'"` // 唯一
	User       string `xorm:"'user'"`
	Ctime      int64  `xorm:"'ctime'"`                          // 创建时间
	Amount     int64  `xorm:"'amount' index INTEGER DEFAULT 0"` // 主链币种余额
}

func (fh Account) TableName() string {
	return "account"
}

// Balance  代币余额
type Balance struct {
	ID       int64  `xorm:"'id' pk autoincr"`
	Address  string `xorm:"'address' DEFAULT '' "`
	Contract string `xorm:"'contract' index"` // 哪种合约
	Amount   int64  `xorm:"'amount' index INTEGER DEFAULT 0"`
}

func (fh Balance) TableName() string {
	return "balance"
}

// Transactions .
type Transactions struct {
	ID          int64  `xorm:"'id' pk autoincr" json:"-"`
	TxID        string `xorm:"'tx_id'" json:"txid"`
	BlockHeight int64  `xorm:"'block_height'" json:"blockheight"`
	PublicKey   string `xorm:"'public_key'"  json:"publickey"` // 公钥新版字段 如果有就是新版
	Address     string `xorm:"'address' index" json:""`
	FromAddress string `xorm:"'from_address'" json:"fromaddress"`
	Contract    string `xorm:"'contract' index"` // 哪种合约
	Amount      string `xorm:"'amount'" json:"amount"`
	Fee         string `xorm:"'fee'" json:"fee"` // 保留字段
	Timestamp   int64  `xorm:"'timestamp'"`
	Type        string `xorm:"'type'"` // send recive collect
}

// DB .
type DB struct {
}

// NewDB 初始化数据库
func NewDB() (*gorm.DB, error) {
	if global.SHOP_DB != nil {
		return global.SHOP_DB, nil
	}
	return global.SHOP_DB, nil
}
func GetProjectRoot() (string, error) {
	// 获取调用此函数的文件的绝对路径
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", os.ErrNotExist
	}

	// 从当前文件目录开始向上查找 go.mod
	dir := filepath.Dir(filename)
	for {
		// 检查当前目录是否有 go.mod
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		// 向上移动一级
		parent := filepath.Dir(dir)
		if parent == dir {
			// 已到达文件系统根目录仍未找到
			break
		}
		dir = parent
	}

	return "", os.ErrNotExist
}

// Sync 同步数据库结构
func (db *DB) Sync() error {
	return nil
}

// InsertAccount 插入数据
func (db *DB) InsertAccount(account *Account) (int64, error) {
	//return db.("address", "private_key", "public_key", "index", "user", "ctime", "amount").Insert(account)
	err := global.SHOP_DB.Model(Account{}).Save(account).Error
	if err != nil {
		global.SHOP_LOG.Error(err.Error())
		return 0, err
	}
	return account.ID, nil
}

// UpdateAccount 更新数据
func (db *DB) UpdateAccount(account *Account) (int64, error) {
	err := global.SHOP_DB.Model(Account{}).Where("address = ? ", account.Address).Save(account).Error
	if err != nil {
		global.SHOP_LOG.Error(err.Error())
		return 0, err
	}
	return account.ID, nil
	//return global.SHOP_DB.Where("address = ? ", account.Address).Cols("amount").Update(account)
}

// GetAccountWithAddr 搜索地址是否存在
func (db *DB) GetAccountWithAddr(addr string) (*Account, error) {
	var tmp Account
	global.SHOP_DB.Where("address = ?", addr).Limit(1).Find(&tmp)

	return &tmp, nil
}

func (db *DB) GetAccountMaxIndex() int {
	var resp map[string]int
	global.SHOP_DB.Table("account").Select("IFNULL(max(`index`),0) as maxid").Find(&resp)
	return resp["maxid"]
}

// GetAccount 获取所有账户
func (db *DB) GetAccount(from int) ([]Account, error) {
	var tmp = make([]Account, 0)
	if from < 0 {
		from = 0
	}
	err := global.SHOP_DB.Find(&tmp).Error
	return tmp, err
}

// GetAccountWithBalance 获取大于minAmount的所有账户
func (db *DB) GetAccountWithBalance(startid int64, count int) ([]Account, error) {
	var tmp = make([]Account, 0)
	err := global.SHOP_DB.Where("id> ?", startid).Limit(count).Find(&tmp).Error
	return tmp, err
}

// SearchBalance 搜索余额记录是否存在
func (db *DB) SearchBalance(contract, address string) (*Balance, error) {
	var tmp Balance
	err := global.SHOP_DB.Where("contract = ? and address =?", contract, address).Find(&tmp).Error
	if err != nil {
		return nil, err
	}
	return &tmp, nil
}

// InsertBalance 插入数据
func (db *DB) InsertBalance(account *Balance) (int64, error) {
	re, _ := db.SearchBalance(account.Contract, account.Address)
	if re != nil {
		err := global.SHOP_DB.Model(Balance{}).Where("id=?", re.ID).Find(&account).Updates(map[string]interface{}{"amount": account.Amount}).Error
		if err != nil {
			global.SHOP_LOG.Error(err.Error())
			return 0, err
		}
	}
	err := global.SHOP_DB.Model(Balance{}).Save(&account).Error
	if err != nil {
		global.SHOP_LOG.Error(err.Error())
		return 0, err
	}
	return account.ID, nil
}

// GetAccountWithContractBalance 获取大于minAmount的所有账户合约余额
func (db *DB) GetAccountWithContractBalance(contract string, minAmount, startid int64, count int) ([]Balance, error) {
	var tmp = make([]Balance, 0)
	err := global.SHOP_DB.Where("contract= ? and amount >= ? and id > ?", contract, minAmount, startid).Order("id").Limit(count).Find(&tmp).Error
	return tmp, err
}

// GetSumContractBalance 获取合约总余额
func (db *DB) GetSumContractBalance(contract string) (map[string]int64, error) {
	var tmp = make(map[string]int64, 0)
	err := global.SHOP_DB.Table("balance").Select("sum(amount) as sumall").Where("contract= ?", contract).Find(&tmp).Error
	return tmp, err
}

// GetTransactions 获取最近交易记录
func (db *DB) GetTransactions(contract, addr string, count, skip int) ([]Transactions, error) {
	var tmp = make([]Transactions, 0)

	if count < 1 || count > 1000 {
		count = 300
	}
	if skip < 0 {
		skip = 0
	}
	skips := strconv.Itoa(skip)
	tmpdb := global.SHOP_DB.Scopes(Paginate.Paginate(skips, global.SHOP_CONFIG.System.PageSize)).Where("(type=? OR type=?) And contract =? ", Send, Receive, contract)
	if addr != "*" && addr != "" {
		tmpdb = tmpdb.Where("addr = ?", addr)
	}
	err := tmpdb.Order("id desc").Find(&tmp).Error
	return tmp, err
}

// GetCollestTransactions 获取指定时间段内归集交易记录
func (db *DB) GetCollestTransactions(sTime, eTime int64, contract string) ([]Transactions, error) {
	var tmp = make([]Transactions, 0)
	if eTime < sTime || eTime < 1 {
		eTime = 0
	}
	tmpdb := global.SHOP_DB.Where("type=? and contract=? and timesmap>=? ", Collect, contract, sTime)
	if eTime > 1 {
		tmpdb = tmpdb.Where("timesmap<=?", eTime)
	}
	err := tmpdb.Order("id desc").Find(&tmp).Error
	return tmp, err
}

// SearchTransactions 搜索交易记录是否存在
func (db *DB) SearchTransactions(txid string) (*Transactions, error) {
	var tmp Transactions
	err := global.SHOP_DB.Where("tx_id = ?", txid).Limit(1).Find(&tmp).Error
	if err != nil {
		global.SHOP_LOG.Error(err.Error())
		return nil, err
	}

	return &tmp, nil
}

// InsertTransactions 插入数据
func (db *DB) InsertTransactions(transactions *Transactions) (int64, error) {
	re, _ := db.SearchTransactions(transactions.TxID)
	if re != nil {
		return 0, nil
	}
	err := global.SHOP_DB.Model(Transactions{}).Save(&transactions).Error
	if err != nil {
		global.SHOP_LOG.Error(err.Error())
		return 0, err
	}
	return transactions.ID, nil
	//return db.Cols("tx_id", "block_height", "address", "public_key", "from_address", "contract",
	//	"amount", "fee", "timestamp", "type").Insert(transactions)
}

// LoadLastBlockHeight 获取最后一次扫描高度 已经扫描到这个高度
func (db *DB) LoadLastBlockHeight() (int64, error) {
	var tmp OtherParam
	err := global.SHOP_DB.Where("key='block'").Limit(1).Find(&tmp).Error
	if err != nil || tmp.Value == "0" {
		return 0, errors.New("没用记录")
	}
	var un int64
	un, err = strconv.ParseInt(tmp.Value, 10, 0)
	return un, err
}

// InsertLastBlockHeight 更新最后一次扫描高度
func (db *DB) InsertLastBlockHeight(num int64) (err error) {

	var tmp = OtherParam{
		Key: "block",
	}
	global.SHOP_DB.Where("key=?", "block").Find(&tmp)
	if tmp.Key != "" {
		err = global.SHOP_DB.Save(&tmp).Error
		if err != nil {
			global.SHOP_LOG.Error(err.Error())
			return err
		}
	} else {
		tmp.Value = strconv.FormatInt(num, 10)
		err = global.SHOP_DB.Model(OtherParam{}).Where("key='block'").Updates(&tmp).Error
	}
	return
}
