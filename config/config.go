package config

type Server struct {
	Local   Local   `mapstructure:"local" json:"local" yaml:"local"`
	System  System  `mapstructure:"system" json:"system" yaml:"system"`
	Mysql   Mysql   `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Zap     Zap     `mapstructure:"zap" json:"zap" yaml:"zap"`
	Redis   Redis   `mapstructure:"redis" json:"redis" yaml:"redis"`
	JWT     JWT     `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Captcha Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	Sms     Sms     `mapstructure:"sms" json:"sms" yaml:"sms"`
	Email   Email   `mapstructure:"email" json:"email" yaml:"email"`
	Wallet  Wallet  `mapstructure:"wallet" json:"wallet" yaml:"wallet"`
}
