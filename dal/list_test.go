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
					"ipAddress": "127.0.0.1",
					"role": "api"
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

		Convey("should deserialize json into ListEvent", func() {
			e, err := DeserializeListEvent(createListMsgJSON)

			So(err, ShouldBeNil)
			So(e, ShouldNotBeNil)
			So(e.Method, ShouldEqual, "create")
			So(e.Resource, ShouldEqual, "list")
			So(e.Servers[0].Hostname, ShouldEqual, "localhost")
			So(e.Servers[0].IPAddress, ShouldEqual, "127.0.0.1")
			So(e.Servers[0].TTL, ShouldEqual, 1)
			So(e.Servers[0].Role, ShouldEqual, "api")
			So(e.Client.ID, ShouldEqual, "C0085E97-9BCB-4BAA-9A1A-C5AFF37FF433")
			So(e.Client.DeviceType, ShouldEqual, "iPhone")
			So(e.Client.OsVersion, ShouldEqual, "7,2")
			So(e.User.ID, ShouldEqual, "0C6B2206-CAD1-4F9C-BE17-74977A257877")
			So(e.User.Name, ShouldEqual, "Greg")
			So(e.Data.ID, ShouldEqual, "12345")
			So(e.Data.Title, ShouldEqual, "Trader Joes")
			So(e.Data.IsHidden, ShouldEqual, true)
			So(e.Data.Category, ShouldEqual, "shopping")
			So(e.Data.Users[0], ShouldEqual, "greg")
			So(e.Data.Users[1], ShouldEqual, "jess")
		})
		Convey("should return error with invalid json", func() {
			e, err := DeserializeListEvent([]byte("asdf"))
			So(err, ShouldNotBeNil)
			So(e, ShouldBeNil)
		})
	})
}

func Test_NewListEvent_unit(t *testing.T) {
	Convey("NewListEvent()", t, func() {
		Convey("should generate Event.ID, Data.ID, and users array", func() {
			e, err := NewListEvent(MsgMethodCreate, []byte(`{
				"method": "create",
				"resource": "list",
				"user": {
					"id": "12345"
				},
				"data": {
					"title": "title"
				}
			}`))
			So(err, ShouldBeNil)
			So(e, ShouldNotBeNil)
			So(e.ID, ShouldNotBeEmpty)
			So(e.Data.ID, ShouldNotBeEmpty)
			So(e.Data.Users[0], ShouldEqual, "12345")
			So(e.Data.Title, ShouldEqual, "title")
			So(e.Method, ShouldEqual, MsgMethodCreate)
			So(e.Resource, ShouldEqual, MsgResourceList)
			So(e.Data.UpdatedAt, ShouldResemble, e.Data.CreatedAt)
		})
		Convey("should update UpdatedAt & use data.id", func() {
			timestamp := time.Now().UTC()
			e, err := NewListEvent(MsgMethodUpdate, []byte(`{
				"method": "update",
				"resource": "list",
				"data": {
					"id": "12345",
					"title": "title"
				}
			}`))
			So(err, ShouldBeNil)
			So(e, ShouldNotBeNil)
			So(e.Method, ShouldEqual, MsgMethodUpdate)
			So(e.Resource, ShouldEqual, MsgResourceList)
			So(e.Data.ID, ShouldEqual, "12345")
			So(e.Method, ShouldEqual, MsgMethodUpdate)
			So(e.Data.UpdatedAt, ShouldHappenAfter, timestamp)
		})
	})
}

func Test_ListEvent_Validate_unit(t *testing.T) {
	Convey("Validate()", t, func() {
		e := &ListEvent{
			Event: Event{
				ID:       "create-list-event-id",
				Method:   "create",
				Resource: "list",
			},
			Servers: []Server{
				Server{
					TTL:       1,
					Hostname:  "create-list-event-server-hostname",
					IPAddress: "create-list-event-server-ipaddress",
					Role:      "api",
				},
			},
			Client: Client{
				ID:         "create-list-event-client-id",
				DeviceType: "create-list-event-client-deveice-type",
			},
			User: User{
				ID: "create-list-event-user-id",
			},
			Data: List{
				ID:       "create-list-event-list-id",
				Title:    "create-list-event-list-title",
				Category: "shopping",
				IsHidden: false,
			},
		}

		// b, _ := json.MarshalIndent(event, "", "  ")
		// fmt.Printf("%s", string(b))

		Convey("Should return no error with base fields", func() {
			err := e.Validate()
			So(err, ShouldBeNil)
		})

		Convey("should return errror when", func() {
			Convey("list id is blank", func() {
				e.Data.ID = ""
				err := e.Validate()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, MissingListIDError)
			})
			Convey("list title is blank", func() {
				e.Data.Title = ""
				err := e.Validate()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, MissingListTitleError)
			})
			Convey("one of the base fields is blank", func() {
				e.ID = ""
				err := e.Validate()
				So(err, ShouldNotBeNil)
			})
			Convey("event method is not valid", func() {
				e.Method = "invalid method type"
				err := e.Validate()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, InvalidEventMethodError)
			})
			Convey("event resource is not \"list\"", func() {
				e.Resource = "invalid resource"
				err := e.Validate()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, InvalidEventResourceError)
			})
		})
	})
}

func Test_List_CRUD_int(t *testing.T) {
	cfg := DALCfg{
		Hosts:       []string{"127.0.0.1"},
		Keyspace:    "demodb",
		Consistency: gocql.Quorum,
	}
	d, _ := NewDAL(cfg)

	Convey("Given a ListEvent struct", t, func() {
		var e = &ListEvent{
			Servers: []Server{
				Server{
					Hostname:  "localhost",
					IPAddress: "127.0.0.1",
					Role:      "api",
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
			Event: Event{
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

		Convey("CreateOrUpdateList()", func() {
			Convey("should save the list", func() {
				err := d.CreateOrUpdateList(e)
				So(err, ShouldBeNil)

				list, err := d.GetList(e)
				So(err, ShouldBeNil)
				So(list, ShouldNotBeNil)
				So(list.ID, ShouldEqual, e.Data.ID)
				So(list.Title, ShouldEqual, e.Data.Title)
				So(list.Category, ShouldEqual, e.Data.Category)
				So(list.IsHidden, ShouldEqual, false)
				So(list.Users, ShouldResemble, e.Data.Users)
				So(list.CreatedAt, ShouldHappenOnOrBefore, e.Data.CreatedAt)
				So(list.UpdatedAt, ShouldHappenOnOrBefore, e.Data.UpdatedAt)
			})
			Convey("should return error when", func() {
				Convey("list ID is blank", func() {
					e.Data.ID = ""
					err := d.CreateOrUpdateList(e)
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldContainSubstring, MissingListIDError)
				})
			})
		})

		Convey("should update the List", func() {
			e.Data.Title = "Target😡  "
			err := d.CreateOrUpdateList(e)
			So(err, ShouldBeNil)

			list, err := d.GetList(e)
			So(err, ShouldBeNil)
			So(list, ShouldNotBeNil)
			So(list.ID, ShouldEqual, e.Data.ID)
			So(list.Title, ShouldEqual, e.Data.Title)
			So(list.Category, ShouldEqual, e.Data.Category)
			So(list.IsHidden, ShouldEqual, false)
			So(list.Users, ShouldResemble, e.Data.Users)
			So(list.CreatedAt, ShouldHappenOnOrBefore, e.Data.CreatedAt)
			So(list.UpdatedAt, ShouldHappenOnOrBefore, e.Data.UpdatedAt)
		})

		Convey("should delete the List", func() {
			err := d.DeleteList(e)
			So(err, ShouldBeNil)

			list, err := d.GetList(e)
			So(err.Error(), ShouldContainSubstring, gocql.ErrNotFound.Error())
			So(list, ShouldBeNil)
		})
	})
}
