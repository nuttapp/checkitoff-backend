package controllers

import (
	"encoding/json"
	"testing"

	"github.com/bitly/go-nsq"
	"github.com/nuttapp/checkitoff-backend/api/config"
	m "github.com/nuttapp/checkitoff-backend/common/models"
	. "github.com/smartystreets/goconvey/convey"
)

var createListEventJSON = []byte(`
	{
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
		},
		"method":  "create",
		"resource": "list"
	}`)

var apiCfg = &config.Config{
	Hostname:           "Kemist",
	IPAddress:          "127.0.0.1",
	NSQProducerTCPAddr: "127.0.0.1:4150",
}

func Test_ListController_Integration(t *testing.T) {
	Convey("Should enqueue CreateListMsg on NSQ", t, func() {
		nsqCfg := nsq.NewConfig()
		err := ListControllerCreate(createListEventJSON, nsqCfg, apiCfg)
		So(err, ShouldEqual, nil)
	})
}

func Test_ListController_Unit(t *testing.T) {
	Convey("ListControllerCreate", t, func() {
		var cle = m.CreateListMsg{
			Client: m.Client{
				ID:         "C0085E97-9BCB-4BAA-9A1A-C5AFF37FF433",
				DeviceType: "iPhone",
				OsVersion:  "7,2",
			},
			User: m.User{
				ID:   "0C6B2206-CAD1-4F9C-BE17-74977A257877",
				Name: "test-user-name",
			},
			Msg: m.Msg{
				Method:   "create",
				Resource: "list",
			},
			Data: m.List{
				Title: "Trader Joes",
			},
		}
		Convey("should return error when missing client id", func() {
			e := cle // create copy
			e.Client.ID = ""
			var b, _ = json.Marshal(e)
			err := ListControllerCreate(b, nil, apiCfg)
			So(err.Error(), ShouldContainSubstring, m.MissingClientIDError)
		})
		Convey("should return error when missing client device type", func() {
			e := cle // create copy
			e.Client.DeviceType = ""
			var b, _ = json.Marshal(e)
			err := ListControllerCreate(b, nil, apiCfg)
			So(err.Error(), ShouldContainSubstring, m.MissingClientDeviceTypeError)
		})
		Convey("should return error when missing user id", func() {
			e := cle // create copy
			e.User.ID = ""
			var b, _ = json.Marshal(e)
			err := ListControllerCreate(b, nil, apiCfg)
			So(err.Error(), ShouldContainSubstring, m.MissingUserIDError)
		})
		Convey("should return error when missing method", func() {
			e := cle // create copy
			e.Method = ""
			var b, _ = json.Marshal(e)
			err := ListControllerCreate(b, nil, apiCfg)
			So(err.Error(), ShouldContainSubstring, m.MissingMsgMethodError)
		})
		Convey("should return error when missing resource", func() {
			e := cle // create copy
			e.Resource = ""
			var b, _ = json.Marshal(e)
			err := ListControllerCreate(b, nil, apiCfg)
			So(err.Error(), ShouldContainSubstring, m.MissingMsgResourceError)
		})
		Convey("should return error when invalid method", func() {
			e := cle // create copy
			e.Method = "poop"
			var b, _ = json.Marshal(e)
			err := ListControllerCreate(b, nil, apiCfg)
			So(err.Error(), ShouldContainSubstring, m.InvalidMsgMethodError)
		})
		Convey("should return error when missing list title", func() {
			e := cle // create copy
			e.Data.Title = ""
			var b, _ = json.Marshal(e)
			err := ListControllerCreate(b, nil, apiCfg)
			So(err.Error(), ShouldContainSubstring, m.MissingListTitleError)
		})
		Convey("should return error when missing server hostname", func() {
			var b, _ = json.Marshal(cle)
			apiCfg := &config.Config{
				Hostname:  "",
				IPAddress: "127.0.0.1",
			}
			err := ListControllerCreate(b, nil, apiCfg)
			So(err.Error(), ShouldContainSubstring, m.MissingServerHostnameError)
		})
		Convey("should return error when missing server ip address", func() {
			var b, _ = json.Marshal(cle)
			apiCfg := &config.Config{
				Hostname:  "localhost",
				IPAddress: "",
			}
			err := ListControllerCreate(b, nil, apiCfg)
			So(err.Error(), ShouldContainSubstring, m.MissingServerIPAddressError)
		})
	})

	Convey("Should deserialize json into CreateListMsg", t, func() {
		event, err := m.DeserializeCreateListMsg(createListEventJSON)
		// fmt.Println(util.ToJSON(event))

		So(event, ShouldNotBeNil)
		So(err, ShouldBeNil)
		// Purposely left out id, it should show up as an empty value
		// because the struct isn't using pointers to fields
		So(event.Data.ID, ShouldEqual, "")
		So(event.Data.Title, ShouldEqual, "Trader Joes")
		So(event.User.Name, ShouldEqual, "Greg")
	})
}
