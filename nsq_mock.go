package nsqclient

import (
	"context"
	"time"

	"github.com/nsqio/go-nsq"
)

type NoMessageDelegate struct{}

func (d *NoMessageDelegate) OnFinish(m *nsq.Message)                                     {}
func (d *NoMessageDelegate) OnRequeue(m *nsq.Message, delay time.Duration, backoff bool) {}
func (d *NoMessageDelegate) OnTouch(m *nsq.Message)                                      {}

func MockMessage(body []byte) *nsq.Message {
	msgID := nsq.MessageID{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'a', 's', 'd', 'f', 'g', 'h'}
	return &nsq.Message{
		ID:          msgID,
		Body:        body,
		NSQDAddress: "",
		Delegate:    &NoMessageDelegate{},
	}
}

func RunMock(ctx context.Context, h *NsqHandler, msg *nsq.Message) error {
	h.runInit(ctx)
	fn := runNsqConsumer(h, false)
	return fn(msg)
}
