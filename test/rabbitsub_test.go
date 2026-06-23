package test

import (
	"context"
	"fmt"
	"testing"
	"time"
	"wallet_chain.com/utils"
	"wallet_chain.com/utils/rabbitmq"
)

//var testChannelName = "pub-sub"

func TestSubscriber(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)

	err := runPublisher(ctx, testChannelName)
	if err != nil {
		t.Log(err)
		return
	}

	err = runSub(ctx, testChannelName, "fanout-queue-1")
	if err != nil {
		t.Log(err)
		return
	}

	err = runSub(ctx, testChannelName, "fanout-queue-2")
	if err != nil {
		t.Log(err)
		return
	}

	<-ctx.Done()
	time.Sleep(time.Millisecond * 100)
}

func runSub(ctx context.Context, channelname string, identifier string) error {
	var subscriberEr error
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {
		defer cancel()
		connection, err := rabbitmq.NewConnection(rabbitmq.DefaultURL)
		if err != nil {
			subscriberEr = err
			return
		}
		s, err := rabbitmq.NewSubscriber(channelname, identifier, connection, rabbitmq.WithConsumerAutoAck(false))
		if err != nil {
			subscriberEr = err
			return
		}
		s.Subscribe(ctx, handler)
	})
	return subscriberEr

}

var handler = func(ctx context.Context, data []byte, tagID string) error {
	fmt.Printf("[received]: tagID=%s, data=%s\n", tagID, data)
	//err := errors.New("something went wrong")
	return nil
}
