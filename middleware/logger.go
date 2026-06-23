package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
	"time"
	"wallet_chain.com/admin/model"
	"wallet_chain.com/global"
	"wallet_chain.com/utils"
	"wallet_chain.com/utils/logger"
)

// gin自定义日志中间件
func LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {

		//开始时间
		startTime := time.Now()
		language := c.GetHeader("Accept-Language")
		allowlanguage := strings.Split(global.SHOP_CONFIG.System.Language_Array, ",")

		if !utils.InArray(language, allowlanguage) {
			language = global.SHOP_CONFIG.System.Language
		}

		c.Set("Language", language)
		fmt.Println(c.Get("Language"))
		//处理请求
		c.Next()
		//结束时间
		endTime := time.Now()
		//执行时间
		latencyTime := endTime.Sub(startTime).Seconds()
		//请求方式
		reqMethod := c.Request.Method
		//请求路由
		reqUri := c.Request.RequestURI
		//状态码
		statusCode := c.Writer.Status()
		//请求IP
		ClientIp := c.ClientIP()
		//用户标识
		UserAgent := c.Request.UserAgent()

		//日志格式
		logger.Logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    ClientIp,
			"req_method":   reqMethod,
			"req_uri":      reqUri,
			"user_agent":   UserAgent,
		}).Info()
		tokenString := c.GetHeader("Authorization")
		// 验证token非空
		var username = ""
		if tokenString != "" {
			token, claims, err := utils.ParseToken(tokenString)
			if err == nil && token.Valid {

				username = claims.UserName
			}
			//如果用户存在 将user信息存入上下文

		}

		go func(username string) {
			var acclog model.MAccesslog
			acclog.Username = username
			acclog.Ip = ClientIp
			acclog.Path = reqUri
			acclog.CreateAt = time.Now()
			acclog.Method = reqMethod
			db, err := geoip2.Open("GeoLite2-City.mmdb")
			if err != nil {
				global.SHOP_LOG.Log(2, err.Error())
			}
			defer db.Close()
			ip := net.ParseIP(acclog.Ip)
			record, err := db.City(ip)
			if err != nil {
				global.SHOP_LOG.Log(2, err.Error())
			}
			global.SHOP_DB.Save(&acclog)
			address := fmt.Sprintf("%s_%s", record.Country.Names["zh-CN"], record.City.Names["zh-CN"])
			if len(record.Subdivisions) > 0 {
				address = fmt.Sprintf("%s_%s", address, record.Subdivisions[0].Names["zh-CN"])
			}
			if len(record.Subdivisions) > 0 {
				global.SHOP_DB.Model(model.MAccesslog{}).Where("id=?", acclog.Id).Updates(model.MAccesslog{
					City:        record.City.Names["zh-CN"],
					Country:     record.Country.Names["zh-CN"],
					Subdivision: record.Subdivisions[0].Names["zh-CN"],
					Address:     address,
				})
			} else {
				global.SHOP_DB.Model(model.MAccesslog{}).Where("id=?", acclog.Id).Updates(model.MAccesslog{
					City:    record.City.Names["zh-CN"],
					Country: record.Country.Names["zh-CN"],
					Address: address,
				})
			}

		}(username)

		// token验证是否失效

	}
}
