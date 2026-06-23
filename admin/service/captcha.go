package service

import (
	"time"

	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"wallet_chain.com/global"
)

// 当开启多服务器部署时，替换下面的配置，使用redis共享存储验证码
// var store = captcha.NewDefaultRedisStore()
var store = base64Captcha.DefaultMemStore

type Captcha struct{}

// Captcha
// @Tags      Base
// @Summary   生成验证码
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=systemRes.SysCaptchaResponse,msg=string}  "生成验证码,返回包括随机数id,base64,验证码长度,是否开启验证码"
// @Router    /base/captcha [post]
func (b *Captcha) Captcha(key string) interface{} {
	// 判断验证码是否开启
	openCaptcha := global.SHOP_CONFIG.Captcha.OpenCaptcha               // 是否开启防爆次数
	openCaptchaTimeOut := global.SHOP_CONFIG.Captcha.OpenCaptchaTimeOut // 缓存超时时间

	v, ok := global.BlackCache.Get(key)
	if !ok {
		global.BlackCache.Set(key, 1, time.Second*time.Duration(openCaptchaTimeOut))
	}

	var oc bool
	if openCaptcha == 0 || openCaptcha < interfaceToInt(v) {
		oc = true
	}
	// 字符,公式,验证码配置
	// 生成默认数字的driver
	driver := base64Captcha.NewDriverDigit(global.SHOP_CONFIG.Captcha.ImgHeight, global.SHOP_CONFIG.Captcha.ImgWidth, global.SHOP_CONFIG.Captcha.KeyLong, 0.9, 80)
	// cp := base64Captcha.NewCaptcha(driver, store.UseWithCtx(c))   // v8下使用redis
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err := cp.Generate()
	if err != nil {
		global.SHOP_LOG.Error("验证码获取失败!", zap.Error(err))
		return err.Error()
	}
	return map[string]interface{}{
		"CaptchaId":     id,
		"PicPath":       b64s,
		"CaptchaLength": global.SHOP_CONFIG.Captcha.KeyLong,
		"OpenCaptcha":   oc,
	}
}

// 类型转换
func interfaceToInt(v interface{}) (i int) {
	switch v := v.(type) {
	case int:
		i = v
	default:
		i = 0
	}
	return
}
