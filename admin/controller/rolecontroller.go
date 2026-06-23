package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"wallet_chain.com/admin/model"
	"wallet_chain.com/admin/model/common/request"
	"wallet_chain.com/admin/service"
	"wallet_chain.com/global"
)

type RoleController struct {
	BaseController
}

func (this *RoleController) Index(ctx *gin.Context) {

	var req request.PageInfo
	var service service.Role
	err := ctx.ShouldBind(&req)
	if err != nil {
		this.ErrorHtml(ctx, err.Error())
		return
	}

	p := req.Page
	if p == 0 {
		p = 1
	}

	size, _ := strconv.Atoi(global.SHOP_CONFIG.System.PageSize)

	req.Count = true
	req.Limit = size
	req.Offset = (p - 1) * size

	rolelist, count := service.GetAll(req)

	Search := map[string]interface{}{
		"page":  p,
		"limit": size,
		"kw":    req.Keyword,
	}
	ctx.HTML(http.StatusOK, "role_index.html", gin.H{
		"status": "200",
		"List":   rolelist,
		"Count":  count,
		"Search": Search,
	})

}
func (this *RoleController) Edit(ctx *gin.Context) {

	if ctx.Request.Method == "GET" {
		var id request.GetById

		err := ctx.ShouldBind(&id)
		if err != nil {
			this.ErrorHtml(ctx, err.Error())
		}
		var server service.Role
		role, err := server.GetRole(id.Uint32())
		if err != nil {
			this.ErrorHtml(ctx, err.Error())
		}

		// 查询权限列表
		roles := server.GetRoles(true)

		options := this.getSelectTree(roles, role.Pid, 0)
		ctx.HTML(http.StatusOK, "role_form.html", gin.H{
			"status":   "200",
			"role":     role,
			"options":  options,
			"IsUpdate": true,
		})
	} else {
		var rolemodel model.Role
		var server service.Role
		err := ctx.ShouldBind(&rolemodel)
		fmt.Println(rolemodel)
		if err != nil {
			this.Error(ctx, err.Error())
			return
		}
		err = server.Save(&rolemodel)
		if err != nil {
			this.Error(ctx, err.Error())
			return
		}
		service.Getmenu()
		this.Success(ctx, "更新成功")

	}
}
func (this *RoleController) Add(ctx *gin.Context) {
	if ctx.Request.Method == "GET" {
		var server service.Role

		// 查询权限列表
		roles := server.GetRoles(true)

		options := this.getSelectTree(roles, 0, 0)
		ctx.HTML(http.StatusOK, "role_form.html", gin.H{
			"status": "200",

			"options":  options,
			"IsUpdate": false,
		})
	} else {
		var rolemodel model.Role
		var server service.Role
		err := ctx.ShouldBind(&rolemodel)

		if err != nil {
			this.Error(ctx, err.Error())
			return
		}
		err = server.Save(&rolemodel)
		if err != nil {
			this.Error(ctx, err.Error())
			return
		}
		service.Getmenu()
		this.Success(ctx, "成功")
	}

}
func (this *RoleController) Delete(ctx *gin.Context) {

	var req request.GetById
	err := ctx.ShouldBind(&req)
	if err != nil {
		this.Error(ctx, err.Error())
		return
	}
	var server service.Role
	err = server.Delete(req.Uint32())
	if err != nil {
		this.Error(ctx, err.Error())
		return
	}
	service.Getmenu()
	this.Success(ctx, "删除成功")

}

// 格式化为select表单树
func (this *RoleController) getSelectTree(roles []*model.Role, pid uint32, level int) string {
	html := ""
	for _, role := range roles {
		sel := ""
		if role.Id == pid {
			sel = " selected"
		}
		html += fmt.Sprintf(`<option value="%d" %s>%s%s</option>`,
			role.Id,
			sel,
			strings.Repeat("-", level*2),
			role.Name,
		)

		if role.Id != pid && len(role.Child) > 0 {
			html += this.getSelectTree(role.Child, pid, level+1)
		}
	}

	return html
}
