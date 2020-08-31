package nsqclient

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/goapt/logger"
	"github.com/nsqio/go-nsq"

	"github.com/goapt/nsqclient/delay"
)

var mu sync.Mutex

type HandleFunc func(log logger.ILogger, message *nsq.Message) error

var nsqGroups = make(map[string][]*NsqHandler)

type NsqHandler struct {
	Connect          string        // 连接的nsq 默认是default
	Topic            string        // nsq topic
	Channel          string        // topic channel
	Size             int           // 并发数MaxInFlight
	MaxAttempts      uint16        // 最大执行次数，默认是100
	OpenChannelTopic bool          // 是否开启独立的topic [Topic.Channel]
	TouchDuration    time.Duration // 多久之后touch一次当前message，保持消息存活，默认不Touch
	Logger           logger.ILogger
	handler          HandleFunc
	initFn           func(ctx context.Context)
	shouldRequeue    func(message *nsq.Message) (bool, time.Duration)
}

func NewNsqHandler(options ...func(*NsqHandler)) *NsqHandler {
	handler := new(NsqHandler)
	for _, option := range options {
		option(handler)
	}
	return handler
}

func (h *NsqHandler) connectName() string {
	if h.Connect == "" {
		h.Connect = "default"
	}

	return h.Connect
}

func (h *NsqHandler) conf() (*Config, error) {
	c, ok := nsqConfigs[h.connectName()]
	if !ok {
		return nil, fmt.Errorf("nsq config not found:%s", h.connectName())
	}
	return &c, nil
}

func (h *NsqHandler) Init(fn func(ctx context.Context)) {
	h.initFn = fn
}

func (h *NsqHandler) runInit(ctx context.Context) {
	if h.initFn != nil {
		h.initFn(ctx)
	}
}

func (h *NsqHandler) isOpenChannelTopic() bool {
	return h.OpenChannelTopic
}

func (h *NsqHandler) getChannelTopic() string {
	return h.Topic + "." + h.Channel
}

func (h *NsqHandler) SetHandle(fn HandleFunc) {
	h.handler = fn
}

func (h *NsqHandler) getMaxAttempts() uint16 {
	mu.Lock()
	defer mu.Unlock()
	if h.MaxAttempts == 0 {
		h.MaxAttempts = 100
	}
	return h.MaxAttempts
}

func (h *NsqHandler) SetShouldRequeue(fn func(message *nsq.Message) (bool, time.Duration)) {
	h.shouldRequeue = fn
}

func (h *NsqHandler) getShouldRequeue(message *nsq.Message) (bool, time.Duration) {
	if h.shouldRequeue == nil {
		return delay.DefaultDelay(message)
	}

	return h.shouldRequeue(message)
}

func (h *NsqHandler) group(group string) {
	mu.Lock()
	defer mu.Unlock()
	nsqGroups[group] = append(nsqGroups[group], h)
}

func (h *NsqHandler) getSize() int {
	if h.Size == 0 {
		return 2
	}
	return h.Size
}
