package nsqclient

import (
	"errors"
	"time"

	"github.com/goapt/nsqclient/internal/pool"
)

type Producer interface {
	Publish(topic string, body []byte) error
	MultiPublish(topic string, body [][]byte) error
	DeferredPublish(topic string, delay time.Duration, body []byte) error
}

var _ Producer = (*producer)(nil)

type producer struct {
	pool pool.Pool
}

func NewProducer(name string) (*producer, error) {
	p, ok := Client(name)
	if !ok {
		return nil, errors.New("nsq producer config not found")
	}
	return &producer{
		pool: p,
	}, nil
}

func (p *producer) Publish(topic string, body []byte) error {
	nsq, err := p.pool.Get()
	if err != nil {
		return err
	}
	defer nsq.Close()

	return retry(2, func() error {
		return nsq.Publish(topic, body)
	})
}

func (p *producer) MultiPublish(topic string, body [][]byte) error {
	nsq, err := p.pool.Get()
	if err != nil {
		return err
	}
	defer nsq.Close()

	return retry(2, func() error {
		return nsq.MultiPublish(topic, body)
	})
}

func (p *producer) DeferredPublish(topic string, delay time.Duration, body []byte) error {
	nsq, err := p.pool.Get()
	if err != nil {
		return err
	}
	defer nsq.Close()

	return retry(2, func() error {
		return nsq.DeferredPublish(topic, delay, body)
	})
}
