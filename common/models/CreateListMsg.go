// CreateListMsg encapsulates the action of a user creating a list
package models

import "errors"

type CreateListMsg struct {
	Server ServerFields `json:"server"`
	Client ClientFields `json:"client"`
	User   UserFields   `json:"user"`
	EventFields
	Data CreateListMsgData `json:"data"`
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
