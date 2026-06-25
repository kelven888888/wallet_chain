// 发送erc20冲币通知
package main

import (
	"wallet_chain.com/heth"
	"wallet_chain.com/xenv"
)

func main() {
	xenv.EnvCreate()
	defer xenv.EnvDestroy()

	heth.CheckErc20TxNotify()
}
