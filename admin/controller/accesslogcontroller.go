package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"wallet_chain.com/admin/model/common/request"
	"wallet_chain.com/admin/service"
	"wallet_chain.com/global"
)

type Acclogctr struct {
	Services service.Saccesslog
	BaseController
}

func (this *Acclogctr) Index(ctx *gin.Context) {

	var req request.PageInfo
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

	result, count := this.Services.GetAll(req)

	Search := map[string]interface{}{
		"page":   p,
		"limit":  size,
		"kw":     req.Keyword,
		"status": req.Status,
	}

	ctx.HTML(http.StatusOK, "accesslog_index.html", gin.H{
		"status": "200",
		"List":   result,
		"Count":  count,
		"Search": Search,
	})
}
