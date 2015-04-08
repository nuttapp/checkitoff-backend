package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"github.com/bitly/go-nsq"
	"github.com/gocql/gocql"
	"github.com/nuttapp/checkitoff-backend/common"
	"github.com/nuttapp/checkitoff-backend/dal"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	conCfg := &common.ConsumerConfig{
		Topic:           "api_events",
		Channel:         "persistence",
		LookupdHTTPaddr: "127.0.0.1:4161",
		Concurrency:     10,
	}

	nsqCfg := nsq.NewConfig()
	nsqCfg.MaxInFlight = 10

	handler := &DALMessageHandler{}
	mh := common.NewMessageConsumer(conCfg, nsqCfg, handler)
	mh.Start()

	for {
		select {
		// On quit signal stop the DAL consumer
		case <-signalChan:
			fmt.Println("")
			mh.Stop()
		case err := <-mh.StopChan:
			if err != nil {
				fmt.Printf("APIConsumer stopped with error: %s", err)
			}
			fmt.Println("END")
			return
		}
	}
}

type DALMessageHandler struct{}

func (mh *DALMessageHandler) HandleMessage(msg *nsq.Message) error {
	var j map[string]interface{}
	err := json.Unmarshal(msg.Body, &j)
	if err != nil {
		fmt.Println("Failed to unmarshal JSON from NSQ message: %s", err)
	}

	cmd := fmt.Sprintf("%s-%s", j["method"], j["resource"])
	fmt.Printf("RECV msg: %s, %s \n", cmd, j["id"])

	switch cmd {
	case "create-list":
		cle, err := dal.DeserializeListMsg(msg.Body)
		if err != nil {
			return err
		}

		err = SaveCreateListMsg(cle)
		if err != nil {
			return err
		}
	}

	// Get list of users to notify
	// Send notification to each of the users
	return nil
}

func SaveInviteUserToListMsg() {
	// addToSetQuery := session.Query(`UPDATE list SET users = users + {'blah2'} WHERE list_id = ?`, cle.EventData.ID)
	// err = addToSetQuery.Exec()
	// if err != nil {
	// 	return err
	// }
}

func SaveCreateListMsg(cle *dal.ListMsg) error {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "demodb"
	cluster.Consistency = gocql.Quorum
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
