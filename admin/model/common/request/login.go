package request

type Login struct {
	Username  string `form:"username" json:"username" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
	Code      string `form:"code" json:"code" binding:"required"`
	CaptchaId string `form:"captchaid" json:"captchaid" binding:"required"`
}
