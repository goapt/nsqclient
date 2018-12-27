package nsqclient

import (
	"sync"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/verystar/golib/debug"
	"github.com/verystar/logger"
	"github.com/verystar/nsqclient/delay"
)

var _ INsqHandler = (*NsqHandler)(nil)
var mu sync.Mutex

type HandleFunc func(debug *debug.DebugTag, log logger.ILogger, message *nsq.Message) error

var NsqGroups = make(map[string][]INsqHandler)

type NsqHandler struct {
	Topic            string
	Channel          string
	Size             int
	MaxAttepts       uint16
	OpenChannelTopic bool // 是否开启独立的topic [Topic.Channel]
	Handler          HandleFunc
	initFn           func()
	shouldRequeue    func(message *nsq.Message) (bool, time.Duration)
}

func NewNsqHandler(options ... func(*NsqHandler)) *NsqHandler {
	handler := new(NsqHandler)
	for _, option := range options {
		option(handler)
	}
	return handler
}

func (h *NsqHandler) Init(fn func()) {
	h.initFn = fn
}

func (h *NsqHandler) RunInit() {
	h.initFn()
}

func (h *NsqHandler) GetTopic() string {
	return h.Topic
}

func (h *NsqHandler) IsOpenChannelTopic() bool {
	return h.OpenChannelTopic
}

func (h *NsqHandler) GetChannelTopic() string {
	return h.Topic + "." + h.Channel
}

func (h *NsqHandler) GetChannel() string {
	return h.Channel
}

func (h *NsqHandler) SetHandle(fn HandleFunc) {
	h.Handler = fn
}

func (h *NsqHandler) GetHandle() HandleFunc {
	return h.Handler
}

func (h *NsqHandler) GetMaxAttepts() uint16 {
	mu.Lock()
	defer mu.Unlock()
	if h.MaxAttepts == 0 {
		h.MaxAttepts = 100
	}
	return h.MaxAttepts
}

func (h *NsqHandler) SetShouldRequeue(fn func(message *nsq.Message) (bool, time.Duration)) {
	h.shouldRequeue = fn
}

func (h *NsqHandler) GetShouldRequeue(message *nsq.Message) (bool, time.Duration) {
	if h.shouldRequeue == nil {
		return delay.DefaultDelay(message)
	}

	return h.shouldRequeue(message)
}

func (h *NsqHandler) Group(group string) {
	mu.Lock()
	defer mu.Unlock()
	NsqGroups[group] = append(NsqGroups[group], h)
}

func (h *NsqHandler) GetSize() int {
	if h.Size == 0 {
		return 1
	}
	return h.Size
}
