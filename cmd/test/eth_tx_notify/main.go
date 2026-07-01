// 检测入账广播
package main

import (
	"wallet_chain.com/cores"
	"wallet_chain.com/global"
	"wallet_chain.com/heth"
	"wallet_chain.com/initialize"
	"wallet_chain.com/xenv"
)

func main() {
	xenv.EnvCreate()
	defer xenv.EnvDestroy()
	cores.Viper()                      // 初始化Viper
	global.SHOP_LOG = cores.Zap()      // 初始化zap日志库
	global.SHOP_DB = initialize.Gorm() // gorm连接数据库
	initialize.OtherInit()
	initialize.Redis()
	heth.CheckTxNotify()
}
