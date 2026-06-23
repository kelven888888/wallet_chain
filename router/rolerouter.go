package router

import (
	"github.com/gin-gonic/gin"
	"wallet_chain.com/admin/controller"
)

func InitRoleRoute(Router *gin.RouterGroup) (R gin.IRoutes) {
	controllers := controller.RoleController{}
	BaserRouter := Router.Group("role")
	{

		BaserRouter.GET("/index", controllers.Index)
		BaserRouter.Any("/edit", controllers.Edit)
		BaserRouter.Any("/add", controllers.Add)
		BaserRouter.POST("/delete", controllers.Delete)

	}
	return BaserRouter
}
