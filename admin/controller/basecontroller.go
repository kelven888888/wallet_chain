package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseController struct {
}

type ResponseData struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (con BaseController) Success(c *gin.Context, res ...interface{}) {
	var code int16 = 200
	if reslen := len(res); reslen == 1 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  res[0],
			"data": []string{},
		})

	} else if reslen == 2 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  res[0],
			"data": res[1],
		})
	} else if reslen == 3 {
		c.JSON(http.StatusOK, gin.H{
			"code":  code,
			"msg":   res[0],
			"data":  res[1],
			"count": res[2],
		})
	} else if reslen == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  "Success",
			"data": []string{},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  res[0],
			"data": []string{},
		})

	}

}

func (con BaseController) Error(c *gin.Context, res ...interface{}) {

	var code int16 = 400
	if reslen := len(res); reslen == 1 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  res[0],
			"data": []string{},
		})

	} else if reslen == 2 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  res[0],
			"data": res[1],
		})
	} else if reslen == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  "Error",
			"data": []string{},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  res[0],
			"data": []string{},
		})

	}
}
func (con BaseController) ErrorHtml(c *gin.Context, res ...interface{}) {

	c.HTML(http.StatusOK, "502.html", gin.H{
		"status": "200",
		"err":    res[0],
	})
}
func (con BaseController) SuccessHtml(c *gin.Context, res ...interface{}) {

	c.HTML(http.StatusOK, "sucess.html", gin.H{
		"status": "200",
		"err":    res[0],
	})
}
