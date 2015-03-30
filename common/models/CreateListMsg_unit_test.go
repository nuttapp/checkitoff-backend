package models

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_CreateListMsg_unit(t *testing.T) {

	Convey("NewCreateListMsg()", t, func() {
		msg := NewCreateListMsg()
		So(msg, ShouldNotBeNil)
	})

	// CreateListMsg.IsReadyToBeSaved()
	Convey("IsReadyToBeSaved()", t, func() {
		msg := &CreateListMsg{
			EventFields: EventFields{
				ID:   "create-list-msg-id",
				Type: CreateListMsgType,
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
			Data: CreateListMsgData{
				List{
					ID:    "create-list-msg-list-id",
					Title: "create-list-msg-list-title",
				},
			},
		}

		Convey("Should return no error with base fields", func() {
			err := msg.IsReadyToBeSaved()
			So(err, ShouldBeNil)
		})

		Convey("should return errror when", func() {
			Convey("list id is blank", func() {
				msg.Data.List.ID = ""
				err := msg.IsReadyToBeSaved()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, MissingListIDError)
			})
			Convey("list title is blank", func() {
				msg.Data.List.Title = ""
				err := msg.IsReadyToBeSaved()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, MissingListTitleError)
			})
			Convey("one of the base fields is blank", func() {
				msg.ID = ""
				err := msg.IsReadyToBeSaved()
				So(err, ShouldNotBeNil)
			})
			Convey("name of message is wrong", func() {
				msg.Type = "invalid-name-for-create-msg"
				err := msg.IsReadyToBeSaved()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, InvalidEventTypeError)
			})
		})
	})
}
