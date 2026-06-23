package controller

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"wallet_chain.com/admin/model"
	"wallet_chain.com/admin/model/common/request"
	"wallet_chain.com/admin/model/common/response"
	"wallet_chain.com/admin/service"
	"wallet_chain.com/global"
	"wallet_chain.com/utils"
)

type Admincontroll struct {
	BaseController
}

func (this *Admincontroll) Adminlist(ctx *gin.Context) {
	var group []model.Group
	global.SHOP_DB.Find(&group)
	ctx.HTML(http.StatusOK, "admin_list.html", gin.H{

		"group": group,
	})

}
func (this *Admincontroll) Getmenu(ctx *gin.Context) {
	server := service.GroupServer{}
	session := sessions.Default(ctx)

	groupId := session.Get("groupId").(uint)

	Menus, err := server.GetMenus(groupId)

	fmt.Printf("%+v", Menus)
	if err != nil {
		fmt.Printf("%+v", err)
	}
	this.Success(ctx, "ok", Menus)
	ctx.HTML(http.StatusOK, "admin_list.html", Menus)

}

func (this *Admincontroll) Add(ctx *gin.Context) {
	method := ctx.Request.Method

	if method == "GET" {
		//var adminUser model.Admin
		var group []model.Group
		//global.SHOP_DB.First(&adminUser)
		global.SHOP_DB.Find(&group)
		ctx.HTML(http.StatusOK, "admin_form.html", gin.H{

			//"admininfo": adminUser,
			"group": group,
		})
	} else {
		var adminUser model.Admin
		if err := ctx.ShouldBindJSON(&adminUser); err == nil {

			fmt.Printf("login-request:%+v\n", adminUser)
			//global.SHOP_DB.Where("username = ?", adminUser.Username).First(&adminUser)

			adminUser.Password = utils.EncryptPassworld(utils.MD5V(adminUser.Password))

			err = global.SHOP_DB.Save(&adminUser).Error
			if err != nil {
				this.Error(ctx, err.Error())
				return
			}

			this.Success(ctx, "添加成功")
		} else {
			this.Error(ctx, err.Error())
		}
	}

}
func (this *Admincontroll) Delete(ctx *gin.Context) {
	var requid request.GetByUserId
	var err = ctx.ShouldBindJSON(&requid)
	if err != nil {
		this.Error(ctx, err.Error())
		return

	}
	if requid.Uint() == 1 {
		this.Error(ctx, "超级管理员1不允许删除")
		return
	}
	var user model.Admin
	err = global.SHOP_DB.Where("id=?", requid.Uint()).Delete(&user).Error
	if err != nil {
		this.Error(ctx, err.Error())
		return
	}
	this.Success(ctx, "删除成功", requid)

}
func (this *Admincontroll) Deletebatch(ctx *gin.Context) {
	var requids request.IdsReq
	var err = ctx.ShouldBindJSON(&requids)
	if err != nil {
		this.Error(ctx, err.Error())
		return

	}
	var user model.Admin
	fmt.Println(requids.Ids)
	if utils.InArray(1, requids.Ids) {
		this.Error(ctx, "超级管理员1不允许删除")
		return
	}
	err = global.SHOP_DB.Where("id in?", requids.Ids).Delete(&user).Error
	if err != nil {
		this.Error(ctx, err.Error())
		return
	}
	this.Success(ctx, "删除成功", requids)

}
func (this *Admincontroll) Edit(ctx *gin.Context) {
	method := ctx.Request.Method

	if method == "GET" {
		var adminUser model.Admin
		var group []model.Group
		if err := ctx.ShouldBind(&adminUser); err == nil {

			fmt.Printf("login-request:%+v\n", adminUser)
			global.SHOP_DB.First(&adminUser)
			global.SHOP_DB.Find(&group)
			fmt.Println(adminUser)
			ctx.HTML(http.StatusOK, "admin_form.html", gin.H{

				"admininfo": adminUser,
				"group":     group,
				"IsUpdate":  true,
			})
		} else {
			this.Error(ctx, "查询错误", err.Error())
		}
	} else {
		var adminUser model.Admin
		if err := ctx.ShouldBindJSON(&adminUser); err == nil {

			fmt.Printf("login-request:%+v\n", adminUser)
			//global.SHOP_DB.Where("username = ?", adminUser.Username).First(&adminUser)
			//id := adminUser.Id
			adminUser.Password = utils.EncryptPassworld(utils.MD5V(adminUser.Password))
			err = global.SHOP_DB.Where("id"+
				" = ?", adminUser.Id).First(&model.Admin{}).Updates(&adminUser).Error

			if err != nil {

				this.Error(ctx, "失败", err.Error())
			}
			//service.RunPublisher(ctx, "shop_message", "修改"+fmt.Sprintf("%v", id)+"成功")
			this.Success(ctx, "修改成功", adminUser)
		} else {
			fmt.Printf(err.Error())
			this.Error(ctx, "修改失败", err.Error())
		}
	}

}

// 退出登录
func (this Admincontroll) Logout(c *gin.Context) {

	session := sessions.Default(c)
	session.Clear() // 清除所有数据
	if err := session.Save(); err != nil {
		this.Error(c, "退出失败")
		return
	}

	//c.JSON(http.StatusOK, gin.H{
	//	"code": 200,
	//})
	this.Success(c, "退出成功")
}
func (this *Admincontroll) Getlist(ctx *gin.Context) {
	var adminUser []model.Admin

	var pageInfo request.PageInfo
	err := ctx.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	var count int64 = 0
	//err = global.SHOP_DB.Limit(pageInfo.Limit).Offset((pageInfo.Page - 1) * pageInfo.Limit).Order("user_id ASC").Find(&adminUser).Error
	query := global.SHOP_DB.Table("nov_admin").Select("nov_admin.*,nov_group.name as groupname").Joins("left join nov_group on nov_admin.group_id = nov_group.id")
	if pageInfo.Account != "" {
		query.Where("nov_admin.account=?", pageInfo.Account)
	}
	if strconv.Itoa(pageInfo.GroupId) != "" && pageInfo.GroupId != 0 {
		query.Where("nov_group.id =?", pageInfo.GroupId)
	}
	query.Where("nov_admin.deleted_at is  null")
	query.Count(&count)
	err = query.Limit(pageInfo.Limit).Offset((pageInfo.Page - 1) * pageInfo.Limit).Order(" nov_admin.id ASC").Find(&adminUser).Error
	if err != nil {
		this.Error(ctx, "查询错误", err.Error())
		fmt.Println(err)
		return
	}
	//global.SHOP_DB.Model(model.Admin{}).Count(&count)
	if err != nil {
		fmt.Println(err)
		this.Error(ctx, "没有数据", err.Error())
		return
	}
	this.Success(ctx, "查询成功", adminUser, count)
	//ctx.String(http.StatusOK, str)
	//this.Success(ctx, "退出成功", str)

}
func (this *Admincontroll) Changepwd(ctx *gin.Context) {
	method := ctx.Request.Method

	if method == "GET" {
		var adminUser model.Admin

		if err := ctx.ShouldBind(&adminUser); err == nil {

			session := sessions.Default(ctx)

			adminId := session.Get("adminId")
			userId := adminId.(uint)
			global.SHOP_DB.Where("id=?", userId).First(&adminUser)

			ctx.HTML(http.StatusOK, "admin_changepwd.html", gin.H{

				"admininfo": adminUser,

				"IsUpdate": true,
			})
		} else {
			this.Error(ctx, "查询错误", err.Error())
		}
	} else {
		var adminUser model.Admin
		if err := ctx.ShouldBindJSON(&adminUser); err == nil {

			session := sessions.Default(ctx)

			adminId := session.Get("adminId")
			userId := adminId.(uint)
			//global.SHOP_DB.Where("username = ?", adminUser.Username).First(&adminUser)
			//id := adminUser.Id
			if adminUser.Password == "" {
				this.Error(ctx, "修改失败,密码为空")
				return
			}
			adminUser.Password = utils.EncryptPassworld(utils.MD5V(adminUser.Password))
			err = global.SHOP_DB.Where("id"+
				" = ?", userId).First(&model.Admin{}).Updates(&adminUser).Error

			if err != nil {

				this.Error(ctx, "失败", err.Error())
				return
			}
			//service.RunPublisher(ctx, "shop_message", "修改"+fmt.Sprintf("%v", id)+"成功")
			session.Clear() // 清除所有数据
			if err := session.Save(); err != nil {
				println(err)

			}
			this.Success(ctx, "修改成功请重新登录", adminUser)

		} else {
			fmt.Println(err.Error())
			this.Error(ctx, "修改失败", err.Error())
			return
		}
	}

}
func (this *Admincontroll) CleanCache(ctx *gin.Context) {

	session := sessions.Default(ctx)

	groupId := session.Get("groupId").(uint)
	if groupId != 1 {
		this.ErrorHtml(ctx, fmt.Sprintf("只有超管有权限%d", groupId))
		return
	}

	action := ctx.Query("action")
	if action == "cleancache" {
		err := utils.BatchDeleteKeys(global.SHOP_REDIS, "showcase_*")
		if err != nil {
			this.ErrorHtml(ctx, "清除redis缓存失败", err.Error())
			return
		}
	} else {
		err := utils.BatchDeleteKeys(global.SHOP_REDIS, "session_*")
		if err != nil {
			this.ErrorHtml(ctx, "清除redis缓存失败", err.Error())
			return
		}
	}

	this.SuccessHtml(ctx, "成功")

}
