package service

import (
	"errors"
	"wallet_chain.com/admin/model"
	"wallet_chain.com/admin/model/common/request"
	"wallet_chain.com/global"
	"wallet_chain.com/utils"
)

type AdminServer struct {
}

func (this *AdminServer) Login(loginreq request.Login) (*model.Admin, error) {
	var adminuser model.Admin
	err := global.SHOP_DB.Where("account = ?", loginreq.Username).First(&adminuser).Error
	if err != nil {
		return nil, errors.New("密码错误")
	}
	if adminuser.Id == 0 && err != nil {
		return nil, errors.New("账号不存在")
	}
	//fmt.Println(utils.EncryptPassworld(utils.MD5V(loginreq.Password)))
	//fmt.Println(adminuser.Password)
	if adminuser.Password != utils.EncryptPassworld(utils.MD5V(loginreq.Password)) {

		return nil, errors.New("密码错误")
	}

	if *adminuser.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}
	return &adminuser, nil

}
