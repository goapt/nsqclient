package nsqclient

import (
	"context"
	"fmt"
	"time"

	"github.com/goapt/logger"
	"github.com/nsqio/go-nsq"

	"github.com/goapt/nsqclient/internal/color"
)

type NsqConsumer struct {
	consumer *nsq.Consumer
	handler  nsq.Handler
	ctx      context.Context
	topic    string
	channel  string
}

func NewNsqConsumer(ctx context.Context, topic, channel string, options ...func(*nsq.Config)) (*NsqConsumer, error) {
	conf := nsq.NewConfig()
	conf.MaxAttempts = 0
	conf.MsgTimeout = 10 * time.Minute         // 默认一个消息最多能处理十分钟，否则就会重新丢入队列
	conf.LookupdPollInterval = 3 * time.Second // 调整consumer的重连间隔时间为3秒
	for _, option := range options {
		option(conf)
	}

	consumer, err := nsq.NewConsumer(topic, channel, conf)
	if err != nil {
		return nil, err
	}
	return &NsqConsumer{
		consumer: consumer,
		ctx:      ctx,
		topic:    topic,
		channel:  channel,
	}, nil
}

func (n *NsqConsumer) AddHandler(handler nsq.Handler) {
	n.handler = handler
}

func (n *NsqConsumer) Run(conf *Config, concurrency int) {
	n.consumer.ChangeMaxInFlight(concurrency)
	n.consumer.AddConcurrentHandlers(n.handler, concurrency)

	if err := n.consumer.ConnectToNSQD(conf.Host + ":" + conf.Port); err != nil {
		logger.Error("nsq:ConnectToNSQD", err)
		return
	}
	for {
		select {
		case <-n.ctx.Done():
			fmt.Println(color.Yellow("[%s] %s,%s", "stop consumer", n.topic, n.channel))
			n.consumer.Stop()
			fmt.Println(color.Yellow("[%s] %s,%s", "stop consumer success", n.topic, n.channel))
			return
		}
	}
}
