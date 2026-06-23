package config

type System struct {
	Env             string `mapstructure:"env" json:"env" yaml:"env"`
	Version         string `mapstructure:"version" json:"version" yaml:"version"`
	Addr            int    `mapstructure:"addr" json:"addr" yaml:"addr"`
	DbType          string `mapstructure:"db-type" json:"dbType" yaml:"db-type"`
	OssType         string `mapstructure:"oss-type" json:"ossType" yaml:"oss-type"`
	UseMultipoint   bool   `mapstructure:"use-multipoint" json:"useMultipoint" yaml:"use-multipoint"`
	SecretKey       string `mapstructure:"secret-key" json:"secret-key" yaml:"secret-key"`
	PageSize        string `mapstructure:"pagesize" json:"pagesize" yaml:"pagesize"`
	LimitTimeIP     int    `mapstructure:"iplimit-time" json:"iplimit-time" yaml:"iplimit-time"`
	LimitCountIP    int    `mapstructure:"iplimit-count" json:"iplimit-count" yaml:"iplimit-count"`
	WebApiURL       string `mapstructure:"web_api_url" json:"web_api_url" yaml:"web_api_url"`
	AesKey          string `mapstructure:"aes_key" json:"aes_key" yaml:"aes_key"`
	AesIv           string `mapstructure:"aes_iv" json:"aes_iv" yaml:"aes_iv"`
	Ai_td_acc       string `mapstructure:"ai_td_acc" json:"ai_td_acc" yaml:"ai_td_acc"`
	POLYGON_API_KEY string `mapstructure:"polygon_api_key" json:"polygon_api_key" yaml:"polygon_api_key"`
	Logo_path       string `mapstructure:"logo_path" json:"logo_path" yaml:"logo_path"`
	Finn_API_KEY    string `mapstructure:"finn_api_key" json:"finn_api_key" yaml:"finn_api_key"`
	Language        string `mapstructure:"language" json:"language" yaml:"language"`
	Language_Array  string `mapstructure:"language_array" json:"language_array" yaml:"language_array"`
}
