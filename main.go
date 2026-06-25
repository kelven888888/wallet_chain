package main

import (
	"embed"
	_ "embed"
	"os"
	"time"
	ethcrontab "wallet_chain.com/cmd/crontab"
	"wallet_chain.com/cores"
	"wallet_chain.com/crondtab"
	"wallet_chain.com/gen"
	"wallet_chain.com/global"
	"wallet_chain.com/initialize"
	trxserver "wallet_chain.com/trx"
)

func timeDifferenceInMinutes(t1, t2 time.Time) int {
	diff := t1.Sub(t2)
	return int(diff.Minutes())
}

//go:embed  views/admin/include/*
var Templatess embed.FS

func main() {
	//channle := make(chan uint, 10)
	//for i := 0; i < 1000; i++ {
	//	select {
	//	case channle <- 1:
	//		go func() {
	//
	//			fmt.Println(i)
	//			time.Sleep(time.Second * 1)
	//			<-channle
	//		}()
	//	}
	//}
	const (
		layoutsDir   = "templates/layouts"
		templatesDir = "views/admin/include/"
		extension    = "/*.html"
	)

	cores.Viper()                      // 初始化Viper
	global.SHOP_LOG = cores.Zap()      // 初始化zap日志库
	global.SHOP_DB = initialize.Gorm() // gorm连接数据库
	initialize.OtherInit()
	initialize.Redis()

	go trxserver.Init()
	go ethcrontab.Init()

	select {}

	//ctx := context.Background()
	////go service.RunPublisher(ctx, "shop_message")
	//go service.RunSubscriber(ctx, "shop_message", "shop_message_queue")
	//time.Local = time.FixedZone("US/Eastern", -4*3600)
	time.Local = time.FixedZone("Asia/Shanghai", 0)
	if global.SHOP_DB != nil {

		// 程序结束前关闭数据库链接
		db, _ := global.SHOP_DB.DB()
		defer db.Close()
	}
	args := os.Args[1:]

	if len(args) != 0 {
		if args[0] == "runcrond" {

			go crondtab.Initcrond()

		}
	}

	if len(args) != 0 {

		if args[0] == "gen" {
			gen := gen.Gen{}
			gen.Gener(args[1], args[2], args[3])
			select {}
		}

	}

	cores.RunWindowsServer()
	select {}
}
