---
title:LXC容器搭建NSQ集群三
categories:
  - GOLANG
---
### LXC 容器搭建nsq集群
- [原文链接](http://ddrv.cn/a/291527)  

#### 使用场景
有两个服务中心，apiserver处理用户的请求，dataserver处理数据，dataserver 有很多个，并且遵循raft高可用协议，写在配置文件中，当用户发起一个请求到apiserver,apiserver会向dataserver中的任意一个请求，这样做的好处是当dataserver中任意一个挂掉，apiserver只需获取当前配置文件中可用的地址即可，存活的dataserver会返回地址。  
#### NSQ当中的用处：
- 每一个dataserver 以生产者的角色，连接nsqd，同时每隔几秒发送topic消息心跳检测，（可以有多个channel，nsq 保证消息投递一次以上，随机挑选channel发送消息，无序性）保证nsqlookupd 知道有哪些nsqd对应的生产者存活。  
- apiserver充当消费者的角色，连接nsqlookupd 4161 http端口，监听topic消息以及channel1，当有消息publish，即通知apiserver,本文即所有存活的dataserver的地址，apiserver 可随机挑选任意一个dataserver进行调用。   
#### 代码

##### dataserver.go

```go
package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"nsq_client/config"
	"time"

	"github.com/nsqio/go-nsq"
)

// Conf *
var Conf *config.Cfgparams

func main() {

	Conf = config.GetConfig()
	// randomAddr := RadomGetNsqd(Conf.NSQ_TCP_ADDRS)
	fmt.Println(Conf)
	// name := viper.GetString("name")
	// x := viper.Get("common.nsqd")
	// fmt.Println(x, reflect.TypeOf(x))

	server := http.Server{
		Addr: Conf.SERVER_ADDR,
	}
	go HeartBeat() //这里开goroutine来发送心跳包
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

	// fmt.Println("Viper get name:", name)
}

// RadomGetNsqd
func RadomGetNsqd(nsqTCPAddrs []string) string {
	rand.Seed(time.Now().UnixNano())
	return nsqTCPAddrs[rand.Intn(len(nsqTCPAddrs))]
}

// HeartBeat
func HeartBeat() {
	var (
		producer *nsq.Producer
		err      error
		ticker   *time.Ticker
	)
	ticker = time.NewTicker(time.Duration(Conf.HEART_BEAT_TIME) * time.Second) //创建一个定时器，这里的时间我都写到配置文件里了，然后用conf包拿出来。我这里是设置的5秒
	for {
		select {
		case <-ticker.C:
			//创建一个生产者，这里的RanDomGetServer()是自定义的一个工具，用来随机获取一个nsqd地址
			if producer, err = nsq.NewProducer(RadomGetNsqd(Conf.NSQ_TCP_ADDRS), nsq.NewConfig()); err != nil {
				panic(err)
			}
			//pulish()接的第一个参数是topic，这个也是自己定义(值为data_server_addr），第二个参数是要发送的消息，这里是本机服务器地址
			err = producer.Publish(Conf.DATA_SERVER_TOPIC, []byte(Conf.SERVER_ADDR))
			if err != nil {
				panic(err)
			}
		}
	}
}
```
##### apiserver.go
```go
package main

import (
	"fmt"
	"net/http"
	"nsq_client/config"
	"sync"
	"time"

	"github.com/nsqio/go-nsq"
)

//创建一个nsq.handler接口实例
var ServerConsumerHandler = &ServerConsumer{DataServerAddrs: make(map[string]time.Time)}

//ServerConsumer实现了nsq.handler接口
type ServerConsumer struct {
	DataServerAddrs map[string]time.Time //保存dataServer进程发过来的服务器地址和接收时间
	rwLocker        sync.RWMutex         //防止多线程同时读写
}

//这里第一个参数是需要绑定topic(data_server_addr),
//第二个参数传入一个string，这就是创建的ch，topic消息队列中的消息会分发到每个ch中.
//每个消费者可以创建不通的ch，也可以共用一个ch，共用一个ch, ch的消息会随机发送给其中一个消费者
//第三个参数是处理message的nsq.handler接口，需要实现一个HanddleMessage(*nsq.Message)error()方法。
func NewConsumer(topic string, chanName string, h nsq.Handler) (consumer *nsq.Consumer, err error) {
	if consumer, err = nsq.NewConsumer(topic, chanName, nsq.NewConfig()); err != nil {
		return nil, err
	}
	consumer.ChangeMaxInFlight(3) //可以根据nsqds数量来配置
	consumer.AddHandler(h)
	err = consumer.ConnectToNSQLookupd(Conf.NSQ_LOOKUPD_ADDR) //todo:读取配置地址
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

//HandleMessage是nsq,handler接口定义方法，必须实现，主要是用来处理传递过来的消息
func (h *ServerConsumer) HandleMessage(message *nsq.Message) error {
	if dataServer := string(message.Body); dataServer != "" {
		h.rwLocker.Lock()
		h.DataServerAddrs[dataServer] = time.Now()
		h.rwLocker.Unlock()
	}
	return nil
}

//监听服务data_server_addr这个消息队列
func MonitorDataServerAddrs() {
	consumer, err := NewConsumer("data_server_addr", "ch1", ServerConsumerHandler)
	if err != nil {
		panic(err)
	}
	//连接到NSQLookupd，它底层会创建连接到每个nsqd.这样就可以监听每个nsqd中的消息
	err = consumer.ConnectToNSQLookupd(Conf.NSQ_LOOKUPD_ADDR) //方法本身会开一个goroutine来检查消息队列
	if err != nil {
		panic(err)
	}
}

//删除过期的服务器地址
func (h *ServerConsumer) removeExpireDatasServers() {
	for {
		h.rwLocker.Lock()
		for dataServer, t := range h.DataServerAddrs {
			//只保存10秒之内发送过来的服务器地址
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(h.DataServerAddrs, dataServer)
			}
		}
		h.rwLocker.Unlock()
		time.Sleep(2 * time.Second)
	}
}

var Conf *config.Cfgparams

func main() {

	Conf = config.GetConfig()
	fmt.Println(Conf)

	server := http.Server{
		Addr: "127.0.0.1:9000",
	}
	MonitorDataServerAddrs() //不用单独开goroutine
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
```
##### config.go
```go
package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Cfgparams struct {
	HEART_BEAT_TIME   int
	SERVER_ADDR       string
	NSQ_LOOKUPD_ADDR  string
	NSQ_TCP_ADDRS     []string
	DATA_SERVER_TOPIC string
}

type Config struct {
	Name string
}

func Init(cfg string) error {
	c := Config{
		Name: cfg,
	}
	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		return err
	}

	c.watchConfig()

	return nil
}

func (c *Config) initConfig() error {
	if c.Name != "" {
		// 如果指定了配置文件，则解析指定的配置文件
		viper.SetConfigFile(c.Name)
	} else {
		// 如果没有指定配置文件，则解析默认的配置文件
		viper.AddConfigPath("conf")
		viper.SetConfigName("config")
	}
	// 设置配置文件格式为YAML
	viper.SetConfigType("yaml")
	// viper解析配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

// 监听配置文件是否改变,用于热更新
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s\n", e.Name)
	})
}

// 获取到配置文件
func GetConfig() *Cfgparams {
	if err := Init("./conf/config.yaml"); err != nil {
		panic(err)
	}
	c := &Cfgparams{}
	// name := viper.GetString("name")
	c.HEART_BEAT_TIME = viper.GetInt("common.heart_beat_time")
	c.SERVER_ADDR = viper.GetString("common.server_addr")
	c.DATA_SERVER_TOPIC = viper.GetString("common.server_topic")
	addr := viper.GetString("common.nsqlookupd.addr")
	port := viper.GetString("common.nsqlookupd.http_port")
	c.NSQ_LOOKUPD_ADDR = addr + ":" + port
	x := viper.Get("common.nsqd")
	// c.NSQ_TCP_ADDRS = viper.Get("common.nsqd").([]string)
	var y []string
	for _, v := range x.([]interface{}) {
		y = append(y, v.(string))
	}
	c.NSQ_TCP_ADDRS = y
	// fmt.Println(c)
	return c

}
```
##### config.yaml
```yaml
name:   nsq
common:
    heart_beat_time: 5
    nsqlookupd:
        name:   s1
        addr:   10.8.0.4
        tcp_port:   4160
        http_port:  4161
    nsqadmin:
        name:   s1
        addr:   10.8.0.4
        port:   4171
    nsqd:
    - 10.8.0.4:4150
    - 10.8.0.5:4150
    - 10.8.0.6:4150
    - 10.8.0.7:4150
    server_topic:   "data_server_addr"
    server_addr:   "127.0.0.1:10001"
```