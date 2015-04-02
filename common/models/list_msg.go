// CreateListMsg encapsulates the action of a user creating a list
package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gocql/gocql"
)

//go:generate stringer -type=ResType
const (
	// A message generated by an application that is purely informational in nature.

	// A message generated by an application that indicates a successful operation
	// Information messages should not require action on the part of a user or administrator.
	ResTypeSuccess = "success"
	// A message generated by an application that represents an known edge case
	// or condition that requires furthur attention.
	// Warning messages generally require some action on the part of a user or administrator.
	ResTypeWarning = "warning"
	// They generally do not disrupt the overal flow of an application.
	// A KnownError should be used when the exception raised is an expected condition.
	ResTypeKnownError = "known-error"
	// Unknown errors represent exceptions that are not handled gracefully.
	// They generally disrupt the overall flow of an application.
	// An UnkownError should be used when the exception raised is an unexpected condition.
	ResTypeUnkownError = "unknown-error"
)

const (
	MsgMethodCreate = "create"
	MsgMethodUpdate = "update"
	MsgMethodPatch  = "patch"
	MsgMethodDelete = "delete"

	MsgResourceList     = "list"
	MsgResourceListItem = "list-item"
)

type CreateListMsg struct {
	Msg
	Server Server `json:"server"`
	Client Client `json:"client"`
	User   User   `json:"user"`
	Data   List   `json:"data"`
}

type UpdateListMsg struct {
	Msg
	Server Server `json:"server"`
	Client Client `json:"client"`
	User   User   `json:"user"`
	Data   List   `json:"data"`
}

func NewCreateListMsg() CreateListMsg {
	createdAt := time.Now().UTC()
	return CreateListMsg{
		Msg: Msg{
			Method:   MsgMethodCreate,
			Resource: MsgResourceList,
			ID:       gocql.TimeUUID().String(),
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
	err := ValidateEvent(&e.Client, &e.User, &e.Msg, &e.Server)
	if err != nil {
		return err
	}
	if len(e.Data.ID) == 0 {
		return errors.New(MissingListIDError)
	}
	if len(e.Data.Title) == 0 {
		return errors.New(MissingListTitleError)
	}
	if e.Method != "create" {
		return errors.New(InvalidMsgMethodError)
	}
	if e.Resource != "list" {
		return errors.New(InvalidMsgResourceError)
	}
	return nil
}
