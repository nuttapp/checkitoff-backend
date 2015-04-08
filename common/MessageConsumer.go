package common

import (
	"fmt"
	"log"
	"os"

	"github.com/bitly/go-nsq"
	"github.com/nuttapp/checkitoff-backend/common/util"
)

const (
	FailedToCreateNSQConsumerError = "Failed to create NSQ consumer"
	FailedToConnectToLookupdError  = "Failed to connect NSQ consumer to lookupd"
)

type ConsumerConfig struct {
	Topic           string // The name of the master topic
	Channel         string // The channel name to use when consuming 1 topic -> N Channels
	LookupdHTTPaddr string // ex: 127.0.0.1:4161
	Concurrency     int    // # of goroutine to spawn for nsq handlers
}

func NewMessageConsumer(cfg *ConsumerConfig, nsqCfg *nsq.Config, handler nsq.Handler) *MessageConsumer {
	consumer := MessageConsumer{
		StopChan: make(chan error),
		Handler:  handler,
		Logger:   log.New(os.Stdout, "", log.LstdFlags),

		requestStopChan: make(chan int),
		nsqCfg:          nsqCfg,
		cfg:             cfg,
	}
	return &consumer
}

// MessageConsumer is a simple wrapper around nsq.Consumer, it mostly
// handles logging, and passing config params. It expects
// that you create your own nsq.Handler
type MessageConsumer struct {
	StopChan chan error  // listen on StopChan to block after calling Stop()
	Handler  nsq.Handler // the struct that will recieve the messages
	Logger   *log.Logger

	requestStopChan chan int
	nsqCfg          *nsq.Config
	cfg             *ConsumerConfig
}

func (mc *MessageConsumer) Start() {
	mc.log("Starting...")
	go mc.ListenForMessages()
}

func (mc *MessageConsumer) Stop() {
	mc.log("Attemping shutdown...")
	mc.requestStopChan <- 1
}

func (mc *MessageConsumer) log(text interface{}) {
	mc.Logger.Printf("[38;5;139mMessageConsumer[0m %s", text)
}

func (mc *MessageConsumer) ListenForMessages() {
	consumer, err := nsq.NewConsumer(mc.cfg.Topic, mc.cfg.Channel, mc.nsqCfg)
	if err != nil {
		err := NewError(util.FailedToCreateNSQConsumerError, err)
		mc.StopChan <- err
		mc.log(err)
		return
	}

	// consumer.AddHandler(mc)
	consumer.AddConcurrentHandlers(mc.Handler, mc.cfg.Concurrency)
	// consumer.SetLogger(log.New(ioutil.Discard, "", log.LstdFlags), nsq.LogLevelInfo)

	err = consumer.ConnectToNSQLookupd(mc.cfg.LookupdHTTPaddr)
	if err != nil {
		err = NewError(util.FailedToConnectToLookupdError, err)
		mc.StopChan <- err
		mc.log(err)
		return
	}

	// ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-mc.requestStopChan:
			// consumer.stop is async so we need to block on consumer.stopChan,
			// which gives our nsq consumer time to gracefully shut down.
			// This isn't strictly necessary but it prevents all inflight (msg currently being
			// handled in goroutines) from timing out.
			consumer.Stop()
		case <-consumer.StopChan:
			mc.log("Shutdown completed")
			mc.StopChan <- nil
			return
			// case <-ticker.C:
			// fmt.Print(".")
		}
	}
}

func NewError(baseText string, err error) error {
	return fmt.Errorf("%s: %s", baseText, err)
}
