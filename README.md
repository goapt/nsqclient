# nsqclient
nsq client for golang

<a href="https://github.com/goapt/nsqclient/actions"><img src="https://github.com/goapt/nsqclient/workflows/build/badge.svg" alt="Build Status"></a>
<a href="https://codecov.io/gh/goapt/nsqclient"><img src="https://codecov.io/gh/goapt/nsqclient/branch/master/graph/badge.svg" alt="codecov"></a>
<a href="https://goreportcard.com/report/github.com/goapt/nsqclient"><img src="https://goreportcard.com/badge/github.com/goapt/nsqclient" alt="Go Report Card
"></a>
<a href="https://pkg.go.dev/github.com/goapt/nsqclient"><img src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square" alt="GoDoc"></a>
<a href="https://opensource.org/licenses/mit-license.php" rel="nofollow"><img src="https://badges.frapsoft.com/os/mit/mit.svg?v=103"></a>


##  Connection
```go
nsqclient.Connect(map[string]Config{
    "default": Config{
        Host:     "10.64.146.231",
        Port:     "4150",
        InitSize: 10,
        MaxSize:  10,
    },
    "other": Config{
        Host:     "10.64.146.22",
        Port:     "4150",
        InitSize: 10,
        MaxSize:  10,
    },
})
```

## Publish
写入nsq为了方便进行单元测试，可以注入nsqclient.Producer接口，然后在单元测试使用mock的实现
```go
// default client
nsq := nsqclient.NewProducer("default")
nsq.Publish("topic","body")

//mock
nsq := nsqclient.NewMockProducer("default")
nsq.Publish("topic","body")
```

## Consumer
```go
// default connect
consumer := &nsqclient.NsqHandler{
    Topic:   "log",
    Channel: "default",
    Size:   10,
}

consumer.SetHandle(func(lg logger.ILogger, message *nsq.Message) error {
    // something
    return nil
})

nsqclient.Register(consumer, "access_log")

nsqclient.Run("access_log", ctx)

// other connect
consumer := &nsqclient.NsqHandler{
    Connect: "other",
    Topic:   "log",
    Channel: "default",
    Size:   10,
}
```