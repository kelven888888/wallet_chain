package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func GrpcAllow() gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("client_ip", c.ClientIP())
		if c.ClientIP() == "127.0.0.1" {
			c.Next()
			return
		}
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "拒绝访问",
			"data": gin.H{},
		})
		return
	}
}
