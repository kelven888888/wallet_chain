package model

import (
	"time"
)

type StockPredictionDayList struct {
	Id         int
	Stock      string    `comment:"股票列表"`
	CreatedAt  time.Time `comment:"创建时间"`
	Type       int       `comment:"1中期 2高爆"`
	Status     string    `comment:"0 未处理1 已处理"`
	DateTime   time.Time `comment:"预测时间年月日 每天一天记录"`
	ChangeVal  float64   `comment:"涨跌额"`
	ChangeRate float64   `comment:"涨跌幅"`
	Price      float64   `comment:"最新价"`
	Desc       string    `comment:"描述内容"`
	Pbannual   float64   `comment:"市净率"`
	Peannual   float64   `comment:"市盈率"`
}

func (*StockPredictionDayList) TableName() string {
	return "stock_prediction_day_list"
}
