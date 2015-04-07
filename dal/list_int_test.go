package dal

import (
	"testing"

	"github.com/gocql/gocql"
	m "github.com/nuttapp/checkitoff-backend/dal/models"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_DAL_List_int(t *testing.T) {
	cfg := DALCfg{
		Hosts:       []string{"127.0.0.1"},
		Keyspace:    "demodb",
		Consistency: gocql.Quorum,
	}
	d, _ := NewDAL(cfg)

	Convey("Given a ListMsg struct", t, func() {
		var msg = &m.ListMsg{
			Servers: []m.Server{
				m.Server{
					Hostname:  "localhost",
					IPAddress: "127.0.0.1",
				},
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
				ID:       "00fff6c8-d9f9-43f4-8d11-dff41ab9008a",
				Title:    "Trader Joes😗 ",
				Category: "shopping",
				IsHidden: false,
				Users:    []string{"0C6B2206-CAD1-4F9C-BE17-74977A257877"},
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
