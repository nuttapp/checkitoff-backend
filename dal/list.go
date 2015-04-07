package dal

import (
	"github.com/gocql/gocql"
	m "github.com/nuttapp/checkitoff-backend/dal/models"
)

func (d *DAL) GetList(msg *m.ListMsg) (*m.List, error) {
	var id gocql.UUID
	var category string
	var title string
	var users []string
	var isHidden bool
	var createdAt gocql.UUID
	var updatedAt gocql.UUID

	q := d.session.Query(
		`SELECT list_id, category, title, users, is_hidden, created_at, updated_at
		 FROM list 
		 WHERE list_id = ? LIMIT 1`, msg.Data.ID)

	err := q.Scan(&id, &category, &title, &users, &isHidden, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	list := &m.List{
		ID:        id.String(),
		Category:  category,
		Title:     title,
		Users:     users,
		CreatedAt: createdAt.Time(),
		UpdatedAt: updatedAt.Time(),
	}

	return list, nil
}

func (d *DAL) CreateOrUpdateList(msg *m.ListMsg) error {
	insertList := d.session.Query(
		`UPDATE list SET category = ?, title = ?, users = ?, is_hidden = ?, created_at = ?, updated_at = ?
		 WHERE list_id = ?`,
		msg.Data.Category, msg.Data.Title, msg.Data.Users, msg.Data.IsHidden,
		gocql.UUIDFromTime(msg.Data.CreatedAt), gocql.UUIDFromTime(msg.Data.UpdatedAt), msg.Data.ID)
	err := insertList.Exec()
	if err != nil {
		return err
	}
	// `INSERT INTO list
	// 	(list_id, category, title, users, is_hidden, created_at, updated_at)
	// VALUES
	// 	(?, ?, ?, ?, ?, ?, ?)`,

	// b, err := json.Marshal(msg)
	// if err != nil {
	// 	return err
	// }

	// msgType := fmt.Sprintf("%s-%s", msg.Method, msg.Resource)
	// insertListMsg := d.session.Query(
	// 	`INSERT INTO list_event (list_id, user_id, event_id, event_type, data) VALUES (?, ?, ?, ?, ?)`,
	// 	msg.Data.ID, msg.User.ID, msg.ID, msgType, b)
	// err = insertListMsg.Exec()
	// if err != nil {
	// 	return err
	// }
	//
	// insertUserTimeline := d.session.Query(
	// 	`INSERT INTO user_timeline (user_id, event_id, event_type, data) VALUES (?, ?, ?, ?)`,
	// 	msg.User.ID, msg.ID, msgType, b)
	// err = insertUserTimeline.Exec()
	// if err != nil {
	// 	return err
	// }

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
