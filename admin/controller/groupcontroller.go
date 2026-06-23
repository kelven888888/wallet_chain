package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"wallet_chain.com/admin/model"
	"wallet_chain.com/admin/model/common/request"
	"wallet_chain.com/admin/service"
	"wallet_chain.com/global"
	"wallet_chain.com/utils"
)

type GroupController struct {
	BaseController
}

func (this *GroupController) Index(ctx *gin.Context) {
	var group []model.Group
	err := global.SHOP_DB.Find(&group).Error
	if err != nil {
		this.ErrorHtml(ctx, err.Error(), err.Error())
		return
	}

	count := len(group)
	ctx.HTML(http.StatusOK, "group_index.html", gin.H{
		"status": "200",
		"List":   group,
		"Count":  count,
	})
}
func (this *GroupController) Delete(ctx *gin.Context) {
	var req request.GetById
	err := ctx.ShouldBind(&req)
	if err != nil {
		this.Error(ctx, err.Error())
		return
	}
	var service service.GroupServer
	err = service.Delete(req.Uint32())
	if err != nil {
		this.Error(ctx, err.Error())
		return
	}
	this.Success(ctx, "删除成功")
}
func (this *GroupController) Edit(ctx *gin.Context) {
	if ctx.Request.Method == "GET" {
		var req request.GetById
		err := ctx.ShouldBind(&req)
		if err != nil {
			this.ErrorHtml(ctx, err.Error())
			return
		}
		var groupservice service.GroupServer
		group, err := groupservice.GetByID(&req)
		if err != nil {
			this.ErrorHtml(ctx, err.Error())
			return
		}

		var roleservice service.Role
		roles := roleservice.GetRoles(true)
		fmt.Println(group.GetRoleIds())
		rolehtml := this.GetRolesHtml(roles, group.GetRoleIds())
		ctx.HTML(http.StatusOK, "group_form.html", gin.H{
			"status":   "200",
			"Group":    group,
			"IsUpdate": true,
			"RoleHtml": rolehtml,
			//"PostUrl":  "admin.GroupController.Edit",
		})
	} else {
		var group model.Group
		err := ctx.ShouldBind(&group)
		if err != nil {
			this.Error(ctx, err.Error())
			return
		}
		group.RoleIds = group.GetRoleString()
		fmt.Println(group.RoleIds)

		//fmt.Printf("%+v", group)
		var groupservice service.GroupServer
		err = groupservice.Save(&group)
		if err != nil {
			this.Error(ctx, err.Error())
			return
		}
		this.Success(ctx, "更新成功")

	}
}
func (this *GroupController) Add(ctx *gin.Context) {
	if ctx.Request.Method == "GET" {
		var group model.Group
		var roleservice service.Role
		roles := roleservice.GetRoles(true)
		rolehtml := this.GetRolesHtml(roles, make([]string, 0))
		ctx.HTML(http.StatusOK, "group_form.html", gin.H{
			"status":   "200",
			"Group":    group,
			"IsUpdate": true,
			"RoleHtml": rolehtml,
			"action":   "add",
			//"PostUrl":  "admin.GroupController.add",
		})
	} else {
		var group model.Group
		err := ctx.ShouldBind(&group)
		if err != nil {
			this.Error(ctx, err.Error())
			return
		}
		group.RoleIds = group.GetRoleString()
		//fmt.Printf("%+v", group)
		var groupservice service.GroupServer
		err = groupservice.Save(&group)
		if err != nil {
			this.Error(ctx, err.Error())
			return
		}
		this.Success(ctx, "添加成功")
	}
}
func (this BaseController) GetRolesHtml(roles []*model.Role, selRoleIds []string) string {
	htmltpl := `
			<tr>
				<td><input class="roles-select" type="checkbox" value="%d" lay-skin="primary" %s title="%s"></td>
				<td>
					<div class="layui-input-block">
					%s
					%s
					</div>
				</td>
			</tr>
		`

	html := ""
	for _, role := range roles {
		childHtml, isDef := this.getRolesChildHtml(role.Child, role.Id, selRoleIds)
		checked := ""
		if isDef {
			checked = "checked"
		}
		shtml, _ := this.getRoleFormatHtml(role, role.Id, selRoleIds)
		html += fmt.Sprintf(htmltpl,
			role.Id,
			checked,
			role.Name,
			shtml,
			childHtml,
		)
	}

	return html
}

// child class html
func (this BaseController) getRolesChildHtml(roles []*model.Role, rootId uint32, selRoleIds []string) (string, bool) {
	html := ""
	isAllchecked := true
	for _, role := range roles {
		// 格式化返回复选框html
		shtml, isChecked := this.getRoleFormatHtml(role, rootId, selRoleIds)

		html += shtml

		if isChecked == false {
			isAllchecked = false
		}

		if len(role.Child) > 0 {
			shtml, isChecked := this.getRolesChildHtml(role.Child, rootId, selRoleIds)

			html += shtml

			if isChecked == false {
				isAllchecked = false
			}
		}
	}

	return html, isAllchecked
}

// 格式化返回复选框html
func (this BaseController) getRoleFormatHtml(role *model.Role, rootId uint32, selRoleIds []string) (string, bool) {
	sel := ""
	isChecked := false
	if (role.IsDefault == 1 && role.Id == 0) || utils.InArray(fmt.Sprint(role.Id), selRoleIds) {
		sel = "checked"
		isChecked = true
	}
	html := fmt.Sprintf(`<input class="role_ids" lay-skin="primary" name="role_ids[]" type="checkbox" value="%d" pid="%d" %s title="%s">`,
		role.Id,
		rootId,
		sel,
		role.Name,
	)

	return html, isChecked
}
