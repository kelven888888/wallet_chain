package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"time"
	"wallet_chain.com/admin/model"
	"wallet_chain.com/admin/service"
	"wallet_chain.com/global"
	"wallet_chain.com/utils"
)

func Captcha(ctx *gin.Context) {
	var service service.Captcha
	var captcha = service.Captcha(ctx.ClientIP())

	utils.Success(ctx, "成功", captcha)

}
func Upload(ctx *gin.Context) {

	//ctx.JSON(http.StatusOK, gin.H{
	//	"error": 0,
	//	"url":   "../../uploads/file/9299961eab1ff4e85e870912f2abb560_20240529172250.jpg",
	//})
	//return
	_, header, err := ctx.Request.FormFile("file")
	if err != nil {
		utils.Fail(ctx, err.Error(), "")
	}
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
		utils.Fail(ctx, "上传文件错误", "")

		return
	}

	types := ctx.Query("type")
	if err != nil {

		utils.Fail(ctx, "上传文件错误", "")
		return
	}
	var service service.FileUploadAndDownloadService
	result, err := service.UploadFile(header, "0", types)
	if err != nil {
		utils.Fail(ctx, "失败", "")
		return
	}
	data := map[string]string{

		"url": fmt.Sprintf(global.SHOP_CONFIG.System.WebApiURL) + result,
	}
	utils.Success(ctx, "成功", data)
}
func SendCode(ctx *gin.Context) {
	var params model.AccountCheckCode
	if err := ctx.ShouldBind(&params); err != nil {
		global.SHOP_LOG.Error(err.Error())
		utils.Fail(ctx, "参数错误", nil)
		return
	}
	length := len(params.Name)
	if length > 30 {
		utils.Fail(ctx, "参数错误", nil)
		return
	}
	key := fmt.Sprintf("code_%s", params.Name)
	result, err := global.SHOP_REDIS.Get(ctx, key).Result()

	if result != "" {
		utils.Fail(ctx, "操作频繁,请稍后再试", "")
		return
	}
	if params.Name == "" {
		utils.Fail(ctx, "请输入账号", "")
		return
	}
	if !utils.IsValidMalaysiaPhone(params.Name) && !utils.ValidateEmail(params.Name) {
		utils.Fail(ctx, "请输入正确的电话号码或邮箱", "")
		return
	}
	if params.Type == 1 {
		var user model.User
		global.SHOP_DB.Where("username=?", params.Name).Find(&user)
		if user.Id > 0 {
			utils.Fail(ctx, "用户已存在", "")
			return
		}
	}
	if params.Type == 3 || params.Type == 5 {
		var user model.User
		global.SHOP_DB.Where("username=?", params.Name).Find(&user)
		if user.Id == 0 {
			utils.Fail(ctx, "用户不存在", "")
			return
		}
	}

	language, _ := ctx.Get("Language")
	params.Language, _ = language.(string)
	var captcha string
	if global.SHOP_CONFIG.System.Env == "debug" {
		captcha = "123456"
	} else {
		captcha = strconv.Itoa(rand.Intn(900000) + 100000)
	}
	params.Captcha = captcha
	seconds := time.Now().Unix()
	params.CreateTime = time.Now()

	// 转换为int类型
	secondsInt := int(seconds) + 60
	params.CreateMap = secondsInt

	global.SHOP_DB.Save(&params)
	err = global.SHOP_REDIS.Set(ctx, key, params.Captcha, time.Second*60).Err()
	if err != nil {
		global.SHOP_LOG.Error(err.Error())
	}

	utils.Success(ctx, "成功", "")
}

// Register 注册
func Register(ctx *gin.Context) {
	DB := global.SHOP_DB.Begin()
	// 获取参数
	var params model.User
	if err := ctx.ShouldBind(&params); err != nil {
		global.SHOP_LOG.Error(err.Error())
		utils.Fail(ctx, "参数错误:", nil)
		return
	}

	//// 测试字符串
	//testStrings := []string{
	//	"Password123",        // 正确，长度9，包含大小写字母和数字
	//	"password",           // 错误，长度不足8位
	//	"PASSWORD1234567",    // 正确，长度9，包含大小写字母和数字
	//	"password1234567890", // 正确，长度10，包含大小写字母和数字
	//	"pass",               // 错误，长度不足8位
	//	"PassWord123456",     // 错误，长度超过16位
	//}
	//
	//// 遍历测试字符串，并使用正则表达式进行匹配
	//for _, str := range testStrings {
	//	if utils.IsValidPasswd(str) {
	//		fmt.Printf("\"%s\" is a valid password.\n", str)
	//	} else {
	//		fmt.Printf("\"%s\" is NOT a valid password.\n", str)
	//	}
	//}

	if !utils.IsValidPasswd(params.Password) {
		utils.Fail(ctx, "密码必须8到16个字符,包含大小写数字及字母", nil)
		return
	}
	if !utils.IsValidTradePasswd(params.TradePassword) {
		utils.Fail(ctx, "资金密码必须6位纯数字", nil)
		return
	}

	if params.Captcha == "" {
		utils.Fail(ctx, "验证码不能为空", nil)
		return
	}
	key := fmt.Sprintf("code_%s", params.Username)
	result, _ := global.SHOP_REDIS.Get(ctx, key).Result()
	if result == "" {
		utils.Fail(ctx, "验证码已过期", nil)
		return
	}
	if params.Captcha != result {
		utils.Fail(ctx, "验证码错误", nil)
		return
	}
	if params.Password != params.ConfirmPassword {
		utils.Fail(ctx, "密码与确认密码不一致", nil)
		return
	}
	if params.TradePassword != params.ConfirmTradePassword {
		utils.Fail(ctx, "资金密码与确认资金密码不一致", nil)
		return
	}
	var path_id = "0,"
	if params.InviteCode != "" {
		var puser model.User
		pid := utils.GetuidfromiCode(params.InviteCode)
		global.SHOP_DB.Where("id=? ", pid).Find(&puser)
		if puser.Id == 0 {
			utils.Fail(ctx, "邀请码错误", nil)
			return
		}
		puser.InviteCount = puser.InviteCount + 1
		global.SHOP_DB.Updates(&puser)
		path_id = puser.PathId
		params.Pid = pid

	}
	var exuser model.User
	global.SHOP_DB.Model(model.User{}).Where("username=?", params.Username).First(&exuser)
	if exuser.Id != 0 {
		utils.Fail(ctx, "用户已存在", nil)
		return
	}
	//if len(params.Username) == 0 {
	//	params.Username = utils.RandomString(10)
	//}
	//
	//if len(params.Usernick) == 0 {
	//	params.Usernick = utils.RandomString(10)
	//}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Fail(ctx, "加密错误", nil)
		return
	}

	params.Password = string(hashPassword)
	hashtradePassword, err := bcrypt.GenerateFromPassword([]byte(params.TradePassword), bcrypt.DefaultCost)
	if err != nil {
		utils.Fail(ctx, "加密错误", nil)
		return
	}

	params.TradePassword = string(hashtradePassword)
	language, _ := ctx.Get("Language")
	params.Language, _ = language.(string)
	params.DateJoined = time.Now()
	params.InviteCode = ""
	if utils.IsValidMalaysiaPhone(params.Username) {
		lens := len(params.AreaCode)
		params.Phone = params.Username[lens:]
	}
	if utils.ValidateEmail(params.Username) {
		params.Email = params.Username
	}
	//注册赠送一次盲盒次数
	params.BlindBoxNum = 1

	// uuid
	//params.UserID = utils.UUID()
	// 创建
	err = DB.Save(&params).Error
	if err != nil {
		DB.Rollback()
		global.SHOP_LOG.Error(err.Error())
		utils.Fail(ctx, "失败", nil)
		return
	}
	userID := params.Id

	// 1. 生成邀请码
	inviteCode := utils.BuildInviteCode(userID)

	if err != nil {
		utils.Fail(ctx, "失败", nil)
		global.SHOP_LOG.Error(err.Error())
	}
	params.InviteCode = inviteCode
	params.PathId = path_id + strconv.Itoa(userID) + ","
	err = DB.Save(&params).Error
	if err != nil {
		DB.Rollback()
		global.SHOP_LOG.Error(err.Error())
		utils.Fail(ctx, "失败", nil)
		return
	}
	var accountfund model.AccountFunds
	accountfund.Username = params.Username
	accountfund.LockFunds = 0
	accountfund.AvaFunds = decimal.NewFromInt(0)
	accountfund.CreateTime = time.Now()
	accountfund.Uid = params.Id
	accountfund.Points = decimal.NewFromInt(0)
	err = DB.Save(&accountfund).Error
	if err != nil {
		DB.Rollback()
		global.SHOP_LOG.Error(err.Error())
		utils.Fail(ctx, "失败", nil)
		return
	}

	// uuid
	//params.UserID = utils.UUID()
	// 创建
	if err != nil {
		utils.Fail(ctx, "失败", nil)
		return
	}
	DB.Commit()
	//var modelcode model.AccountCheckCode
	//global.SHOP_DB.Where("name=?", params.Username).Delete(&modelcode)
	global.SHOP_REDIS.Del(ctx, key)
	//返回结果
	utils.Success(ctx, "成功", nil)
}
func Login(ctx *gin.Context) {
	// 初始化数据库句柄
	DB := global.SHOP_DB
	// 定义使用模型
	var params model.User
	var user model.User
	// 绑定获取请求参数
	if err := ctx.ShouldBind(&params); err != nil {
		utils.Fail(ctx, "参数错误", nil)
		return
	}
	if len(params.Username) == 0 {
		utils.Fail(ctx, "用户名不能为空", nil)
		return
	}
	if len(params.Password) == 0 {
		utils.Fail(ctx, "密码不能为空", nil)
		return
	}
	if len(params.Password) < 6 {
		utils.Fail(ctx, "密码不能小于6位！", nil)
		return
	}
	// 获取用户
	DB.Where("username = ?", params.Username).First(&user)
	if user.Id == 0 {
		utils.Fail(ctx, "该用户未注册", nil)
		return
	}
	if *user.IsActive == 0 {
		utils.Fail(ctx, "该用户已禁用", nil)
		return
	}
	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		utils.Fail(ctx, "密码错误", nil)
		return
	}
	// 生成token
	token, tokenErr := utils.ReleaseToken(strconv.Itoa(user.Id))
	if tokenErr != nil {
		utils.Fail(ctx, "生成token失败", nil)
		return
	}
	key := fmt.Sprintf("login_%d", user.Id)
	dr, err := utils.ParseDuration(global.SHOP_CONFIG.JWT.ExpiresTime)
	if err != nil {
		panic(err)
	}
	global.SHOP_REDIS.Set(ctx, key, token, dr)
	// 获取 本机真实IP
	ip, _ := utils.ExternalIp()
	//user.LoginIp = ip.String()
	//user.LastLogin = time.Now()
	// 更新
	resultErr := DB.Model(model.User{}).Where("id=?", user.Id).Updates(model.User{
		LoginIp:   ip.String(),
		LastLogin: time.Now(),
	}).Error

	if resultErr != nil {
		utils.Fail(ctx, "登录失败", nil)
		return
	}
	//返回结果
	utils.Success(ctx, "登录成功", gin.H{
		"token": token,
	})
}
func Findpwd(ctx *gin.Context) {
	var params struct {
		ConfirmPassword string `json:"confirm_password" form:"confirm_password" `
		Password        string `json:"password" form:"password"  `
		OldPassword     string `json:"old_password" form:"old_password"  `
		UserName        string `json:"username" form:"username"  `
		Captcha         string `json:"captcha" form:"captcha"  `
	}
	if err := ctx.ShouldBind(&params); err != nil {
		utils.Fail(ctx, "参数错误", nil)
		return
	}
	var user model.User
	DB := global.SHOP_DB

	if params.Password != params.ConfirmPassword {
		utils.Fail(ctx, "密码与确认密码不一致", nil)
		return
	}
	if !utils.IsValidPasswd(params.Password) {
		utils.Fail(ctx, "密码必须8到16个字符,包含大小写数字及字母", nil)
		return
	}
	if params.Captcha == "" {
		utils.Fail(ctx, "验证码不能为空", nil)
		return
	}
	key := fmt.Sprintf("code_%s", params.UserName)
	result, _ := global.SHOP_REDIS.Get(ctx, key).Result()
	if result == "" {
		utils.Fail(ctx, "验证码已过期", nil)
		return
	}
	if params.Captcha != result {
		utils.Fail(ctx, "验证码错误", nil)
		return
	}
	err := DB.Where("username = ?", params.UserName).First(&user).Error
	if err != nil {
		utils.Fail(ctx, "用户不存在", nil)
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Fail(ctx, "加密错误", nil)
		return
	}
	user.Password = string(hashPassword)
	resErr := DB.Save(&user).Error
	if resErr != nil {
		utils.Fail(ctx, "失败", nil)
		return
	}
	global.SHOP_REDIS.Del(ctx, key)
	utils.Success(ctx, "成功", nil)
	return
}

func Agreement(ctx *gin.Context) {
	language, _ := ctx.Get("Language")
	Languages, _ := language.(string)
	var req model.Agreement
	if err := ctx.ShouldBind(&req); err != nil {
		utils.Fail(ctx, "参数错误", nil)
		return
	}
	if language == "" {
		Languages = global.SHOP_CONFIG.System.Language
	}
	// 获取参数
	var agree model.Agreement
	global.SHOP_DB.Where("language=? and `key`=?", Languages, req.Key).Find(&agree)
	if agree.Id == 0 {
		utils.Fail(ctx, "没有记录", nil)
		return
	}
	utils.Success(ctx, "成功", agree)
	return

}
