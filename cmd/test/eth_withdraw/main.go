// 检测提现
package main

import (
	"wallet_chain.com/heth"
	"wallet_chain.com/xenv"
)

func main() {
	xenv.EnvCreate()
	defer xenv.EnvDestroy()

	heth.CheckWithdraw()
}
