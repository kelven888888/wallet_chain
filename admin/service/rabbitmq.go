package service

import (
	"context"
	"fmt"
	"time"
	"wallet_chain.com/global"
	"wallet_chain.com/utils"
	"wallet_chain.com/utils/rabbitmq"
)

var testChannelName = "pub-sub"
var Connect *rabbitmq.Connection

func inits() {

	connection, err := rabbitmq.NewConnection(rabbitmq.DefaultURL)
	if err != nil {
		return
	}
	Connect = connection
	//defer connection.Close()

}
func RunPublisher(ctx context.Context, channelname string, msg string) {

	var publisherErr error
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {
		defer cancel()

		p, err := rabbitmq.NewPublisher(channelname, Connect)
		if err != nil {
			publisherErr = err
			return
		}

		data := []byte(msg)
		err = p.Publish(ctx, data)
		if err != nil {
			publisherErr = err
			return

		}

	})
	if publisherErr != nil {
		fmt.Printf(publisherErr.Error())
	}
}
func RunSubscriber(ctx context.Context, channelName string, identifier string) {
	var subscriberErr error
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {
		defer cancel()

		s, err := rabbitmq.NewSubscriber(channelName, identifier, Connect, rabbitmq.WithConsumerAutoAck(false))
		if err != nil {
			subscriberErr = err
			return
		}

		s.Subscribe(ctx, handler)

	})
	if subscriberErr != nil {
		fmt.Printf(subscriberErr.Error())
	}

}

type Msg struct {
	Id       int
	Msg      string
	CreateAt string `gorm:"column:created_at"`
}

func (Msg) TableName() string {
	return "nov_msg"
}

var handler = func(ctx context.Context, data []byte, tagID string) error {
	fmt.Printf("[received]: tagID=%s, data=%s\n", tagID, data)
	var msg Msg
	global.SHOP_DB.AutoMigrate(&msg)
	msg.Msg = string(data)
	msg.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	err := global.SHOP_DB.Save(&msg).Error
	if err != nil {
		global.SHOP_LOG.Error(err.Error())
	}
	return nil
}
