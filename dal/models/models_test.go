package models

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Models_unit(t *testing.T) {
	Convey("ValidateMsg()", t, func() {
		c := &Client{
			ID:         "C0085E97-9BCB-4BAA-9A1A-C5AFF37FF433",
			DeviceType: "iPhone",
			OsVersion:  "7,2",
		}
		u := &User{
			ID:   "0C6B2206-CAD1-4F9C-BE17-74977A257877",
			Name: "test-user-name",
		}
		m := &Msg{
			ID:       "b920a5c0-d8f1-11e4-85fc-358ed15ecf72",
			Method:   "create",
			Resource: "list",
		}
		s := []Server{
			Server{
				Hostname:  "localhost",
				IPAddress: "127.0.0.1",
			},
		}

		Convey("should return error when missing client id", func() {
			c.ID = ""
			err := ValidateMsg(c, u, m, s)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, MissingClientIDError)
		})
		Convey("should return error when missing client device type", func() {
			c.DeviceType = ""
			err := ValidateMsg(c, u, m, s)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, MissingClientDeviceTypeError)
		})
		Convey("should return error when missing user id", func() {
			u.ID = ""
			err := ValidateMsg(c, u, m, s)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, MissingUserIDError)
		})
		Convey("should return error when missing method", func() {
			m.Method = ""
			err := ValidateMsg(c, u, m, s)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, MissingMsgMethodError)
		})
		Convey("should return error when missing resource", func() {
			m.Resource = ""
			err := ValidateMsg(c, u, m, s)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, MissingMsgResourceError)
		})
		Convey("should return error when missing server hostname", func() {
			s[0].Hostname = ""
			err := ValidateMsg(c, u, m, s)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, MissingServerHostnameError)
		})
		Convey("should return error when missing server ip address", func() {
			s[0].IPAddress = ""
			err := ValidateMsg(c, u, m, s)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, MissingServerIPAddressError)
		})
		//
		// Convey("Should deserialize json into CreateListMsg", t, func() {
		// 	event, err := m.DeserializeCreateListMsg(createListMsgJSON)
		// 	// fmt.Println(util.ToJSON(event))
		//
		// 	So(event, ShouldNotBeNil)
		// 	So(err, ShouldBeNil)
		// 	// Purposely left out id, it should show up as an empty value
		// 	// because the struct isn't using pointers to fields
		// 	So(event.Data.ID, ShouldEqual, "")
		// 	So(event.Data.Title, ShouldEqual, "Trader Joes")
		// 	So(event.User.Name, ShouldEqual, "Greg")
	})
}
