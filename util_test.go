package nsqclient

import (
	"errors"
	"testing"
)

func Test_retry(t *testing.T) {
	type args struct {
		num int
		fn  func() error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{3, func() error {
				return nil
			}},
			wantErr: false,
		},
		{
			name: "test",
			args: args{2, func() error {
				return errors.New("xxxx")
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := retry(tt.args.num, tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("retry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
