package dal

import (
	"encoding/json"
	"errors"
	"fmt"
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

func NewListMsg(msgMethod string, msgJSON []byte) (*ListMsg, error) {
	msgID, err := gocql.RandomUUID()
	if err != nil {
		return nil, err
	}

	msg, err := DeserializeListMsg(msgJSON)
	if err != nil {
		return nil, err
	}

	switch msgMethod {
	case MsgMethodCreate:
		listID, err := gocql.RandomUUID()
		if err != nil {
			return nil, err
		}

		msg.ID = msgID.String()
		msg.Data.ID = listID.String()
		msg.Data.CreatedAt = time.Now().UTC()
		msg.Data.UpdatedAt = msg.Data.CreatedAt

		// add the user creating the list to users of the list
		msg.Data.Users = append(msg.Data.Users, msg.User.ID)
	case MsgMethodUpdate:
		msg.ID = msgID.String()
		msg.Data.UpdatedAt = time.Now().UTC()
	}

	return msg, nil
}

// DeserializeCreateListMsg deserializes a JSON serialized CreateListMsg struct
func DeserializeListMsg(jsonText []byte) (*ListMsg, error) {
	var event ListMsg
	err := json.Unmarshal(jsonText, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (e *ListMsg) ValidateMsg() error {
	err := ValidateMsg(e.Client, e.User, e.Msg, e.Servers)
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

func GetListCQL(msg *ListMsg) (string, []interface{}) {
	cql := `SELECT list_id, category, title, users, is_hidden, created_at, updated_at
			FROM list 
			WHERE list_id = ? LIMIT 1`
	params := []interface{}{msg.Data.ID}
	return cql, params
}

func (d *DAL) GetList(msg *ListMsg) (*List, error) {
	var id gocql.UUID
	var category string
	var title string
	var users []string
	var isHidden bool
	var createdAt time.Time
	var updatedAt time.Time

	cql, params := GetListCQL(msg)
	q := d.session.Query(cql, params...)

	err := q.Scan(&id, &category, &title, &users, &isHidden, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	list := &List{
		ID:        id.String(),
		Category:  category,
		Title:     title,
		Users:     users,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return list, nil
}

func CreateOrUpdateListCQL(msg *ListMsg) (string, []interface{}, error) {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return "", []interface{}{}, fmt.Errorf("DAL.CreateOrUpdateListCQL: %s\n", err.Error())
	}

	cql := `UPDATE list SET category=?, title=?, users=?, is_hidden=?, created_at=?, updated_at=?, msg=?
			 WHERE list_id = ?`

	params := []interface{}{msg.Data.Category, msg.Data.Title, msg.Data.Users, msg.Data.IsHidden,
		msg.Data.CreatedAt, msg.Data.UpdatedAt, msgBytes,
		msg.Data.ID}

	return cql, params, nil
}

func (d *DAL) CreateOrUpdateList(msg *ListMsg) error {
	err := msg.ValidateMsg()
	if err != nil {
		return fmt.Errorf("DAL.CreateOrUpdateList: %s\n", err.Error())
	}

	cql, params, err := CreateOrUpdateListCQL(msg)
	if err != nil {
		return fmt.Errorf("DAL.CreateOrUpdateList: %s\n", err.Error())
	}

	insertList := d.session.Query(cql, params...)
	err = insertList.Exec()
	if err != nil {
		return fmt.Errorf("DAL.CreateOrUpdateList: %s\n", err.Error())
	}

	return nil
}

func (d *DAL) DeleteList(msg *ListMsg) error {
	insertList := d.session.Query(`DELETE FROM list WHERE list_id = ?`, msg.Data.ID)
	err := insertList.Exec()
	if err != nil {
		return err
	}
	return nil
}
