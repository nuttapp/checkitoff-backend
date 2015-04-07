// CreateListMsg encapsulates the action of a user creating a list
package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gocql/gocql"
)

type List struct {
	ID        string    `json:"id"`
	Category  string    `json:"category"`
	Title     string    `json:"title"`
	Users     []string  `json:"users"`
	IsHidden  bool      `json:"isHidden"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ListMsg struct {
	Msg
	Servers []Server `json:"servers"`
	Client  Client   `json:"client"`
	User    User     `json:"user"`
	Data    List     `json:"data"`
	Err     error    `json:"error"`
}

func NewCreateListMsg() ListMsg {
	createdAt := time.Now().UTC()
	return ListMsg{
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
func DeserializeCreateListMsg(jsonText []byte) (*ListMsg, error) {
	var event ListMsg
	err := json.Unmarshal(jsonText, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (e *ListMsg) IsReadyToBeSaved() error {
	err := ValidateMsg(&e.Client, &e.User, &e.Msg, e.Servers)
	if err != nil {
		return err
	}
	if len(e.Data.ID) == 0 {
		return errors.New(MissingListIDError)
	}
	if len(e.Data.Title) == 0 {
		return errors.New(MissingListTitleError)
	}
	isValidMethod := e.Method == MsgMethodCreate || e.Method == MsgMethodUpdate || e.Method == MsgMethodDelete
	if !isValidMethod {
		return errors.New(InvalidMsgMethodError)
	}
	if e.Resource != MsgResourceList {
		return errors.New(InvalidMsgResourceError)
	}
	return nil
}
