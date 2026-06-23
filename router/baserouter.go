package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"wallet_chain.com/admin/controller"
)

func InitPublicRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	controllers := controller.Publiccontroll{}
	wallcontrollers := controller.WalletCtr{}
	Router.GET("/", func(context *gin.Context) {
		context.Abort()
		context.Redirect(302, "/public/login")

	})

	Router.Any("/wallet/payoutcalback", wallcontrollers.PayoutBack)

	BaserRouter := Router.Group("public")
	{
		BaserRouter.GET("/ping", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "pong",
			})

		})
		BaserRouter.POST("/message/index", func(context *gin.Context) {
			fmt.Println(1)

		})
		BaserRouter.GET("/register", controllers.Register)
		BaserRouter.GET("/captcha", controllers.Captcha)
		BaserRouter.GET("/login", controllers.Login)
		BaserRouter.POST("/loginsubmit", controllers.Loginsubmit)
		BaserRouter.POST("/upload", controllers.Upload)
		BaserRouter.POST("/uploadeditor", controllers.Uploadeditor)

	}
	return BaserRouter
}
