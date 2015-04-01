// CreateListMsg encapsulates the action of a user creating a list
package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gocql/gocql"
)

type CreateListMsg struct {
	Server Server `json:"server"`
	Client Client `json:"client"`
	User   User   `json:"user"`
	EventFields
	Data CreateListMsgData `json:"data"`
}

type UpdateListMsg struct {
	Server Server `json:"server"`
	Client Client `json:"client"`
	User   User   `json:"user"`
	EventFields
	Data CreateListMsgData `json:"data"`
}

// updates fields of a particular list
type UpdateListMsgData struct {
	ID    *string
	Title *string
	Users *[]string
}

// CreateListMsgData wraps any fields specific to this event
type CreateListMsgData struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string
	Users     []string
}

func NewCreateListMsg() CreateListMsg {
	createdAt := time.Now().UTC()
	return CreateListMsg{
		EventFields: EventFields{
			Type: CreateListMsgType,
			ID:   gocql.TimeUUID().String(),
		},
		Data: CreateListMsgData{
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
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
	if len(e.Data.ID) == 0 {
		return errors.New(MissingListIDError)
	}
	if len(e.Data.Title) == 0 {
		return errors.New(MissingListTitleError)
	}
	if e.Type != CreateListMsgType {
		return errors.New(InvalidEventTypeError)
	}
	return nil
}
