package nsqclient

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/goapt/logger"
	gonsq "github.com/nsqio/go-nsq"

	"github.com/goapt/nsqclient/internal/color"
)

func Run(group string, ctx context.Context) {
	stop := make(chan struct{})
	defer close(stop)
	if !runMulti(ctx, group) {
		panic("未匹配到NSQ的运行时参数")
	}

	for {
		select {
		case <-ctx.Done():
			// 给3秒的处理时间
			fmt.Println(color.Yellow("%s", "----------------------------------"))
			fmt.Println(color.Yellow("%s , %s", "| give nsq consumer 3 second", "... |"))
			fmt.Println(color.Yellow("%s", "----------------------------------"))
			time.AfterFunc(time.Second*3, func() {
				stop <- struct{}{}
			})
			goto END
		}
	}
END:
	<-stop
}

func Register(group string, h ...*NsqHandler) {
	for _, v := range h {
		v.group(group)
	}
}

func runMulti(ctx context.Context, group string) bool {

	if _, check := nsqGroups[group]; !check {
		return false
	}

	for _, h := range nsqGroups[group] {
		h.runInit(ctx)

		fn := runNsqConsumer(h, false)
		go runHandler(ctx, h, false, fn)
		if h.isOpenChannelTopic() {
			fn := runNsqConsumer(h, true)
			go runHandler(ctx, h, true, fn)
		}
	}
	return true
}

func runNsqConsumer(h *NsqHandler, isChannelTopic bool) gonsq.HandlerFunc {
	if h.Logger == nil {
		h.Logger = logger.NewLogger(func(config *logger.Config) {
			config.LogName = h.Topic
		})
	}

	var fn gonsq.HandlerFunc = func(m *gonsq.Message) error {
		m.DisableAutoResponse()
		done := make(chan struct{})

		defer func() {
			if h.TouchDuration > 0 {
				done <- struct{}{}
			}
			if err := recover(); err != nil {
				h.Logger.Error("[Nsq Consumer Handler Recover]%s%s", err, string(stack(1)))
				should, t := h.getShouldRequeue(m)
				if should {
					m.Requeue(t)
				}
			}
		}()

		// 如果开启了定时touch，则定时touch message，让服务端保持消息活跃
		if h.TouchDuration > 0 {
			go func() {
				t := time.NewTicker(h.TouchDuration)

				for {
					select {
					case <-t.C:
						m.Touch()
					case <-done:
						t.Stop()
						return
					}
				}
			}()
		}

		err := h.handler(h.Logger, m)
		if err != nil {
			errdata := map[string]string{
				"error":   err.Error(),
				"channel": h.getChannelTopic(),
				"data":    string(m.Body),
			}
			if _, ok := err.(*debugError); ok {
				h.Logger.Data(errdata).Debugf("[NSQ Consumer Error] "+h.getChannelTopic()+" %s", err.Error())
			} else {
				h.Logger.Data(errdata).Errorf("[NSQ Consumer Error] "+h.getChannelTopic()+" %s", err.Error())
			}

			should, t := h.getShouldRequeue(m)
			if should {
				// 如果任务开启了channel topic功能，并且当前handler的topic不是channel topic，则写入到channel topic
				if h.isOpenChannelTopic() && !isChannelTopic {
					nsq, err := NewProducer(h.Connect)
					if err != nil {
						logger.Error("[NSQ Push Channel Topic] new producer error", err)
						m.RequeueWithoutBackoff(t)
					} else {
						err := nsq.Publish(h.getChannelTopic(), m.Body)
						if err != nil {
							logger.Error("[NSQ Push Channel Topic] publish error", err)
							m.RequeueWithoutBackoff(t)
						} else {
							// 丢入子topic之后，主topic将消息结束掉
							m.Finish()
						}
					}

				} else {
					m.RequeueWithoutBackoff(t)
					// m.Requeue(t)
				}
				return nil
			}
		}
		m.Finish()
		return nil
	}
	return fn
}

func runHandler(ctx context.Context, h *NsqHandler, isChannelTopic bool, fn gonsq.HandlerFunc) {
	var topic string
	if isChannelTopic {
		topic = h.getChannelTopic()
	} else {
		topic = h.Topic
	}

	manager, err := NewNsqConsumer(ctx, topic, h.Channel, func(nc *gonsq.Config) {
		nc.MaxAttempts = h.getMaxAttempts()
	})

	if err != nil {
		log.Println("NSQ Consumer err:", err)
	}

	manager.AddHandler(fn)
	conf, err := h.conf()
	if err != nil {
		logger.Fatal("NSQ config error", err)
	}

	log.Println(fmt.Sprintf("NSQ Consumer Run:topic[%s] channel[%s] concurrent[%d]", topic, h.Channel, h.getSize()))
	go manager.Run(conf, h.getSize())
}

// Stack gets the call stack
func stack(calldepth int) []byte {
	var (
		e             = make([]byte, 1<<16) // 64k
		nbytes        = runtime.Stack(e, false)
		ignorelinenum = 2*calldepth + 1
		count         = 0
		startIndex    = 0
	)
	for i := range e {
		if e[i] == '\n' {
			count++
			if count == ignorelinenum {
				startIndex = i + 1
			}
		}
	}
	return e[startIndex:nbytes]
}
