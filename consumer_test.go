package nsqclient

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/goapt/logger"
	"github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestNsqConsumer_RunHandler(t *testing.T) {
	consumer := &NsqHandler{
		Topic:            "test",
		Channel:          "ai",
		OpenChannelTopic: true,
	}

	consumer.SetHandle(func(log logger.ILogger, message *nsq.Message) error {
		assert.Equal(t, []byte(`hello`), message.Body)
		fmt.Println("====>", string(message.Body))
		return nil
	})

	Register("test", consumer)
	ctx, cancel := context.WithCancel(context.Background())

	p, err := NewProducer("default")
	require.NoError(t, err)
	err = p.Publish("test", []byte(`hello`))
	require.NoError(t, err)
	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()
	Run("test", ctx)
}
