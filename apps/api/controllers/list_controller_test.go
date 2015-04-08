package controllers

import (
	"testing"

	"github.com/bitly/go-nsq"
	"github.com/nuttapp/checkitoff-backend/apps/api/config"
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

func Test_ListController_Integration(t *testing.T) {
	var apiCfg = &config.Config{
		Hostname:           "Kemist",
		IPAddress:          "127.0.0.1",
		NSQProducerTCPAddr: "127.0.0.1:4150",
	}

	Convey("Should enqueue CreateListMsg on NSQ", t, func() {
		nsqCfg := nsq.NewConfig()
		err := ListControllerCreate(createJSON, nsqCfg, apiCfg)
		So(err, ShouldEqual, nil)
	})
}
