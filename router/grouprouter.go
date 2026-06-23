package router

import (
	"github.com/gin-gonic/gin"
	"wallet_chain.com/admin/controller"
)

func InitGroupRoute(Router *gin.RouterGroup) (R gin.IRoutes) {
	var GroupController = controller.GroupController{}
	RoleRouter := Router.Group("/group/")
	{
		RoleRouter.GET("/index", GroupController.Index)
		RoleRouter.Any("/edit", GroupController.Edit)
		RoleRouter.Any("/add", GroupController.Add)
		RoleRouter.Any("/delete", GroupController.Delete)
		/*	RoleRouter.Any("/roleedit", GroupController.Roleedit)
			RoleRouter.GET("/getlist", GroupController.Getlist)
			RoleRouter.GET("/getmenu", GroupController.Getmenu)
			RoleRouter.GET("/roleadd", GroupController.RoleAdd)
			RoleRouter.POST("/roledel", GroupController.RoleDel)*/

	}
	return RoleRouter

}
