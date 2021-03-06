package dal

import (
	"errors"

	"github.com/gocql/gocql"
)

const (
	// Base event
	MissingEventIDError       = "Event ID cannot be empty"
	MissingEventMethodError   = "Event Method cannot be empty"
	MissingEventResourceError = "Event Resource cannot be empty"
	InvalidEventMethodError   = "Unsupported Event Method for the given Resource"
	InvalidEventResourceError = "Invalid Event Resource"

	// Client
	MissingClientIDError         = "Event Client ID cannot be empty"
	MissingClientDeviceTypeError = "Event Client DeviceType cannot be empty"
	// User
	MissingUserIDError = "Event User ID cannot be empty"
	// Server
	MissingServersField         = "Event Servers array cannot be empty"
	MissingServerHostnameError  = "Event Server Hostname cannot be empty"
	MissingServerIPAddressError = "Event Server IPAddress cannot be empty"
	MissingServerRoleError      = "Event Server Role cannot be empty"
	// List
	MissingListIDError    = "List ID cannot be emtpy"
	MissingListTitleError = "List Title cannot be empty"
	// ListItem
)

//go:generate stringer -type=ResType
const (
	// A message generated by an application that is purely informational in nature.
	// Information messages should not require action on the part of a user or administrator.
	LogInformation = "information"
	// A message generated by an application that represents an known edge case
	// or condition that requires furthur attention.
	// Warning messages generally require some action on the part of a user or administrator.
	LogWarning = "warning"
	// They generally do not disrupt the overal flow of an application.
	// A KnownError should be used when the exception raised is an expected condition.
	LogKnownError = "known-error"
	// Unknown errors represent exceptions that are not handled gracefully.
	// They generally disrupt the overall flow of an application.
	// An UnkownError should be used when the exception raised is an unexpected condition.
	UnkownError = "unknown-error"
)

const (
	MsgMethodCreate = "create"
	MsgMethodUpdate = "update"
	MsgMethodPatch  = "patch"
	MsgMethodDelete = "delete"

	// for ListItem
	MsgMethodCheck   = "check"
	MsgMethodUncheck = "uncheck"

	MsgResourceList     = "list"
	MsgResourceListItem = "list-item"
)

type Server struct {
	TTL       int    `json:"ttl"`       // incrementing # of the hop
	Hostname  string `json:"hostname"`  //
	IPAddress string `json:"ipAddress"` // ipv4 addr
	Role      string `json:"role"`      // api/persistor/logger/db
}

type Client struct {
	DeviceType string `json:"deviceType"`
	ID         string `json:"id"`
	OsVersion  string `json:"osVersion"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Event struct {
	Method   string `json:"method"`
	Resource string `json:"resource"`
	ID       string `json:"id"`
}

// ValidateEvent examples the fields that every event needs and throws an error if they're blank
func ValidateEvent(c Client, u User, e Event, s []Server) error {
	if len(c.ID) == 0 {
		return errors.New(MissingClientIDError)
	}
	if len(c.DeviceType) == 0 {
		return errors.New(MissingClientDeviceTypeError)
	}
	if len(u.ID) == 0 {
		return errors.New(MissingUserIDError)
	}
	if len(e.Method) == 0 {
		return errors.New(MissingEventMethodError)
	}
	if len(e.Resource) == 0 {
		return errors.New(MissingEventResourceError)
	}
	if len(e.ID) == 0 {
		return errors.New(MissingEventIDError)
	}
	if len(s) == 0 {
		return errors.New(MissingServersField)
	}
	for _, server := range s {
		if len(server.Hostname) == 0 {
			return errors.New(MissingServerHostnameError)
		}
		if len(server.IPAddress) == 0 {
			return errors.New(MissingServerIPAddressError)
		}
		if len(server.Role) == 0 {
			return errors.New(MissingServerRoleError)
		}
	}
	return nil
}

type DALCfg struct {
	Hosts       []string
	Keyspace    string
	Consistency gocql.Consistency
}

type DAL struct {
	Cfg     DALCfg
	session *gocql.Session
}

func NewDAL(cfg DALCfg) (*DAL, error) {
	d := &DAL{
		Cfg: cfg,
	}
	err := d.initSession()
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *DAL) initSession() error {
	if len(d.Cfg.Hosts) == 0 {
		return errors.New("No hosts were provided for cassandra cluster")
	}
	if len(d.Cfg.Keyspace) == 0 {
		return errors.New("No keyspace was provided")
	}
	if d.session == nil {
		cluster := gocql.NewCluster(d.Cfg.Hosts...)
		cluster.Keyspace = d.Cfg.Keyspace
		cluster.Consistency = d.Cfg.Consistency

		session, err := cluster.CreateSession()
		// fmt.Println(session)
		if err != nil {
			return err
		}
		// fmt.Printf("session: %s\n", session)
		d.session = session
	}
	return nil
}

func (d *DAL) HandleListEvent(e *ListEvent) error {
	err := d.initSession()
	if err != nil {
		return err
	}

	switch e.Method {
	case MsgMethodCreate, MsgMethodUpdate:
		return d.CreateOrUpdateList(e)

	}
	return nil
}
