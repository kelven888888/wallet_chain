package model

type Setting struct {
	Key     string `json:"key" gorm:"primary_key"`
	Comment string `json:"comment" gorm:"column:comment"`
	Value   string `json:"value" gorm:"column:value;"`
}

func (Setting) TableName() string {
	return "setting"
}
