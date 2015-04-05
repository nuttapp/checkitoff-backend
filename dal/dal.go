package dal

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gocql/gocql"
	m "github.com/nuttapp/checkitoff-backend/dal/models"
)

type DALCfg struct {
	Hosts       []string
	Keyspace    string
	Consistency gocql.Consistency
}

type DAL struct {
	Cfg     DALCfg
	session *gocql.Session
}

func NewDAL(cfg DALCfg) (*DAL, error) {
	d := &DAL{
		Cfg: cfg,
	}
	err := d.initSession()
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *DAL) initSession() error {
	if len(d.Cfg.Hosts) == 0 {
		return errors.New("No hosts were provided for cassandra cluster")
	}
	if len(d.Cfg.Keyspace) == 0 {
		return errors.New("No keyspace was provided")
	}
	if d.session == nil {
		cluster := gocql.NewCluster(d.Cfg.Hosts...)
		cluster.Keyspace = d.Cfg.Keyspace
		cluster.Consistency = d.Cfg.Consistency

		session, err := cluster.CreateSession()
		// fmt.Println(session)
		if err != nil {
			return err
		}
		// fmt.Printf("session: %s\n", session)
		d.session = session
	}
	return nil
}

func (d *DAL) HandleListMsg(msg *m.ListMsg) error {
	err := d.initSession()
	if err != nil {
		return err
	}

	switch msg.Method {
	case m.MsgMethodCreate:
		return d.CreateList(msg)

	}
	return nil
}

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

func (d *DAL) CreateList(msg *m.ListMsg) error {
	insertList := d.session.Query(
		`INSERT INTO list (list_id, title, created_at, updated_at, users) VALUES (?, ?, ?, ?, ?)`,
		msg.Data.ID, msg.Data.Title, msg.Data.CreatedAt, msg.Data.UpdatedAt, []string{msg.User.ID})
	err := insertList.Exec()
	if err != nil {
		return err
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	msgType := fmt.Sprintf("%s-%s", msg.Method, msg.Resource)
	insertListMsg := d.session.Query(
		`INSERT INTO list_event (list_id, user_id, event_id, event_type, data) VALUES (?, ?, ?, ?, ?)`,
		msg.Data.ID, msg.User.ID, msg.ID, msgType, b)
	err = insertListMsg.Exec()
	if err != nil {
		return err
	}

	insertUserTimeline := d.session.Query(
		`INSERT INTO user_timeline (user_id, event_id, event_type, data) VALUES (?, ?, ?, ?)`,
		msg.User.ID, msg.ID, msgType, b)
	err = insertUserTimeline.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (d *DAL) UpdateList(msg *m.ListMsg) error {
	if len(msg.ID) == 0 {
		return fmt.Errorf("DAL.UpdateList: %s", m.MissingListIDError)
	}
	return nil
}
