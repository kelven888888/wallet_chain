package config

type Wallet struct {
	Url          string  `mapstructure:"url" json:"url" yaml:"url"`
	Appkey       string  `mapstructure:"appkey" json:"appkey" yaml:"appkey"`
	Secretkey    string  `mapstructure:"secretkey" json:"secretkey" yaml:"secretkey"`
	CallbackURL  string  `mapstructure:"callback_url" json:"callback_url" yaml:"callback_url"`
	LimitAvaAddr string  `mapstructure:"limit_ava_addr" json:"limit_ava_address" yaml:"limit_ava_addr"`
	Pid          int64   `mapstructure:"pid" json:"pid" yaml:"pid"`
	Withdrawfee  float64 `mapstructure:"withdrawfee" json:"withdrawfee"`
}
