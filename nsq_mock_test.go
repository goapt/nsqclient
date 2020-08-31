package nsqclient

import (
	"context"
	"errors"
	"testing"

	"github.com/goapt/logger"
	"github.com/nsqio/go-nsq"
)

var exampleConsumer = func() *NsqHandler {
	consumer := &NsqHandler{
		Topic:   "test",
		Channel: "test",
	}

	// job handler
	consumer.SetHandle(func(log logger.ILogger, message *nsq.Message) error {
		if string(message.Body) != "hello" {
			return errors.New("test message is not hello")
		}
		return nil
	})
	return consumer
}

func TestRunMock(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type args struct {
		ctx context.Context
		h   *NsqHandler
		msg *nsq.Message
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				ctx: ctx,
				h:   exampleConsumer(),
				msg: MockMessage([]byte("hello")),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RunMock(tt.args.ctx, tt.args.h, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("RunMock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
