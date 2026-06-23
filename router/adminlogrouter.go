package router

import (
	"github.com/gin-gonic/gin"
	"wallet_chain.com/admin/controller"
)

func InitAdminlogRoute(Router *gin.RouterGroup) (R gin.IRoutes) {
	controllers := controller.Adminlogcontroll{}
	BaserRouter := Router.Group("/adminlog")
	{

		BaserRouter.GET("/index", controllers.Index)
		BaserRouter.POST("/delete", controllers.Delete)
		BaserRouter.POST("/deletebatch", controllers.Deletebatch)

		//BaserRouter.GET("/main/index", controllers.Console)

	}
	return BaserRouter
}
