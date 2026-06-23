package middleware

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"strings"
	"wallet_chain.com/utils"

	"wallet_chain.com/admin/model"
	"wallet_chain.com/admin/service"
	"wallet_chain.com/global"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// 后台中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		session := sessions.Default(c)

		adminId := session.Get("adminId")
		sessionid := utils.MD5V(session.ID())
		result, err := global.SHOP_REDIS.Get(c, fmt.Sprintf("adminlogin_%d", adminId)).Result()
		fmt.Println(result, err, sessionid)
		if err == redis.Nil {
			session.Clear() // 清除所有数据
			if err := session.Save(); err != nil {
				println(err)

			}
			c.Abort()
			c.Redirect(302, "/public/login")
		}
		if result != sessionid {
			session.Clear() // 清除所有数据
			if err := session.Save(); err != nil {
				println(err)

			}
			c.Abort()
			c.Redirect(302, "/public/login")
		}
		errCode := 0

		if adminId == nil || adminId == 0 {
			c.Abort()
			c.Redirect(302, "/public/login")
		} else {
			//验证权限
			admin := model.Admin{}
			admin.Id = adminId.(uint)
			err := global.SHOP_DB.Table("nov_admin").Select("nov_admin.*,nov_group.role_ids ").Joins("left join nov_group on nov_admin.group_id=nov_group.id  ").Find(&admin)

			if err != nil {
				fmt.Println(err.Error)
			}

			if admin.Id == 0 || *admin.Status != 1 {

				c.Abort()
				c.JSON(200, gin.H{
					"code": errCode,
					"msg":  "账号不存在或被禁用",
				})

			} else {

				path := c.Request.URL.Path

				urlsli := strings.Split(path, "/")

				key := strings.Join(urlsli[len(urlsli)-2:], "")

				if key == "adminlogout" {
					c.Next()
				} else {

					v, ok := global.BlackCache.Get(key)
					if ok {
						println(v)
					} else {
						c.Abort()
						c.JSON(200, gin.H{
							"code": errCode,
							"msg":  "无权操作BlackCache",
						})
						return
					}

					//继续验证权限
					//res := true //model.AuthCheck(path, admin.Roles)
					var server service.Role
					res := server.AuthCheck(key, admin.GroupId)
					fmt.Println(admin.GroupId, "groul-----------------")
					if !res {
						c.Abort()
						c.JSON(200, gin.H{
							"path": key,
							"code": errCode,
							"msg":  "无权操作",
						})
						return
					} else {
						c.Set("adminId", adminId)
						//	c.Set("adminRoles", admin.RoleId)
					}
				}
			}
		}

		//c.Next()
		//终止后续方法执行，但中间件后面的内容会继续执行
		//c.Abort()
	}
}
