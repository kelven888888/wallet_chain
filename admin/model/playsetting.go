package model

type MPlayConfig struct {
	Model
	RateArr      string `gorm:"column:rate_arr" comment:"概率配置"`                                     // 概率配置
	RewardArr    string `gorm:"column:reward_arr" comment:"概率配置"`                                   // 奖励配置
	SingelStatus int    `gorm:"column:singel_status" types:"radio" text:"关闭,启用" range:"0,1"`        // 0关1开
	DoubleStatus int    `gorm:"column:double_status" types:"radio" text:"关闭,启用" range:"0,1"`        // 0关1开
	Price        int    `gorm:"column:price" comment:"参加价格"`                                        // 参加价格
	PointArr     string `gorm:"column:point_arr" `                                                  // 积分配置
	GoodsIds     string `gorm:"column:goods_ids" `                                                  // 产品配置
	Type         int    `gorm:"column:type" types:"radio" text:"高爆,保底,魔王,PK,5随机" range:"1,2,3,4,5"` // 1高爆2保底3魔王4pk5随机
	Name         string `gorm:"column:name" comment:"名称"`
}

func (m *MPlayConfig) TableName() string {
	return "play_config"
}
