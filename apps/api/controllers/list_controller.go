package controllers

import (
	"encoding/json"
	"errors"

	"github.com/bitly/go-nsq"
	"github.com/gocql/gocql"
	"github.com/nuttapp/checkitoff-backend/api/config"
	"github.com/nuttapp/checkitoff-backend/common/util"
	m "github.com/nuttapp/checkitoff-backend/dal/models"
)

const (
	ProducerConnectionError = "Failed to connect to NSQ producer"
	ProducerPublishError    = "Failed to publish to NSQ producer"
)

const (
	CreateListMsgJSONUnmarshalError = "Failed to unmarshal json into CreateListMsg"
	CreateListMsgJSONMarshalError   = "Failed to marshal CreateListMsg into json"
	CreateListMsgValidationError    = "Validation failed for CreateListMsg"
)

func ListControllerCreate(jsonText []byte, nsqCfg *nsq.Config, apiCfg *config.Config) error {
	var event m.ListMsg
	err := json.Unmarshal(jsonText, &event)
	if err != nil {
		return util.NewError(CreateListMsgJSONUnmarshalError, err)
	}

	event.ID = gocql.TimeUUID().String()

	listID, err := gocql.RandomUUID()
	if err != nil {
		return errors.New("Failed to create UUID")
	}

	event.Data.ID = listID.String()

	if apiCfg == nil {
		return errors.New("apiCfg cannot be nil")
	}

	event.Server.Hostname = apiCfg.Hostname
	event.Server.IPAddress = apiCfg.IPAddress

	err = event.IsReadyToBeSaved()
	if err != nil {
		return util.NewError(CreateListMsgValidationError, err)
	}

	producer, err := nsq.NewProducer(apiCfg.NSQProducerTCPAddr, nsqCfg)
	if err != nil {
		return util.NewError(ProducerConnectionError, err)
	}

	b, err := json.Marshal(event)
	if err != nil {
		return util.NewError(CreateListMsgJSONMarshalError, err)
	}

	err = producer.Publish("api_events", b)
	if err != nil {
		return util.NewError(ProducerPublishError, err)
	}

	// - validate json
	// ? authenticate the user
	// - Create create-list-event on NSQ
	// - enqueue create-list-event on nsqd topic api_events
	// wait for create-list-response on topic create-response
	// Ask if the item on response-create something that we enqueued
	return nil
}

// hj, ok := w.(http.Hijacker)
// if !ok {
// 	http.Error(w, "httpserver doesn't support hijacking", http.StatusInternalServerError)
// 	fmt.Println("httpserver doesn't support hijacking")
// 	return
// }
// conn, bufrw, err := hj.Hijack()
// if err != nil {
// 	fmt.Println(err)
// 	http.Error(w, err.Error(), http.StatusInternalServerError)
// 	return
// }
// _ = conn
// _ = bufrw
