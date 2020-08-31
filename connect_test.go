package nsqclient

import (
	"os"
	"testing"
)

func TestConnect(t *testing.T) {
	nsqHost := "10.64.146.231"
	if os.Getenv("DRONE") == "true" {
		nsqHost = "nsqd"
	}

	Connect(map[string]Config{
		"default": {
			Host:     nsqHost,
			Port:     "4150",
			InitSize: 1,
			MaxSize:  2,
		},
	})

	_, ok := Client("default")

	if !ok {
		t.Errorf("nsq connect error")
	}
}
