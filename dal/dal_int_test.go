package dal

import (
	"testing"

	"github.com/gocql/gocql"
	m "github.com/nuttapp/checkitoff-backend/dal/models"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_DAL_List(t *testing.T) {
	cfg := DALCfg{
		Hosts:       []string{"127.0.0.1"},
		Keyspace:    "demodb",
		Consistency: gocql.Quorum,
	}
	d, err := NewDAL(cfg)

	Convey("List", t, func() {
		So(err, ShouldBeNil)

		var msg = &m.ListMsg{
			Server: m.Server{
				Hostname:  "localhost",
				IPAddress: "127.0.0.1",
			},
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
				ID:       gocql.TimeUUID().String(),
				Method:   m.MsgMethodCreate,
				Resource: m.MsgResourceList,
			},
			Data: m.List{
				ID:       "00FFF6C8-D9F9-43F4-8D11-DFF41AB9008A",
				Title:    "Trader Joes😗",
				Category: "shopping",
				IsHidden: false,
				Users:    []string{"0C6B2206-CAD1-4F9C-BE17-74977A257877"},
			},
		}

		So(d, ShouldNotBeNil)

		Convey("save/get", func() {
			err = d.CreateList(msg)
			So(err, ShouldBeNil)

			list, err := d.GetList(msg)
			So(err, ShouldBeNil)
			So(list, ShouldNotBeNil)
			So(msg.Data.ID, ShouldEqual, "00FFF6C8-D9F9-43F4-8D11-DFF41AB9008A")
			So(msg.Data.Title, ShouldEqual, msg.Data.Title)
			So(msg.Data.Category, ShouldEqual, msg.Data.Category)
			So(msg.Data.IsHidden, ShouldEqual, false)
			So(msg.Data.Users, ShouldResemble, msg.Data.Users)
		})

		Convey("update", func() {
			d.UpdateList(msg)
			So(err, ShouldBeNil)
		})

		Convey("delete", func() {
			// err := SaveCreateListMsg(cle)
			d.CreateList(msg)
			So(err, ShouldBeNil)
		})
	})
}
