// erc20 提币
package main

import (
	"wallet_chain.com/heth"
	"wallet_chain.com/xenv"
)

func main() {
	xenv.EnvCreate()
	defer xenv.EnvDestroy()

	heth.CheckErc20Withdraw()
}
