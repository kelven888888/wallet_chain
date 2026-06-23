package model

type Config struct {
	Id    uint
	Key   string `comment:"健值名"`
	Value string `comment:"健值内容"`
	ModelTime
}

func (Config) TableName() string {
	return "nov_config"
}
