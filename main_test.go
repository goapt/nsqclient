package nsqclient

import "os"

func init() {
	nsqHost := "127.0.0.1"
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
}
