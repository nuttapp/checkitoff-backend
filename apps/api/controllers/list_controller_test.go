package controllers

import (
	"testing"

	"github.com/bitly/go-nsq"
	"github.com/nuttapp/checkitoff-backend/apps/api"
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
	ctx := api.NewContext("test")
	ctx.Cfg.NSQ.PubTopic = "api_messages"
	ctx.Cfg.NSQ.SubTopic = "api_messages"
	ctx.Cfg.NSQ.SubChannel = "api_messages"

	Convey("Should enqueue ListMsg on NSQ", t, func() {
		apiServer := api.APIServer{Ctx: ctx}
		apiServer.Start()
		// apiServer.CreateTopic("Test_ListController")

		err := ListControllerCreate(createJSON, ctx)
		So(err, ShouldBeNil)

		apiServer.Stop()

		// const Test_ListController = "Test_ListController"
		// const NSQChannel = "test"
		// consumer, err := nsq.NewConsumer(ctx.Cfg.NSQ.PubTopic, NSQChannel, ctx.NSQCfg)
		// So(err, ShouldBeNil)
		// //
		// var logBuf bytes.Buffer
		// logger := log.New(&logBuf, "", log.LstdFlags)
		// consumer.SetLogger(logger, nsq.LogLevelDebug)
		//
		// th := &testHandler{
		// 	testChan: make(chan *nsq.Message),
		// }
		// consumer.AddHandler(th)
		// err = consumer.ConnectToNSQLookupd(ctx.Cfg.NSQ.LookupdHTTPAddr)
		// So(err, ShouldBeNil)
		//
		// dequeuedMsg := <-th.testChan
		// msg, err := dal.DeserializeListMsg(dequeuedMsg.Body)
		// So(err, ShouldBeNil)
		// So(msg, ShouldNotBeNil)
		// So(msg.Data.Title, ShouldEqual, "Trader Joes")
		// So(msg.Data.ID, ShouldNotBeEmpty)
		// So(msg.ID, ShouldNotBeEmpty)
		//
		// // waiting for consumer to stop take about 10ms
		// consumer.Stop()
		// <-consumer.StopChan
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
