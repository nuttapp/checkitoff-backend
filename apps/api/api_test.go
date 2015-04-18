package api

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_API_int(t *testing.T) {
	context := InitContext()
	Convey("API", t, func() {
		Convey("should Ping NSQ", func() {
			err := context.Producer.Ping()
			So(err, ShouldBeNil)
		})
	})
}
