package utils

const (
	Rechargetype = iota + 1
	Adminlogtype
	Withdrayapplytype
	Withdrayapplyrefuetype
)

// 1:注册2修改密码,3忘记密码 4认证
const (
	Reg_code = iota + 1
	Change_pwd_code
	Forgot_pwd_code
	Auth_code
	Recharge_Msg
	Whthdraw_Success_Msg
	Whthdraw_Fail_Msg
	Fund_Redemption_Successful
	Quan_Redemption_Successful
	Trade_Account_Pass
	Trade_Account_Reject
	BankRecharge_Msg_Wait
	BankRecharge_Msg_Success
	Review_Passed
	Review_Rejected
)
