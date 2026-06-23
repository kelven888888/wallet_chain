package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"wallet_chain.com/global"
	"wallet_chain.com/utils"
)

// 验证token
func CheckToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取 Authorization header 头
		tokenString := ctx.GetHeader("Authorization")
		// 验证token非空
		if tokenString == "" {
			utils.Response(ctx, http.StatusUnauthorized, 403, "未登录", nil)
			ctx.Abort()
			return
		}
		// token验证是否失效
		token, claims, err := utils.ParseToken(tokenString)
		if err != nil || !token.Valid {
			utils.Response(ctx, http.StatusUnauthorized, 403, "登录已过期", nil)
			ctx.Abort()
			return
		}
		key := fmt.Sprintf("login_%s", claims.UserId)
		val, _ := global.SHOP_REDIS.Get(ctx, key).Result()
		if val == "" || val != tokenString {
			utils.Response(ctx, http.StatusUnauthorized, 403, "登录已过期", nil)
			ctx.Abort()
			return
		}
		//如果用户存在 将user信息存入上下文
		ctx.Set("user_id", claims.UserId)
		ctx.Set("user_name", claims.UserName)
		ctx.Next()
	}
}
