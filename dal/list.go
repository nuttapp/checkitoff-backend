package dal

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	m "github.com/nuttapp/checkitoff-backend/dal/models"
)

func GetListCQL(msg *m.ListMsg) (string, []interface{}) {
	cql := `SELECT list_id, category, title, users, is_hidden, created_at, updated_at
			FROM list 
			WHERE list_id = ? LIMIT 1`
	params := []interface{}{msg.Data.ID}
	return cql, params
}

func (d *DAL) GetList(msg *m.ListMsg) (*m.List, error) {
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

	list := &m.List{
		ID:        id.String(),
		Category:  category,
		Title:     title,
		Users:     users,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return list, nil
}

func CreateOrUpdateListCQL(msg *m.ListMsg) (string, []interface{}, error) {
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

func (d *DAL) CreateOrUpdateList(msg *m.ListMsg) error {
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

func (d *DAL) DeleteList(msg *m.ListMsg) error {
	insertList := d.session.Query(`DELETE FROM list WHERE list_id = ?`, msg.Data.ID)
	err := insertList.Exec()
	if err != nil {
		return err
	}
	return nil
}
