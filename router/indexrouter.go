package router

import (
	"github.com/gin-gonic/gin"
	"wallet_chain.com/admin/controller"
)

func InitIndexRoute(Router *gin.RouterGroup) (R gin.IRoutes) {
	controllers := controller.IndexController{}
	Router.GET("/main/index", controllers.Index)
	Router.Any("/home/index", controllers.Console)

	return Router
}
