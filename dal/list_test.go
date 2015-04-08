package dal

import (
	"testing"
	"time"

	"github.com/gocql/gocql"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_DeserializeistMsg_unit(t *testing.T) {
	Convey("DeserializeListMsg()", t, func() {
		var createListMsgJSON = []byte(`
			{
				"method":  "create",
				"resource": "list",
				"servers": [{
					"ttl": 1,
					"hostname": "localhost",
					"ipAddress": "127.0.0.1"
				}],
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
					"title": "Trader Joes",
					"category": "shopping",
					"isHidden": true,
					"users": ["greg", "jess"]
				}
			}`)

		Convey("should deserialize json into ListMsg", func() {
			msg, err := DeserializeListMsg(createListMsgJSON)

			So(err, ShouldBeNil)
			So(msg, ShouldNotBeNil)
			So(msg.Method, ShouldEqual, "create")
			So(msg.Resource, ShouldEqual, "list")
			So(msg.Servers[0].Hostname, ShouldEqual, "localhost")
			So(msg.Servers[0].IPAddress, ShouldEqual, "127.0.0.1")
			So(msg.Servers[0].TTL, ShouldEqual, 1)
			So(msg.Client.ID, ShouldEqual, "C0085E97-9BCB-4BAA-9A1A-C5AFF37FF433")
			So(msg.Client.DeviceType, ShouldEqual, "iPhone")
			So(msg.Client.OsVersion, ShouldEqual, "7,2")
			So(msg.User.ID, ShouldEqual, "0C6B2206-CAD1-4F9C-BE17-74977A257877")
			So(msg.User.Name, ShouldEqual, "Greg")
			So(msg.Data.ID, ShouldEqual, "12345")
			So(msg.Data.Title, ShouldEqual, "Trader Joes")
			So(msg.Data.IsHidden, ShouldEqual, true)
			So(msg.Data.Category, ShouldEqual, "shopping")
		})
		Convey("should return error with invalid json", func() {
			msg, err := DeserializeListMsg([]byte("asdf"))
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
		msg := &ListMsg{
			Msg: Msg{
				ID:       "create-list-msg-id",
				Method:   "create",
				Resource: "list",
			},
			Servers: []Server{
				Server{
					Hostname:  "create-list-msg-server-hostname",
					IPAddress: "create-list-msg-server-ipaddress",
				},
			},
			Client: Client{
				ID:         "create-list-msg-client-id",
				DeviceType: "create-list-msg-client-deveice-type",
			},
			User: User{
				ID: "create-list-msg-user-id",
			},
			Data: List{
				ID:       "create-list-msg-list-id",
				Title:    "create-list-msg-list-title",
				Category: "shopping",
				IsHidden: false,
			},
		}

		// b, _ := json.MarshalIndent(msg, "", "  ")
		// fmt.Printf("%s", string(b))

		Convey("Should return no error with base fields", func() {
			err := msg.ValidateMsg()
			So(err, ShouldBeNil)
		})

		Convey("should return errror when", func() {
			Convey("list id is blank", func() {
				msg.Data.ID = ""
				err := msg.ValidateMsg()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, MissingListIDError)
			})
			Convey("list title is blank", func() {
				msg.Data.Title = ""
				err := msg.ValidateMsg()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, MissingListTitleError)
			})
			Convey("one of the base fields is blank", func() {
				msg.ID = ""
				err := msg.ValidateMsg()
				So(err, ShouldNotBeNil)
			})
			Convey("msg method is not valid", func() {
				msg.Method = "invalid method type"
				err := msg.ValidateMsg()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, InvalidMsgMethodError)
			})
			Convey("msg resource is not \"list\"", func() {
				msg.Resource = "invalid resource"
				err := msg.ValidateMsg()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, InvalidMsgResourceError)
			})
		})
	})
}

func Test_DAL_List_int(t *testing.T) {
	cfg := DALCfg{
		Hosts:       []string{"127.0.0.1"},
		Keyspace:    "demodb",
		Consistency: gocql.Quorum,
	}
	d, _ := NewDAL(cfg)

	Convey("Given a ListMsg struct", t, func() {
		var msg = &ListMsg{
			Servers: []Server{
				Server{
					Hostname:  "localhost",
					IPAddress: "127.0.0.1",
				},
			},
			Client: Client{
				ID:         "C0085E97-9BCB-4BAA-9A1A-C5AFF37FF433",
				DeviceType: "iPhone",
				OsVersion:  "7,2",
			},
			User: User{
				ID:   "0C6B2206-CAD1-4F9C-BE17-74977A257877",
				Name: "test-user-name",
			},
			Msg: Msg{
				ID:       gocql.TimeUUID().String(),
				Method:   MsgMethodCreate,
				Resource: MsgResourceList,
			},
			Data: List{
				ID:        "00fff6c8-d9f9-43f4-8d11-dff41ab9008a",
				Title:     "Trader Joes😗 ",
				Category:  "shopping",
				IsHidden:  false,
				Users:     []string{"0C6B2206-CAD1-4F9C-BE17-74977A257877"},
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			},
		}

		Convey("should save the List", func() {
			err := d.CreateOrUpdateList(msg)
			So(err, ShouldBeNil)

			list, err := d.GetList(msg)
			So(err, ShouldBeNil)
			So(list, ShouldNotBeNil)
			So(list.ID, ShouldEqual, msg.Data.ID)
			So(list.Title, ShouldEqual, msg.Data.Title)
			So(list.Category, ShouldEqual, msg.Data.Category)
			So(list.IsHidden, ShouldEqual, false)
			So(list.Users, ShouldResemble, msg.Data.Users)
			So(list.CreatedAt, ShouldHappenOnOrBefore, msg.Data.CreatedAt)
			So(list.UpdatedAt, ShouldHappenOnOrBefore, msg.Data.UpdatedAt)
		})

		Convey("should update the List", func() {
			msg.Data.Title = "Target😡  "
			err := d.CreateOrUpdateList(msg)
			So(err, ShouldBeNil)

			list, err := d.GetList(msg)
			So(err, ShouldBeNil)
			So(list, ShouldNotBeNil)
			So(list.ID, ShouldEqual, msg.Data.ID)
			So(list.Title, ShouldEqual, msg.Data.Title)
			So(list.Category, ShouldEqual, msg.Data.Category)
			So(list.IsHidden, ShouldEqual, false)
			So(list.Users, ShouldResemble, msg.Data.Users)
			So(list.CreatedAt, ShouldHappenOnOrBefore, msg.Data.CreatedAt)
			So(list.UpdatedAt, ShouldHappenOnOrBefore, msg.Data.UpdatedAt)
		})

		Convey("should delete the List", func() {
			err := d.DeleteList(msg)
			So(err, ShouldBeNil)

			list, err := d.GetList(msg)
			So(err.Error(), ShouldContainSubstring, gocql.ErrNotFound.Error())
			So(list, ShouldBeNil)
		})
	})
}
