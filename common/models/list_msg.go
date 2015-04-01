// CreateListMsg encapsulates the action of a user creating a list
package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gocql/gocql"
)

type CreateListMsg struct {
	EventFields
	Server Server `json:"server"`
	Client Client `json:"client"`
	User   User   `json:"user"`
	Data   List   `json:"data"`
}

type UpdateListMsg struct {
	EventFields
	Server Server `json:"server"`
	Client Client `json:"client"`
	User   User   `json:"user"`
	Data   List   `json:"data"`
}

func NewCreateListMsg() CreateListMsg {
	createdAt := time.Now().UTC()
	return CreateListMsg{
		EventFields: EventFields{
			Type: CreateListMsgType,
			ID:   gocql.TimeUUID().String(),
		},
		Data: List{
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
