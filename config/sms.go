package config

type Sms struct {
	Url       string `mapstructure:"url" json:"url" yaml:"url"`
	Appkey    string `mapstructure:"appkey" json:"appkey" yaml:"appkey"`
	Secretkey string `mapstructure:"secretkey" json:"secretkey" yaml:"secretkey"`
	Appcode   string `mapstructure:"appcode" json:"appcode" yaml:"appcode"`
}
