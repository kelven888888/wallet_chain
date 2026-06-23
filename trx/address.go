package trx

import (
	"crypto/ecdsa"
	"errors"
	"github.com/smirkcat/hdwallet"
)

func SearchAccount(addr string) (*Account, error) {
	var ac *Account
	var err error

	ac, err = dbengine.GetAccountWithAddr(addr)
	return ac, err
}

func NewPrivateKey() (int, *ecdsa.PrivateKey, error) {
	if IsMulti {
		return 0, nil, errors.New("not suppot new addr is_multi true ")
	}
	index := dbengine.GetAccountMaxIndex() + 1
	ac, err := hdwallet.NewPrivateKeyIndex(index)
	return index, ac, err
}
