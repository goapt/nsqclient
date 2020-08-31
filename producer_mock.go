package nsqclient

import (
	"encoding/json"
	"log"
	"time"
)

var _ Producer = (*producerMock)(nil)

type producerMock struct {
	name string
}

func NewMockProducer(name string) (*producerMock, error) {
	return &producerMock{
		name: name,
	}, nil
}

func (p *producerMock) Publish(topic string, body []byte) error {
	log.Printf("nsq:%s,topic:%s,body:%s", p.name, topic, string(body))
	return nil
}

func (p *producerMock) MultiPublish(topic string, body [][]byte) error {
	b, _ := json.Marshal(body)
	log.Printf("nsq:%s,topic:%s,body:%s", p.name, topic, string(b))
	return nil
}

func (p *producerMock) DeferredPublish(topic string, delay time.Duration, body []byte) error {
	log.Printf("nsq:%s,topic:%s,delay:%d,body:%s", p.name, topic, delay, string(body))
	return nil
}
