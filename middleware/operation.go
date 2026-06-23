package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"wallet_chain.com/admin/model"

	//"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"wallet_chain.com/admin/service"
	"wallet_chain.com/global"
)

var operationRecordService = service.OperationRecordService{}

var respPool sync.Pool
var bufferSize = 1024

func init() {
	respPool.New = func() interface{} {
		return make([]byte, bufferSize)
	}
}

func OperationRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data []byte
		var body []byte
		var userId uint
		if c.Request.Method != http.MethodGet {

			var err error
			body, err = io.ReadAll(c.Request.Body)
			if err != nil {
				global.SHOP_LOG.Error("read body from request error:", zap.Error(err))
			} else {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		} else {

		}
		session := sessions.Default(c)

		adminId := session.Get("adminId")
		userId = adminId.(uint)
		admin := model.Admin{}
		global.SHOP_DB.Where("id=?", userId).Find(&admin)
		username := admin.Account
		m := make(map[string]interface{})
		err := json.Unmarshal(body, &m)
		if err != nil {
			fmt.Println(err.Error())
		}
		_, ok := m["password"]
		if ok {
			m["password"] = "******"
		}
		_, ok = m["tradepassword"]
		if ok {
			m["tradepassword"] = "******"
		}

		if len(m) == 0 {
			c.Request.ParseForm()
			for key, value := range c.Request.PostForm {
				m[key] = value
			}
			data, _ = json.Marshal(&m)
		} else {
			data, _ = json.Marshal(&m)

		}

		record := model.SysOperationRecord{
			//Ip: c.ClientIP(),
			Ip:       "127.0.0.1",
			Method:   c.Request.Method,
			Path:     c.Request.URL.Path,
			Agent:    c.Request.UserAgent(),
			Body:     string(data),
			UserID:   userId,
			UserName: username,
		}
		//fmt.Println(c.Request.URL.Path)
		//fmt.Println(string(body))

		// 上传文件时候 中间件日志进行裁断操作
		if strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
			if len(record.Body) > bufferSize {
				// 截断
				newBody := respPool.Get().([]byte)
				copy(newBody, record.Body)
				record.Body = string(newBody)
				defer respPool.Put(newBody)
			}
		}

		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer
		now := time.Now()

		c.Next()

		latency := time.Since(now)
		record.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		record.Status = c.Writer.Status()
		record.Latency = latency
		//record.Resp = writer.body.String()

		if strings.Contains(c.Writer.Header().Get("Pragma"), "public") ||
			strings.Contains(c.Writer.Header().Get("Expires"), "0") ||
			strings.Contains(c.Writer.Header().Get("Cache-Control"), "must-revalidate, post-check=0, pre-check=0") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/force-download") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/octet-stream") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/vnd.ms-excel") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/download") ||
			strings.Contains(c.Writer.Header().Get("Content-Disposition"), "attachment") ||
			strings.Contains(c.Writer.Header().Get("Content-Transfer-Encoding"), "binary") {
			if len(record.Resp) > bufferSize {
				// 截断
				newBody := respPool.Get().([]byte)
				copy(newBody, record.Resp)
				record.Resp = string(newBody)
				defer respPool.Put(newBody)
			}
		}
		if c.Request.Method != http.MethodGet {
			//if c.Request.Method != "sd" {
			if err := operationRecordService.CreateSysOperationRecord(record); err != nil {
				global.SHOP_LOG.Error("create operation record error:", zap.Error(err))
			}
		}
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
