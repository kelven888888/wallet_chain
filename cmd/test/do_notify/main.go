// 发送通知
package main

import (
	"wallet_chain.com/app"
	"wallet_chain.com/xenv"
)

func main() {
	xenv.EnvCreate()
	defer xenv.EnvDestroy()

	app.CheckDoNotify()
}
