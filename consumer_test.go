package nsqclient

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/goapt/logger"
	"github.com/nsqio/go-nsq"
)

func TestNsqConsumer_AddHandler(t *testing.T) {
	consumer := &NsqHandler{
		Topic:            "test",
		Channel:          "ai",
		OpenChannelTopic: true,
	}

	consumer.SetHandle(func(log logger.ILogger, message *nsq.Message) error {
		fmt.Println(string(message.Body))
		return errors.New("error")
	})

	ctx, _ := context.WithCancel(context.Background())
	err := RunMock(ctx, consumer, MockMessage([]byte("hello")))

	if err != nil {
		t.Error(err)
	}
}
