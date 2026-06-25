package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"wallet_chain.com/admin/model"
	"wallet_chain.com/admin/service"
	"wallet_chain.com/global"
	"wallet_chain.com/utils"
)

type IndexController struct {
	BaseController
}

func (this *IndexController) Index(ctx *gin.Context) {
	server := service.GroupServer{}
	session := sessions.Default(ctx)

	groupId := session.Get("groupId").(uint)

	Menus, err := server.GetMenus(groupId)
	if err != nil {

	}

	ctx.HTML(http.StatusOK, "main_index.html", gin.H{

		"Menus": Menus,
	})

}
func (this *IndexController) Console(ctx *gin.Context) {
	if ctx.Request.Method == "GET" {
		server := service.GroupServer{}
		session := sessions.Default(ctx)

		groupId := session.Get("groupId").(uint)

		Menus, err := server.GetMenus(groupId)
		var menu []model.Role
		for _, v := range Menus {
			for _, vs := range v.Child {
				menu = append(menu, *vs)
			}
		}
		if err != nil {

		}
		date := utils.Us_datatime()
		var heartbeat model.Heartbeat
		global.SHOP_DB.Where("type=1").Order("id desc ").First(&heartbeat)
		var crondlognetval model.CrontabLog
		global.SHOP_DB.Where("type=2").Order("id desc ").Limit(1).First(&crondlognetval)
		var crondlogaiorder model.CrontabLog
		global.SHOP_DB.Where("type=3").Order("id desc ").Limit(1).First(&crondlogaiorder)
		var crondlogredeem model.CrontabLog
		global.SHOP_DB.Where("type=4").Order("id desc ").Limit(1).First(&crondlogredeem)

		var crondlogsettle model.CrontabLog
		global.SHOP_DB.Where("type=1").Order("id desc ").Limit(1).First(&crondlogsettle)

		//var showcast model.QuantShowcaseData
		//global.SHOP_DB.Where("qs_date=?", date[0:10]).Order("id desc ").Limit(1).First(&showcast)
		marketkey := fmt.Sprintf("market_status_%s", date[0:10])
		marketstatus, _ := global.SHOP_REDIS.Get(ctx, marketkey).Result()
		marketopen := 0
		if marketstatus == "1" {
			marketopen = 1
		}
		var undoWithdrawcount int64
		//global.SHOP_DB.Model(model.UsdtWithdrawModel{}).Where("status=0").Count(&undoWithdrawcount)

		var rechargecount int64
		//global.SHOP_DB.Model(model.FundRecharge{}).Count(&rechargecount)

		totalnum := undoWithdrawcount + rechargecount

		// 设置session数据
		session.Set("noticnum", totalnum)

		// 保存session数据
		session.Save()

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
		marketserverkey := fmt.Sprintf("market_status_%s", date[0:10])
		marketserver, _ := global.SHOP_REDIS.Get(ctx, marketserverkey).Result()

		ctx.HTML(http.StatusOK, "console.html", gin.H{

			"Menus":            menu,
			"Date":             date,
			"Heartbeattime":    heartbeat.UpdatedAt.Format(time.DateTime),
			"updatenetvaltime": crondlognetval.CreatedAt.Format(time.DateTime),
			"aioudertime":      crondlogaiorder.CreatedAt.Format(time.DateTime),
			"airedeemtime":     crondlogredeem.CreatedAt.Format(time.DateTime),
			"aisettletime":     crondlogsettle.CreatedAt.Format(time.DateTime),

			"marketopen":        marketopen,
			"undoWithdrawcount": undoWithdrawcount,

			"dates":        dates,
			"vistecount":   vistecount,
			"ipvistecount": ipvistecount,
			"marketserver": marketserver,
		})
	} else {

		//var undoWithdrawcount int64
		//global.SHOP_DB.Model(model.UsdtWithdrawModel{}).Where("status=0").Count(&undoWithdrawcount)
		//
		//var undotradeapplyaccount int64
		//
		//var rechargecount int64
		//global.SHOP_DB.Model(model.FundRecharge{}).Count(&rechargecount)
		//totalnum := undoWithdrawcount + undotradeapplyaccount + rechargecount
		//
		//session := sessions.Default(ctx)
		//noticnum := session.Get("noticnum").(int64)
		//fmt.Println(noticnum, totalnum)
		//if noticnum != totalnum {
		//	ctx.String(200, fmt.Sprintf("%d", totalnum))
		//} else {
		//	ctx.String(200, fmt.Sprintf("%s", "0"))
		//}
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
