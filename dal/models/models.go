package models

import (
	"errors"
	"time"
)

const (
	InvalidMsgMethodError   = "Invalid msg method"
	InvalidMsgResourceError = "Invalid msg resource"

	// base msg fields
	MissingClientIDError         = "Msg client id cannot be empty"
	MissingClientDeviceTypeError = "Msg client deviceType cannot be empty"
	MissingUserIDError           = "Msg user id cannot be empty"
	MissingMsgIDError            = "Msg id cannot be empty"
	MissingMsgMethodError        = "Msg method cannot be empty"
	MissingMsgResourceError      = "Msg resource cannot be empty"
	MissingServerHostnameError   = "Msg server hostname cannot be empty"
	MissingServerIPAddressError  = "Msg server ip address cannot be empty"

	MissingListIDError    = "List id cannot be emtpy"
	MissingListTitleError = "List title cannot be empty"
)

type Server struct {
	Hostname  string `json:"hostname"`
	IPAddress string `json:"ipAddress"`
}

type Client struct {
	DeviceType string `json:"deviceType"`
	ID         string `json:"id"`
	OsVersion  string `json:"osVersion"`
}

type List struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Title     string    `json:"title"`
	Users     []string  `json:"users"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Msg struct {
	Method   string `json:"method"`
	Resource string `json:"resource"`
	ID       string `json:"id"`
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
		return errors.New(MissingMsgIDError)
	}
	if len(s.Hostname) == 0 {
		return errors.New(MissingServerHostnameError)
	}
	if len(s.IPAddress) == 0 {
		return errors.New(MissingServerIPAddressError)
	}
	return nil
}
