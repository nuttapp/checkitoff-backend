package models

import "errors"

const (
	// Catch fat finger mistakes where an event type is not equal to it's intended event type
	// ex: a CreateListMsg should never have an EventType "UpdateListEvent"
	InvalidMsgMethodError   = "Invalid msg method"
	InvalidMsgResourceError = "Invalid msg resource"

	// base event fields
	MissingClientIDError         = "Msg client id cannot be empty"
	MissingClientDeviceTypeError = "Msg client deviceType cannot be empty"
	MissingUserIDError           = "Msg user id cannot be empty"
	MissingEventIDError          = "Msg id cannot be empty"
	MissingMsgMethodError        = "Msg method cannot be empty"
	MissingMsgResourceError      = "Msg resource cannot be empty"
	MissingServerHostnameError   = "Msg server hostname cannot be empty"
	MissingServerIPAddressError  = "Msg server ip address cannot be empty"

	MissingListIDError    = "List id cannot be emtpy"
	MissingListTitleError = "List title cannot be empty"
)

type Event interface {
	IsReadyToBeSaved() error
}

// ValidateEvent examples the fields that every event needs and throws an error if they're blank
func ValidateMsg(c *Client, u *User, e *Msg, s *Server) error {
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
		return errors.New(MissingMsgMethodError)
	}
	if len(e.Resource) == 0 {
		return errors.New(MissingMsgResourceError)
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
