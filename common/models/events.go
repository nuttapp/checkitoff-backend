package models

import "errors"

const (
	// Catch fat finger mistakes where an event type is not equal to it's intended event type
	// ex: a CreateListMsg should never have an EventType "UpdateListEvent"
	InvalidEventTypeError = "Invalid event type"

	// base event fields
	MissingClientIDError         = "Client id cannot be empty"
	MissingClientDeviceTypeError = "Client deviceType cannot be empty"
	MissingUserIDError           = "User id cannot be empty"
	MissingEventIDError          = "Event id cannot be empty"
	MissingEventTypeError        = "Event type cannot be empty"
	MissingServerHostnameError   = "Server hostname cannot be empty"
	MissingServerIPAddressError  = "Server ip address cannot be empty"

	MissingListIDError    = "List id cannot be emtpy"
	MissingListTitleError = "List title cannot be empty"

	CreateListMsgType = "create-list"
)

type Event interface {
	IsReadyToBeSaved() error
}

// ValidateEvent examples the fields that every event needs and throws an error if they're blank
func ValidateEvent(c *Client, u *User, e *EventFields, s *Server) error {
	if len(c.ID) == 0 {
		return errors.New(MissingClientIDError)
	}
	if len(c.DeviceType) == 0 {
		return errors.New(MissingClientDeviceTypeError)
	}
	if len(u.ID) == 0 {
		return errors.New(MissingUserIDError)
	}
	if len(e.Type) == 0 {
		return errors.New(MissingEventTypeError)
	}
	if len(e.ID) == 0 {
		return errors.New(MissingEventIDError)
	}
	if len(s.Hostname) == 0 {
		return errors.New(MissingServerHostnameError)
	}
	if len(s.IPAddress) == 0 {
		return errors.New(MissingServerIPAddressError)
	}
	return nil
}
