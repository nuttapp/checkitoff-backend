package controllers

import (
	"bytes"
	"log"
	"testing"

	"github.com/bitly/go-nsq"
	"github.com/nuttapp/checkitoff-backend/apps/api/config"
	"github.com/nuttapp/checkitoff-backend/common"
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
	conCfg := &common.ConsumerConfig{
		Topic:           "Test_ListController",
		Channel:         "test",
		LookupdHTTPaddr: "127.0.0.1:4161",
		Concurrency:     1,
	}

	apiCfg := &config.Config{
		Hostname:           "localhost",
		IPAddress:          "127.0.0.1",
		NSQProducerTCPAddr: "127.0.0.1:4150",
		NSQTopic:           "Test_ListController",
	}

	nsqCfg := nsq.NewConfig()

	Convey("Should enqueue CreateListMsg on NSQ", t, func() {
		err := ListControllerCreate(createJSON, nsqCfg, apiCfg)
		So(err, ShouldBeNil)

		th := &testHandler{
			testChan: make(chan *nsq.Message, 5),
		}
		con := common.NewMessageConsumer(conCfg, nsqCfg, th)

		var b bytes.Buffer
		con.Logger = log.New(&b, "", 0)
		con.Start()

		dequeuedMsg := <-th.testChan
		msg, err := dal.DeserializeListMsg(dequeuedMsg.Body)
		So(err, ShouldBeNil)
		So(msg, ShouldNotBeNil)
		So(msg.Data.Title, ShouldEqual, "Trader Joes")
	})

}

type testHandler struct {
	testChan chan *nsq.Message
}

func (mh *testHandler) HandleMessage(msg *nsq.Message) error {
	mh.testChan <- msg
	return nil
}
