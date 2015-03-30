// CreateListMsg encapsulates the action of a user creating a list
package models

import (
	"encoding/json"
	"errors"

	"github.com/gocql/gocql"
)

type CreateListMsg struct {
	Server Server `json:"server"`
	Client Client `json:"client"`
	User   User   `json:"user"`
	EventFields
	Data CreateListMsgData `json:"data"`
}

// CreateListMsgData wraps any fields specific to this event
type CreateListMsgData struct {
	List List
}

func NewCreateListMsg() CreateListMsg {
	return CreateListMsg{
		EventFields: EventFields{
			Type: CreateListMsgType,
			ID:   gocql.TimeUUID().String(),
		},
	}
}

// DeserializeCreateListMsg deserializes a JSON serialized CreateListMsg struct
func DeserializeCreateListMsg(jsonText []byte) (*CreateListMsg, error) {
	var event CreateListMsg
	err := json.Unmarshal(jsonText, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (e *CreateListMsg) IsReadyToBeSaved() error {
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
	if e.Type != CreateListMsgType {
		return errors.New(InvalidEventTypeError)
	}
	return nil
}
