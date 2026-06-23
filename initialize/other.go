package initialize

import (
	"fmt"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"strings"

	"wallet_chain.com/global"
	"wallet_chain.com/utils"
)

func OtherInit() {
	dr, err := utils.ParseDuration(global.SHOP_CONFIG.JWT.ExpiresTime)
	if err != nil {
		panic(err)
	}
	global.BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(dr),
	)
	allowlanguage := strings.Split(global.SHOP_CONFIG.System.Language_Array, ",")
	for _, v := range allowlanguage {
		jsons, err := utils.ReadFileContent(fmt.Sprintf("./language/%s.json", v))
		if err != nil {
			fmt.Println(err.Error())
		}
		result := strings.Join(jsons, "")                                   // 使用空格作为分隔符
		global.BlackCache.SetDefault(fmt.Sprintf("language_%s", v), result) // 输
	}

}
