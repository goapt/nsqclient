package nsqclient

import (
	"log"
	"os"
	"strings"

	"github.com/nsqio/go-nsq"

	"github.com/goapt/nsqclient/internal/pool"
)

type Config struct {
	Host     string `toml:"host" json:"host"`
	Port     string `toml:"port" json:"port"`
	InitSize int    `toml:"init_size" json:"init_size"`
	MaxSize  int    `toml:"max_size" json:"max_size"`
}

var (
	nsqConfigs map[string]Config
	nsqList    map[string]pool.Pool
	errs       []string
)

func Connect(configs map[string]Config) {
	defer func() {
		if len(errs) > 0 {
			panic("[nsq] " + strings.Join(errs, "\n"))
		}
	}()

	nsqConfigs = configs
	nsqList = make(map[string]pool.Pool)

	for name, conf := range configs {
		n, err := newProducerPool(conf.Host+":"+conf.Port, conf.InitSize, conf.MaxSize)
		log.Printf("[nsq] connect:%s:%s:%s", name, conf.Host, conf.Port)
		if err == nil {
			nsqList[name] = n
			// 支持ip:port寻址
			nsqList[conf.Host+":"+conf.Port] = n
		} else {
			errs = append(errs, err.Error())
		}
	}
}

func Client(name ...string) (pool.Pool, bool) {
	key := "default"
	if len(name) > 0 {
		key = name[0]
	}
	n, ok := nsqList[key]
	return n, ok
}

// CreateNSQProducerPool create a nwq producer pool
func newProducerPool(addr string, initSize, maxSize int, options ...func(*nsq.Config)) (pool.Pool, error) {
	factory := func() (*nsq.Producer, error) {
		// TODO 这里应该执行ping方法来确定连接是正常的否则不应该创建conn
		return newProducer(addr, options...)
	}
	nsqPool, err := pool.NewChannelPool(initSize, maxSize, factory)
	if err != nil {
		return nil, err
	}
	return nsqPool, nil
}

// CreateNSQProducer create nsq producer
func newProducer(addr string, options ...func(*nsq.Config)) (*nsq.Producer, error) {
	cfg := nsq.NewConfig()
	for _, option := range options {
		option(cfg)
	}

	producer, err := nsq.NewProducer(addr, cfg)
	if err != nil {
		return nil, err
	}
	producer.SetLogger(log.New(os.Stderr, "", log.Flags()), nsq.LogLevelError)
	return producer, nil
}
