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
	Cfg DALCfg
}

func (d *DAL) SaveCreateListMsg(cle *m.ListMsg) error {
	if len(d.Cfg.Hosts) == 0 {
		return errors.New("No hosts were provided for cassandra cluster")
	}
	if len(d.Cfg.Keyspace) == 0 {
		return errors.New("No keyspace was provided")
	}
	cluster := gocql.NewCluster(d.Cfg.Hosts...)
	cluster.Keyspace = d.Cfg.Keyspace
	cluster.Consistency = d.Cfg.Consistency
	session, _ := cluster.CreateSession()
	defer session.Close()

	createdAt := gocql.TimeUUID()
	updatedAt := gocql.TimeUUID()

	insertList := session.Query(`INSERT INTO list (list_id, title, created_at, updated_at, users) VALUES (?, ?, ?, ?, ?)`,
		cle.Data.ID, cle.Data.Title, createdAt, updatedAt, []string{cle.User.ID})
	err := insertList.Exec()
	if err != nil {
		return err
	}

	b, err := json.Marshal(cle)
	if err != nil {
		return err
	}

	msgType := fmt.Sprintf("%s-%s", cle.Method, cle.Resource)
	insertListMsg := session.Query(`INSERT INTO list_event (list_id, user_id, event_id, event_type, data) VALUES (?, ?, ?, ?, ?)`,
		cle.Data.ID, cle.User.ID, cle.ID, msgType, b)
	err = insertListMsg.Exec()
	if err != nil {
		return err
	}

	insertUserTimeline := session.Query(`INSERT INTO user_timeline (user_id, event_id, event_type, data) VALUES (?, ?, ?, ?)`,
		cle.User.ID, cle.ID, msgType, b)
	err = insertUserTimeline.Exec()
	if err != nil {
		return err
	}

	return nil
}
