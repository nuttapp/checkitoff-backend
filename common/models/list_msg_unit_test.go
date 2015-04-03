package models

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_DeserializeCreateListMsg_unit(t *testing.T) {
	Convey("DeserializeCreateListMsg()", t, func() {
		var createListEventJSON = []byte(`
			{
				"method":  "create",
				"resource": "list",
				"server": {
					"hostname": "localhost",
					"ipAddress": "127.0.0.1"
				},
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
					"id": "12345",
					"title": "Trader Joes"
				}
			}`)

		Convey("should deserialize json into CreateListMsg", func() {
			msg, err := DeserializeCreateListMsg(createListEventJSON)

			So(err, ShouldBeNil)
			So(msg, ShouldNotBeNil)
			So(msg.Method, ShouldEqual, "create")
			So(msg.Resource, ShouldEqual, "list")
			So(msg.Server.Hostname, ShouldEqual, "localhost")
			So(msg.Server.IPAddress, ShouldEqual, "127.0.0.1")
			So(msg.Client.ID, ShouldEqual, "C0085E97-9BCB-4BAA-9A1A-C5AFF37FF433")
			So(msg.Client.DeviceType, ShouldEqual, "iPhone")
			So(msg.Client.OsVersion, ShouldEqual, "7,2")
			So(msg.User.ID, ShouldEqual, "0C6B2206-CAD1-4F9C-BE17-74977A257877")
			So(msg.User.Name, ShouldEqual, "Greg")
			So(msg.Data.ID, ShouldEqual, "12345")
			So(msg.Data.Title, ShouldEqual, "Trader Joes")
		})
		Convey("should return error with invalid json", func() {
			msg, err := DeserializeCreateListMsg([]byte("asdf"))
			So(err, ShouldNotBeNil)
			So(msg, ShouldBeNil)
		})
	})
}

func Test_CreateListMsg_unit(t *testing.T) {

	Convey("NewCreateListMsg()", t, func() {
		msg := NewCreateListMsg()
		So(msg, ShouldNotBeNil)
	})

	// CreateListMsg.IsReadyToBeSaved()
	Convey("IsReadyToBeSaved()", t, func() {
		msg := &CreateListMsg{
			Msg: Msg{
				ID:       "create-list-msg-id",
				Method:   "create",
				Resource: "list",
			},
			Server: Server{
				Hostname:  "create-list-msg-server-hostname",
				IPAddress: "create-list-msg-server-ipaddress",
			},
			Client: Client{
				ID:         "create-list-msg-client-id",
				DeviceType: "create-list-msg-client-deveice-type",
			},
			User: User{
				ID: "create-list-msg-user-id",
			},
			Data: List{
				ID:    "create-list-msg-list-id",
				Title: "create-list-msg-list-title",
			},
		}

		Convey("Should return no error with base fields", func() {
			err := msg.IsReadyToBeSaved()
			So(err, ShouldBeNil)
		})

		Convey("should return errror when", func() {
			Convey("list id is blank", func() {
				msg.Data.ID = ""
				err := msg.IsReadyToBeSaved()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, MissingListIDError)
			})
			Convey("list title is blank", func() {
				msg.Data.Title = ""
				err := msg.IsReadyToBeSaved()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, MissingListTitleError)
			})
			Convey("one of the base fields is blank", func() {
				msg.ID = ""
				err := msg.IsReadyToBeSaved()
				So(err, ShouldNotBeNil)
			})
			Convey("msg method is not valid", func() {
				msg.Method = "invalid method type"
				err := msg.IsReadyToBeSaved()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, InvalidMsgMethodError)
			})
			Convey("msg resource is not \"list\"", func() {
				msg.Resource = "invalid resource"
				err := msg.IsReadyToBeSaved()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, InvalidMsgResourceError)
			})
		})
	})
}
