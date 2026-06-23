package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"wallet_chain.com/admin/model/common/request"
	"wallet_chain.com/admin/service"
	"wallet_chain.com/global"
)

type Adminlogcontroll struct {
	BaseController
}

func (this *Adminlogcontroll) Index(ctx *gin.Context) {
	var req request.PageInfo
	var service service.AdminlogServer
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

	adminloglist, count := service.GetAll(req)

	Search := map[string]interface{}{
		"page":  p,
		"limit": size,
	}
	ctx.HTML(http.StatusOK, "adminlog_index.html", gin.H{
		"status": "200",
		"List":   adminloglist,
		"Count":  count,
		"Search": Search,
	})

}
func (this *Adminlogcontroll) Delete(ctx *gin.Context) {
	var req request.GetById
	err := ctx.ShouldBind(&req)
	if err != nil {
		this.ErrorHtml(ctx, err.Error())
		return
	}
	var service service.AdminlogServer
	err = service.Delete(req)
	if err != nil {
		this.ErrorHtml(ctx, err.Error())
		return
	}
	this.Success(ctx, "删除成功")

}
func (this *Adminlogcontroll) Deletebatch(ctx *gin.Context) {
	var req request.IdsReq
	err := ctx.ShouldBind(&req)
	if err != nil {
		this.ErrorHtml(ctx, err.Error())
		return
	}
	var service service.AdminlogServer
	err = service.Deletebatch(req)
	if err != nil {
		this.ErrorHtml(ctx, err.Error())
		return
	}
	this.Success(ctx, "删除成功")
}
