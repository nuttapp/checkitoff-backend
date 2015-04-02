package main

import (
	"testing"

	"github.com/gocql/gocql"
	m "github.com/nuttapp/checkitoff-backend/common/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMain_Integration(t *testing.T) {
	var cle = &m.CreateListMsg{
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
			ID:    "00FFF6C8-D9F9-43F4-8D11-DFF41AB9008A",
			Title: "Trader Joes",
		},
	}

	Convey("SaveCreateListMsg", t, func() {
		cluster := gocql.NewCluster("127.0.0.1")
		cluster.Keyspace = "demodb"
		cluster.Consistency = gocql.Quorum
		session, _ := cluster.CreateSession()
		defer session.Close()

		Convey("Should save to database", func() {
			err := SaveCreateListMsg(cle)
			So(err, ShouldBeNil)
		})
	})
}
