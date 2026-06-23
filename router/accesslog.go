package router

import (
	"github.com/gin-gonic/gin"
	"wallet_chain.com/admin/controller"
)

func InitAccesslogRoute(Router *gin.RouterGroup) (R gin.IRoutes) {
	controllers := controller.Acclogctr{}
	BaserRouter := Router.Group("accesslog")
	{

		BaserRouter.GET("/index", controllers.Index)

	}
	return BaserRouter
}
