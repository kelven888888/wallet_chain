package test

import (
	"context"
	"fmt"
	"testing"
	"time"
	"wallet_chain.com/utils"
	"wallet_chain.com/utils/rabbitmq"
)

var testChannelName = "pub-sub"

func runPublisher(ctx context.Context, channelname string) error {

	var publisherErr error
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {
		defer cancel()
		connection, err := rabbitmq.NewConnection(rabbitmq.DefaultURL)
		if err != nil {
			publisherErr = err
			return
		}
		defer connection.Close()

		p, err := rabbitmq.NewPublisher(channelname, connection)
		if err != nil {
			publisherErr = err
			return
		}
		for i := 1; i < 100; i++ {
			data := []byte("helloword" + time.Now().Format("2006-01-02 15:04:05.000"))
			err = p.Publish(ctx, data)
			if err != nil {
				publisherErr = err
				return
			}
		}
	})
	return publisherErr
}
func TestPublisher(t *testing.T) {
	err := runPublisher(context.Background(), testChannelName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
