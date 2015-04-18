package controllers

import (
	"encoding/json"
	"errors"

	"github.com/bitly/go-nsq"
	"github.com/nuttapp/checkitoff-backend/apps/api"
	"github.com/nuttapp/checkitoff-backend/common/util"
	"github.com/nuttapp/checkitoff-backend/dal"
)

const (
	ProducerConnectionError = "Failed to connect to NSQ producer"
	ProducerPublishError    = "Failed to publish to NSQ producer"
)

const (
	CreateListMsgJSONMarshalError = "Failed to marshal CreateListMsg into json"
	CreateListMsgValidationError  = "Validation failed for CreateListMsg"
)

func ListControllerCreate(jsonText []byte, c *api.APIContext) error {
	apiCfg := c.APICfg
	nsqCfg := c.NSQCfg

	if c.APICfg == nil {
		return errors.New("apiCfg cannot be nil")
	}

	msg, err := dal.NewListMsg(dal.MsgMethodCreate, jsonText)
	if err != nil {
		return err
	}

	server := dal.Server{
		Hostname:  apiCfg.Hostname,
		IPAddress: apiCfg.IPAddress,
		Role:      apiCfg.Role,
	}
	msg.Servers = append(msg.Servers, server)

	err = msg.ValidateMsg()
	if err != nil {
		return util.NewError(CreateListMsgValidationError, err)
	}

	producer, err := nsq.NewProducer(apiCfg.NSQProducerTCPAddr, nsqCfg)
	if err != nil {
		return util.NewError(ProducerConnectionError, err)
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return util.NewError(CreateListMsgJSONMarshalError, err)
	}

	err = producer.Publish(apiCfg.NSQPubTopic, b)
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
