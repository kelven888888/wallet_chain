package crondtab

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"wallet_chain.com/admin/model"
	"wallet_chain.com/global"
)

var ctx = context.Background()
var staticviste = func() {
	var dates StringSlice
	var vistecount Int64Slice
	var ipvistecount Int64Slice
	timeten := time.Now().AddDate(0, 0, -10)
	savakeyipcount := "savakeyipcount"
	savakeyipvistecount := "pv_ipvistecount"
	savadate := "savadate"
	result1, _ := global.SHOP_REDIS.Get(ctx, savakeyipcount).Result()
	result2, _ := global.SHOP_REDIS.Get(ctx, savakeyipvistecount).Result()
	result3, _ := global.SHOP_REDIS.Get(ctx, savadate).Result()
	if result1 != "" && result2 != "" && result3 != "" {
		json.Unmarshal([]byte(result1), &vistecount)
		json.Unmarshal([]byte(result2), &ipvistecount)
		json.Unmarshal([]byte(result3), &dates)
		fmt.Println(dates, "缓存-----------------------")
	} else {
		for i := 0; i <= 10; i++ {
			datetimenow := timeten.AddDate(0, 0, i).Format("2006-01-02")
			//dateyesterday := timeten.AddDate(0, 0, i-1).Format("2006-01-02")
			datetomorrow := timeten.AddDate(0, 0, i+1).Format("2006-01-02")
			var count int64
			global.SHOP_DB.Model(model.MAccesslog{}).Where("create_at>? and create_at<?", datetimenow, datetomorrow).Count(&count)
			var ipcountarr []int64
			var ipcount int64
			global.SHOP_DB.Model(model.MAccesslog{}).Select("DISTINCT(ip)").Where("create_at>? and create_at<?", datetimenow, datetomorrow).Find(&ipcountarr)
			vistecount = append(vistecount, count)
			ipcount = int64(len(ipcountarr))
			ipvistecount = append(ipvistecount, ipcount)
			dates = append(dates, datetimenow)

		}

		ipvistecount, _ := ipvistecount.MarshalJSON()
		vistecounts, _ := vistecount.MarshalJSON()
		datess, _ := dates.MarshalJSON()
		expire := time.Minute * 5
		err := global.SHOP_REDIS.Set(ctx, savakeyipcount, vistecounts, expire).Err()
		if err != nil {
			global.SHOP_LOG.Log(2, err.Error())
		}
		err = global.SHOP_REDIS.Set(ctx, savakeyipvistecount, ipvistecount, expire).Err()
		if err != nil {
			global.SHOP_LOG.Log(2, err.Error())
		}
		err = global.SHOP_REDIS.Set(ctx, savadate, datess, expire).Err()
		if err != nil {
			global.SHOP_LOG.Log(2, err.Error())
		}

	}

}

type Int64Slice []int64

func (s Int64Slice) MarshalJSON() ([]byte, error) {
	// 实现自定义的 JSON 编组逻辑
	// 例如，直接返回 JSON 数组格式
	return json.Marshal([]int64(s))
}

type StringSlice []string

func (s StringSlice) MarshalJSON() ([]byte, error) {
	// 实现自定义的 JSON 编组逻辑
	// 例如，直接返回 JSON 数组格式
	return json.Marshal([]string(s))
}
