---
title:LXC容器搭建NSQ集群二
categories:
  - GOLANG
---
### LXC 容器搭建nsq集群
#### go-nsq
- 生产者角色和消费者角色

```yaml
version: '2' # 高版本支持3
services:
  nsqlookupd:
    image: nsqio/nsq
    command: /nsqlookupd
    ports:
      - "4160:4160" # tcp
      - "4161:4161" # http

  nsqd:
    image: nsqio/nsq
    # 广播地址不填的话默认就是hostname(或虚拟机名称)，不加就变成lxc的容器名，那样 lookupd 会连接不上，所以直接写IP
    command: /nsqd --broadcast-address=10.220.151.50 --lookupd-tcp-address=nsqlookupd:4160
    depends_on:
      - nsqlookupd
    ports:
      - "4150:4150" # tcp
      - "4151:4151" # http

  nsqadmin:
    image: nsqio/nsq
    command: /nsqadmin --lookupd-http-address=nsqlookupd:4161
    depends_on:
      - nsqlookupd
    ports:
      - "4171:4171" # http
```
##### producer.go
```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nsqio/go-nsq"
)

func main() {
	config := nsq.NewConfig()
	//直连的是nsqd
	p, err := nsq.NewProducer("10.8.0.6:4150", config)

	if err != nil {
		log.Panic(err)
	}

	for i := 1000; i < 2000; i++ {
		msg := fmt.Sprintf("num-%d", i)
		log.Println("Pub2:" + msg)
		err = p.Publish("testTopic", []byte(msg))
		if err != nil {
			log.Panic(err)
		}
		time.Sleep(time.Second * 1)
	}

	p.Stop()
}

```
##### consumer.go
```go
package main

import (
	"log"
	"sync"

	"github.com/nsqio/go-nsq"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1000)

	config := nsq.NewConfig()
	c, _ := nsq.NewConsumer("testTopic", "ch", config)
	c.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Printf("Got a message: %s", message.Body)
		wg.Done()
		return nil
	}))

	// 1.直连nsqd
	// err := c.ConnectToNSQD("127.0.0.1:4150")

	// 2.通过 nsqlookupd 服务发现
	err := c.ConnectToNSQLookupd("10.8.0.4:4161")
	if err != nil {
		log.Panic(err)
	}
	wg.Wait()
}

```
