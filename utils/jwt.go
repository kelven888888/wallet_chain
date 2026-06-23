package utils

import (
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"time"
	"wallet_chain.com/global"
)

type JWT struct {
	SigningKey  string `mapstructure:"signing-key" json:"signing-key" yaml:"signing-key"`    // jwt签名
	ExpiresTime string `mapstructure:"expires-time" json:"expires-time" yaml:"expires-time"` // 过期时间
	BufferTime  string `mapstructure:"buffer-time" json:"buffer-time" yaml:"buffer-time"`    // 缓冲时间
	Issuer      string `mapstructure:"issuer" json:"issuer" yaml:"issuer"`                   // 签发者
}

var jwtKey = []byte(viper.GetString("server.jwtKey"))

type Claims struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

// ReleaseToken 生成jwt token
func ReleaseToken(userId string) (string, error) {
	ip, _ := ExternalIp()
	dr, err := ParseDuration(global.SHOP_CONFIG.JWT.ExpiresTime)
	if err != nil {
		panic(err)
	}

	expireTime := time.Now().Add(dr)

	claims := &Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    ip.String(),
			Subject:   "token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

// ParseToken 解析jwt token
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}
