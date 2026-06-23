package core

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
	"wallet_chain.com/admin/service"
	"wallet_chain.com/global"
	"wallet_chain.com/initialize"

	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer() {

	Router := initialize.IninRoute()
	//Router.Static("/form-generator", "./resource/page")
	//Router.StaticFS("assets", http.FS(Static))
	//// 设置模板资源
	//Router.SetHTMLTemplate(template.Must(template.New("").ParseFS(assets.Templates, "templates/**/*"))
	address := fmt.Sprintf(":%d", global.SHOP_CONFIG.System.Addr)
	s := initServer(address, Router)
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	global.SHOP_LOG.Info("server run success on ", zap.String("address", address))
	service.Getmenu()
	fmt.Printf(`
	
	默认自动化文档地址:http://localhost%s/swagger/index.html
	默认前端文件运行地址:http://localhost%s
`, address, address)
	global.SHOP_LOG.Error(s.ListenAndServe().Error())

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	HttpServerStop()
	log.Println("Server exiting")

}
