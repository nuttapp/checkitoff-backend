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

type ListEvent struct {
	Event
	Servers []Server `json:"servers"`
	Client  Client   `json:"client"`
	User    User     `json:"user"`
	Data    List     `json:"data"`
	Err     error    `json:"error"`
}

// NewListEvent constructs ListEvent structs
// It takes the json sent by a client, and creates a valid ListEvent struct,
// which can be sent to NSQ and consumed by our apps (persistor, logger etc..)
// You should not try to construct a ListEvent manually
func NewListEvent(eMethod string, eJSON []byte) (*ListEvent, error) {
	e, err := DeserializeListEvent(eJSON)
	if err != nil {
		return nil, err
	}

	eID, err := gocql.RandomUUID()
	if err != nil {
		return nil, err
	}

	switch eMethod {
	case MsgMethodCreate:
		e.ID = eID.String()

		listID, err := gocql.RandomUUID()
		if err != nil {
			return nil, err
		}
		e.Data.ID = listID.String()

		e.Data.CreatedAt = time.Now().UTC()
		e.Data.UpdatedAt = e.Data.CreatedAt

		// add the user creating the list to users of the list
		e.Data.Users = append(e.Data.Users, e.User.ID)
	case MsgMethodUpdate:
		e.ID = eID.String()
		e.Data.UpdatedAt = time.Now().UTC()
	}

	return e, nil
}

// DeserializeListEvent deserializes a ListEvent formatted as JSON
func DeserializeListEvent(eJSON []byte) (*ListEvent, error) {
	var event ListEvent
	err := json.Unmarshal(eJSON, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (e *ListEvent) Validate() error {
	err := ValidateEvent(e.Client, e.User, e.Event, e.Servers)
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
		return errors.New(InvalidEventMethodError)
	}
	if e.Resource != MsgResourceList {
		return errors.New(InvalidEventResourceError)
	}
	return nil
}

func GetListCQL(e *ListEvent) (string, []interface{}) {
	cql := `SELECT list_id, category, title, users, is_hidden, created_at, updated_at
			FROM list 
			WHERE list_id = ? LIMIT 1`
	params := []interface{}{e.Data.ID}
	return cql, params
}

func (d *DAL) GetList(e *ListEvent) (*List, error) {
	var id gocql.UUID
	var category string
	var title string
	var users []string
	var isHidden bool
	var createdAt time.Time
	var updatedAt time.Time

	cql, params := GetListCQL(e)
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

func CreateOrUpdateListCQL(e *ListEvent) (string, []interface{}, error) {
	eBytes, err := json.Marshal(e)
	if err != nil {
		return "", []interface{}{}, fmt.Errorf("DAL.CreateOrUpdateListCQL: %s\n", err.Error())
	}

	cql := `UPDATE list SET category=?, title=?, users=?, is_hidden=?, created_at=?, updated_at=?, msg=?
			 WHERE list_id = ?`

	params := []interface{}{e.Data.Category, e.Data.Title, e.Data.Users, e.Data.IsHidden,
		e.Data.CreatedAt, e.Data.UpdatedAt, eBytes,
		e.Data.ID}

	return cql, params, nil
}

func (d *DAL) CreateOrUpdateList(e *ListEvent) error {
	err := e.Validate()
	if err != nil {
		return fmt.Errorf("DAL.CreateOrUpdateList: %s\n", err.Error())
	}

	cql, params, err := CreateOrUpdateListCQL(e)
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

func (d *DAL) DeleteList(e *ListEvent) error {
	insertList := d.session.Query(`DELETE FROM list WHERE list_id = ?`, e.Data.ID)
	err := insertList.Exec()
	if err != nil {
		return err
	}
	return nil
}
