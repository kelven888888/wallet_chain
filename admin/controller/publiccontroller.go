package controller

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"net/http"
	"time"
	"wallet_chain.com/admin/model"
	requests "wallet_chain.com/admin/model/common/request"
	"wallet_chain.com/admin/service"
	"wallet_chain.com/global"
	"wallet_chain.com/utils"
)

var store = base64Captcha.DefaultMemStore

type Publiccontroll struct {
	BaseController
}

func (this *Publiccontroll) Register(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{})

}
func (this *Publiccontroll) Captcha(ctx *gin.Context) {
	var service service.Captcha
	var captcha = service.Captcha(ctx.ClientIP())
	this.Success(ctx, "成功", captcha)

}
func (this *Publiccontroll) Login(ctx *gin.Context) {
	session := sessions.Default(ctx)

	adminId := session.Get("adminId")
	//fmt.Println(adminId)
	if adminId != nil {
		ctx.Abort()
		ctx.Redirect(302, "/admin/main/index")

	}
	ctx.HTML(http.StatusOK, "admin_login.html", nil)

}

func (this *Publiccontroll) Loginsubmit(ctx *gin.Context) {
	var loginreq requests.Login
	err := ctx.ShouldBind(&loginreq)

	if err != nil {
		this.Error(ctx, "参数错误", err.Error())
		return
	}

	if !store.Verify(loginreq.CaptchaId, loginreq.Code, true) {
		this.Error(ctx, "验证码错误", nil)
		return
	}
	var adminserver service.AdminServer
	var adminuser *model.Admin
	adminuser, err = adminserver.Login(loginreq)
	//println(global.SHOP_CONFIG.System.SecretKey)

	if err != nil {
		this.Error(ctx, err.Error(), nil)
		return
	}
	session := sessions.Default(ctx)

	// 设置session数据
	session.Set("adminId", adminuser.Id)
	session.Set("groupId", adminuser.GroupId)
	// 保存session数据
	session.Save()
	sessionid := utils.MD5V(session.ID())
	err = global.SHOP_REDIS.Set(ctx, fmt.Sprintf("adminlogin_%d", adminuser.Id), sessionid, 3600*24*7*time.Second).Err()
	if err != nil {
		fmt.Println(sessionid, fmt.Sprintf("adminlogin_%d", adminuser.Id), "---------------------------------------------------------------", err, "---------------------------------------------------------------")
	}
	//var data = make(map[string]string)
	//data["access_token"] = "c262e61cd13ad99fc650e6908c7e5e65b63d2f32185ecfed6b801ee3fbdd5c0a"
	this.Success(ctx, "登录成功", nil)

}

func (this *Publiccontroll) Upload(ctx *gin.Context) {

	//ctx.JSON(http.StatusOK, gin.H{
	//	"error": 0,
	//	"url":   "../../uploads/file/9299961eab1ff4e85e870912f2abb560_20240529172250.jpg",
	//})
	//return
	_, header, err := ctx.Request.FormFile("file")
	contentType := header.Header.Get("Content-Type")
	allowedTypes := []string{"application/pdf", "image/jpeg", "image/png", "image/gif"}
	isAllowed := false
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		this.Error(ctx, "上传文件错误", nil)
		return
	}

	types := ctx.Query("type")
	if err != nil {
		this.Error(ctx, "上传文件错误", nil)
		return
	}
	var service service.FileUploadAndDownloadService
	result, err := service.UploadFile(header, "0", types)
	if err != nil {
		this.Error(ctx, err.Error(), nil)
		return
	}
	this.Success(ctx, "上传成功",
		gin.H{
			"url": result,
		},
	)
}
func (this *Publiccontroll) Uploadeditor(ctx *gin.Context) {

	//ctx.JSON(http.StatusOK, gin.H{
	//	"error": 0,
	//	"url":   "../../uploads/file/9299961eab1ff4e85e870912f2abb560_20240529172250.jpg",
	//})
	//return
	_, header, err := ctx.Request.FormFile("imgFile")
	contentType := header.Header.Get("Content-Type")
	allowedTypes := []string{"application/pdf", "image/jpeg", "image/png", "image/gif"}
	isAllowed := false
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		this.Error(ctx, "上传文件错误", nil)
		return
	}
	types := ctx.Query("type")
	if err != nil {
		this.Error(ctx, "上传文件错误", nil)
		return
	}
	var service service.FileUploadAndDownloadService
	result, err := service.UploadFile(header, "0", types)
	if err != nil {
		this.Error(ctx, err.Error(), nil)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"error": 0,
		"url":   "../../" + result,
	})
	return
}
