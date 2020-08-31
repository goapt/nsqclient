package nsqclient

import (
	"reflect"
	"testing"
	"time"

	"github.com/goapt/nsqclient/internal/pool"
)

func TestNewProducer(t *testing.T) {
	type args struct {
		name string
	}

	p, ok := Client("default")

	if !ok {
		t.Fatal("default client not found")
	}

	tests := []struct {
		name    string
		args    args
		want    *producer
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				name: "default",
			},
			want: &producer{
				pool: p,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProducer(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProducer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProducer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_producer_Publish(t *testing.T) {
	type fields struct {
		pool pool.Pool
	}
	type args struct {
		topic string
		body  []byte
	}

	p, ok := Client("default")

	if !ok {
		t.Fatal("default client not found")
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "publish",
			fields: fields{
				pool: p,
			},
			args: args{
				topic: "test",
				body:  []byte("test"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &producer{
				pool: tt.fields.pool,
			}
			if err := p.Publish(tt.args.topic, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_producer_MultiPublish(t *testing.T) {
	type fields struct {
		pool pool.Pool
	}
	type args struct {
		topic string
		body  [][]byte
	}

	p, ok := Client("default")

	if !ok {
		t.Fatal("default client not found")
	}
	body := make([][]byte, 0)
	body = append(body, []byte("test"))
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "publish",
			fields: fields{
				pool: p,
			},
			args: args{
				topic: "test",
				body:  body,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &producer{
				pool: tt.fields.pool,
			}
			if err := p.MultiPublish(tt.args.topic, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("MultiPublish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_producer_DeferredPublish(t *testing.T) {
	type fields struct {
		pool pool.Pool
	}
	type args struct {
		topic string
		delay time.Duration
		body  []byte
	}
	p, ok := Client("default")

	if !ok {
		t.Fatal("default client not found")
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "test",
			fields: fields{pool: p},
			args: args{
				topic: "test",
				delay: 10 * time.Second,
				body:  []byte("test"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &producer{
				pool: tt.fields.pool,
			}
			if err := p.DeferredPublish(tt.args.topic, tt.args.delay, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("DeferredPublish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
