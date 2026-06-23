package config

type Email struct {
	Smtp          string `mapstructure:"Smtp" json:"smtp" yaml:"smtp"`
	EmailUserName string `mapstructure:"emailusername" json:"emailusername" yaml:"emailusername"`
	EmailPwd      string `mapstructure:"emailpwd" json:"emailpwd" yaml:"emailpwd"`
	From          string `mapstructure:"from" json:"from" yaml:"from"`
	Port          int    `mapstructure:"port" json:"port" yaml:"port"`
	Subject       string `mapstructure:"subject" json:"subject" yaml:"subject"`
}
