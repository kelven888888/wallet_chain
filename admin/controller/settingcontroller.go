package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"wallet_chain.com/admin/model"
	"wallet_chain.com/global"
)

type SettingController struct {
	BaseController
}

type SettingPoster struct {
	Website       map[string]string `json:"website"`
	EnableWebsite bool              `json:"enable_website"`
	App           map[string]string `json:"app"`
	EnableApp     bool              `json:"enable_app"`
}

func (SettingPoster) Key() string {
	return "poster"
}

func (h SettingController) Poster(ctx *gin.Context) {
	db := global.SHOP_DB
	req := &SettingPoster{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		h.ErrorHtml(ctx, err.Error())
		return
	}

	if ctx.Request.Method == "GET" {
		setting := &model.Setting{}
		err = db.Where("`key` = ?", req.Key()).First(setting).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			h.ErrorHtml(ctx, err.Error())
			return
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			b, _ := json.Marshal(req)
			setting = &model.Setting{
				Key:     req.Key(),
				Comment: "弹窗海报",
				Value:   string(b),
			}
			err = db.Create(setting).Error
			if err != nil {
				h.ErrorHtml(ctx, err.Error())
				return
			}
		} else {
			_ = json.Unmarshal([]byte(setting.Value), req)
		}

		var language []model.Language
		global.SHOP_DB.Model(model.Language{}).Find(&language)

		ctx.HTML(200, "setting_poster.html", gin.H{
			"poster":   req,
			"language": language,
		})

		return
	}
	if ctx.Request.Method == "POST" {
		b, _ := json.Marshal(req)
		err = db.Table("setting").Where("`key` = ?", req.Key()).Update("value", string(b)).Error
		if err != nil {
			h.Error(ctx, err.Error())
			return
		}

		h.Success(ctx)
		return
	}
}
