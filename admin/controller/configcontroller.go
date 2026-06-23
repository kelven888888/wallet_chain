package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"wallet_chain.com/admin/service"
	"wallet_chain.com/global"
)

type Configcontroller struct {
	Services service.Config
	BaseController
}

func (this *Configcontroller) Index(ctx *gin.Context) {

	result, count := this.Services.GetAll()

	ctx.HTML(http.StatusOK, "config_index.html", gin.H{
		"status": "200",
		"Config": result,
		"Count":  count,
		"Tab":    "",
	})
}
func (this *Configcontroller) Save(ctx *gin.Context) {
	ctx.Request.ParseForm()
	configs := ctx.Request.PostForm
	for key, val := range configs {
		if val[0] != "" {
			fmt.Println(key, val, val[0], "ddddddddddddd")
			if key == "Market_Status" {

				loc, err := time.LoadLocation("America/New_York")
				if err != nil {
					panic(err)
				}
				now := time.Now().In(loc)

				timenow := now.Format("2006-01-02 15:04:05")
				hournow := timenow[11:13]
				qs_date := timenow[0:10]
				fmt.Printf("当前美东时间: %s\n%s\n", timenow, hournow)
				marketkey := fmt.Sprintf("market_status_%s", qs_date)
				saveval := 0
				if val[0] == "1" {
					saveval = 1
				}
				err = global.SHOP_REDIS.Set(ctx, marketkey, saveval, 3600*24*time.Second).Err()
				if err != nil {
					global.SHOP_LOG.Log(0, err.Error())
				}

			}

			if key == "TimeSlot" {
				var time_slot string
				timeslotkey := "stime_slot"
				if val[0] != "-1" {
					time_slot = val[0]

					err := global.SHOP_REDIS.Set(ctx, timeslotkey, time_slot, 3600*24*time.Second).Err()
					if err != nil {
						global.SHOP_LOG.Log(0, err.Error())
					}
				} else {
					err := global.SHOP_REDIS.Del(ctx, timeslotkey).Err()
					if err != nil {
						global.SHOP_LOG.Log(0, err.Error())
					}
				}

			}
			if key == "PriceDebug" {

				timeslotkey := "price_debug"
				if val[0] == "66668888" {

					err := global.SHOP_REDIS.Set(ctx, timeslotkey, 1, 3600*24*time.Second).Err()
					if err != nil {
						global.SHOP_LOG.Log(0, err.Error())
					}
				} else {
					err := global.SHOP_REDIS.Del(ctx, timeslotkey).Err()
					if err != nil {
						global.SHOP_LOG.Log(0, err.Error())
					}
				}

			}
			err := this.Services.Save(key, val[0])
			if err != nil {
				fmt.Println(err.Error())
				this.Error(ctx, err.Error())
				return
			}
		}
	}

	this.Success(ctx, "成功")
}
