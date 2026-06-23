package router

import (
	"github.com/gin-gonic/gin"
	"wallet_chain.com/admin/controller"
)

func InitConfigRoute(Router *gin.RouterGroup) (R gin.IRoutes) {
	controllers := controller.Configcontroller{}
	BaserRouter := Router.Group("/config")
	{

		BaserRouter.GET("/index", controllers.Index)
		BaserRouter.POST("/save", controllers.Save)
		//BaserRouter.POST("/delete", controllers.Delete)
		//BaserRouter.POST("/deletebatch", controllers.Deletebatch)

		//BaserRouter.GET("/main/index", controllers.Console)

	}
	return BaserRouter
}
