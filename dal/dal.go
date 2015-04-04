package dal

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gocql/gocql"
	m "github.com/nuttapp/checkitoff-backend/common/models"
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

	if d == nil {
		cluster := gocql.NewCluster(d.Cfg.Hosts...)
		cluster.Keyspace = d.Cfg.Keyspace
		cluster.Consistency = d.Cfg.Consistency

		session, err := cluster.CreateSession()
		if err != nil {
			return err
		}
		d.session = session
	}
	return nil
}

func (d *DAL) SaveCreateListMsg(msg *m.ListMsg) error {
	err := d.initSession()
	if err != nil {
		return err
	}
	createdAt := gocql.TimeUUID()
	updatedAt := gocql.TimeUUID()

	insertList := d.session.Query(
		`INSERT INTO list (list_id, title, created_at, updated_at, users) VALUES (?, ?, ?, ?, ?)`,
		msg.Data.ID, msg.Data.Title, createdAt, updatedAt, []string{msg.User.ID})
	err = insertList.Exec()
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
