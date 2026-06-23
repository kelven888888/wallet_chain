package router

import (
	"github.com/gin-gonic/gin"
	"wallet_chain.com/admin/controller"
)

func InitSettingRoute(Router *gin.RouterGroup) (R gin.IRoutes) {
	controllers := controller.SettingController{}
	BaserRouter := Router.Group("/setting")
	{
		BaserRouter.Any("/poster", controllers.Poster)
	}
	return BaserRouter
}
