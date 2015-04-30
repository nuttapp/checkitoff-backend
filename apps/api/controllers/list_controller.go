package controllers

import (
	"encoding/json"
	"errors"
	"fmt"

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

func ListControllerCreate(jsonText []byte, ctx *api.APIContext) error {
	if ctx.Cfg == nil {
		return errors.New("apiCfg cannot be nil")
	}

	msg, err := dal.NewListMsg(dal.MsgMethodCreate, jsonText)

	if err != nil {
		return err
	}

	server := dal.Server{
		Hostname:  ctx.Cfg.Server.Hostname,
		IPAddress: ctx.Cfg.Server.IPAddress,
		Role:      ctx.Cfg.Server.Role,
	}
	msg.Servers = append(msg.Servers, server)

	err = msg.ValidateMsg()
	if err != nil {
		return util.NewError(CreateListMsgValidationError, err)
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return util.NewError(CreateListMsgJSONMarshalError, err)
	}

	replyChan, err := ctx.PublishMsg(msg.ID, b)
	if err != nil {
		return util.NewError(ProducerPublishError, err)
	}
	foo := <-replyChan
	fmt.Println(foo)
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
