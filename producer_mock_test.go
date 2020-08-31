package nsqclient

import (
	"reflect"
	"testing"
	"time"
)

func TestNewMockProducer(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *producerMock
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				"test",
			},
			want: &producerMock{
				name: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMockProducer(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMockProducer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMockProducer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_producerMock_Publish(t *testing.T) {
	type fields struct {
		name string
	}
	type args struct {
		topic string
		body  []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				name: "test",
			},
			args:    args{"test", []byte("test")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &producerMock{
				name: tt.fields.name,
			}
			if err := p.Publish(tt.args.topic, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_producerMock_MultiPublish(t *testing.T) {
	type fields struct {
		name string
	}
	type args struct {
		topic string
		body  [][]byte
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
			name: "test",
			fields: fields{
				name: "test",
			},
			args:    args{"test", body},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &producerMock{
				name: tt.fields.name,
			}
			if err := p.MultiPublish(tt.args.topic, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("MultiPublish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_producerMock_DeferredPublish(t *testing.T) {
	type fields struct {
		name string
	}
	type args struct {
		topic string
		delay time.Duration
		body  []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "test",
			fields: fields{name: "test"},
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
			p := &producerMock{
				name: tt.fields.name,
			}
			if err := p.DeferredPublish(tt.args.topic, tt.args.delay, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("DeferredPublish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
