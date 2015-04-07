package dal

import (
	"errors"

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
	case m.MsgMethodCreate, m.MsgMethodUpdate:
		return d.CreateOrUpdateList(msg)

	}
	return nil
}
