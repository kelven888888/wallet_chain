package middleware

import (
	"github.com/gin-gonic/gin"
)

// 全局中间件
func AllowMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		//中间件协程不能直接使用context 要复制后使用
		//newc := c.Copy()

		//中间件协程不需要使用 sync.WriterGroup 等待协程,需注意静态文件加载，也会经过中间件
		// go func() {

		// 	time.Sleep(time.Second * 2)

		// 	fmt.Println("协程", time.Now().Unix(), newc.Request.Method)
		// }()

		//设置数据
		//c.Set("test1", "test1")

		c.Next()
	}
}
