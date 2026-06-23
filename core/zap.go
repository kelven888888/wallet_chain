package core

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"wallet_chain.com/core/internal"
	"wallet_chain.com/global"
	"wallet_chain.com/utils"
)

// Zap 获取 zap.Logger
// Author [SliverHorn](https://github.com/SliverHorn)
func Zap() (logger *zap.Logger) {
	if ok, _ := utils.PathExists(global.SHOP_CONFIG.Zap.Director); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", global.SHOP_CONFIG.Zap.Director)
		_ = os.Mkdir(global.SHOP_CONFIG.Zap.Director, os.ModePerm)
	}
	levels := global.SHOP_CONFIG.Zap.Levels()
	length := len(levels)
	cores := make([]zapcore.Core, 0, length)
	for i := 0; i < length; i++ {
		core := internal.NewZapCore(levels[i])
		cores = append(cores, core)
	}
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段,如：添加一个服务器名称
	filed := zap.Fields(zap.String("application", "shop"))
	logger = zap.New(zapcore.NewTee(cores...), caller, development, filed)
	//logger = zap.New(zapcore.NewTee(cores...))
	if global.SHOP_CONFIG.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}
