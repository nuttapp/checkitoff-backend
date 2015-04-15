package controllers

import (
	"bytes"
	"log"
	"testing"

	"github.com/bitly/go-nsq"
	"github.com/nuttapp/checkitoff-backend/apps/api/config"
	"github.com/nuttapp/checkitoff-backend/dal"
	. "github.com/smartystreets/goconvey/convey"
)

var createJSON = []byte(`
	{
		"method":  "create",
		"resource": "list",
		"client": {
			"id": "C0085E97-9BCB-4BAA-9A1A-C5AFF37FF433",
			"deviceType": "iPhone",
			"osVersion": "7,2"
		},
		"user": {
			"id": "0C6B2206-CAD1-4F9C-BE17-74977A257877",
			"name": "Greg"
		},
		"data": {
			"title": "Trader Joes"
		}
	}`)

func Test_ListController_int(t *testing.T) {
	const (
		NSQTopic           = "Test_ListController"
		NSQChannel         = "test"
		NSQLookupdHTTPAddr = "127.0.0.1:4161"
	)

	apiCfg := &config.Config{
		Hostname:           "localhost",
		IPAddress:          "127.0.0.1",
		Role:               "api",
		NSQProducerTCPAddr: "127.0.0.1:4150",
		NSQTopic:           "Test_ListController",
	}

	nsqCfg := nsq.NewConfig()

	Convey("Should enqueue/dequeue ListMsg on NSQ", t, func() {
		err := ListControllerCreate(createJSON, nsqCfg, apiCfg)
		So(err, ShouldBeNil)

		th := &testHandler{
			testChan: make(chan *nsq.Message),
		}

		consumer, err := nsq.NewConsumer(NSQTopic, NSQChannel, nsqCfg)
		So(err, ShouldBeNil)

		var logBuf bytes.Buffer
		logger := log.New(&logBuf, "", log.LstdFlags)
		consumer.SetLogger(logger, nsq.LogLevelDebug)

		consumer.AddHandler(th)
		err = consumer.ConnectToNSQLookupd(NSQLookupdHTTPAddr)
		So(err, ShouldBeNil)

		dequeuedMsg := <-th.testChan
		msg, err := dal.DeserializeListMsg(dequeuedMsg.Body)
		So(err, ShouldBeNil)
		So(msg, ShouldNotBeNil)
		So(msg.Data.Title, ShouldEqual, "Trader Joes")
		So(msg.Data.ID, ShouldNotBeEmpty)
		So(msg.ID, ShouldNotBeEmpty)

		// waiting for consumer to stop take about 10ms
		consumer.Stop()
		<-consumer.StopChan
		// fmt.Println(logBuf.String())
	})
}

type testHandler struct {
	testChan chan *nsq.Message
}

func (mh *testHandler) HandleMessage(msg *nsq.Message) error {
	mh.testChan <- msg
	return nil
}
