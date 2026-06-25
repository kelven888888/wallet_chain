// erc20零钱整理
package main

import (
	"wallet_chain.com/heth"
	"wallet_chain.com/xenv"
)

func main() {
	xenv.EnvCreate()
	defer xenv.EnvDestroy()

	heth.CheckErc20TxOrg()
}
