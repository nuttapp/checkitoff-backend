package models

import (
	"encoding/json"
	"errors"

	"github.com/gocql/gocql"
)

const (
	// Catch fat finger mistakes where an event type is not equal to it's intended event type
	// ex: a CreateListEvent should never have an EventType "UpdateListEvent"
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

	CreateListEventType = "create-list"
)

type Event interface {
	IsReadyToBeSaved() error
}

func NewCreateListEvent() CreateListEvent {
	return CreateListEvent{
		EventFields: EventFields{
			Type: CreateListEventType,
			ID:   gocql.TimeUUID().String(),
		},
	}
}

// CreateListEvent encapsulates the action of a user creating a list
type CreateListEvent struct {
	Server ServerFields `json:"server"`
	Client ClientFields `json:"client"`
	User   UserFields   `json:"user"`
	EventFields
	Data CreateListEventData `json:"data"`
}

func (e *CreateListEvent) IsReadyToBeSaved() error {
	err := ValidateEvent(&e.Client, &e.User, &e.EventFields, &e.Server)
	if err != nil {
		return err
	}
	if len(e.Data.List.ID) == 0 {
		return errors.New(MissingListIDError)
	}
	if len(e.Data.List.Title) == 0 {
		return errors.New(MissingListTitleError)
	}
	if e.Type != CreateListEventType {
		return errors.New(InvalidEventTypeError)
	}
	return nil
}

func ValidateEvent(c *ClientFields, u *UserFields, e *EventFields, s *ServerFields) error {
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

// CreateListEventData wraps any fields specific to this event
type CreateListEventData struct {
	List ListFields
}

// DeserializeCreateListEvent deserializes a JSON serialized CreateListEvent struct
func DeserializeCreateListEvent(jsonText []byte) (*CreateListEvent, error) {
	var event CreateListEvent
	err := json.Unmarshal(jsonText, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}
