package route

import (
	"github.com/gin-gonic/gin"
	"wallet_chain.com/middleware"
)

func Routers(r *gin.Engine) {
	api := r.Group("/api", middleware.LoggerToFile())
	api.Use(middleware.LoggerToFile())
	{

	}
}
