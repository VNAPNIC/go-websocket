/**
 * Created by GoLand.
 * User: nankai
 * Date: 2019-07-27
 * Time: 14:38
 */

package websocket

import (
	"encoding/json"
	"fmt"
	"go-websocket/common"
	"go-websocket/models"
	"sync"
)

type DisposeFunc func(client *Client, seq string, message []byte) (code uint32, msg string, data interface{})

var (
	handlers        = make(map[string]DisposeFunc)
	handlersRWMutex sync.RWMutex
)

// register
func Register(key string, value DisposeFunc) {
	handlersRWMutex.Lock()
	defer handlersRWMutex.Unlock()
	handlers[key] = value

	return
}

func getHandlers(key string) (value DisposeFunc, ok bool) {
	handlersRWMutex.RLock()
	defer handlersRWMutex.RUnlock()

	value, ok = handlers[key]

	return
}

// Data processing
func ProcessData(client *Client, message []byte) {

	fmt.Println("Processing data", client.Addr, string(message))

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Processing data stop", r)
		}
	}()

	request := &models.Request{}

	err := json.Unmarshal(message, request)
	if err != nil {
		fmt.Println("Processing data json Unmarshal", err)
		client.SendMsg([]byte("Invalid data"))

		return
	}

	requestData, err := json.Marshal(request.Data)
	if err != nil {
		fmt.Println("Processing data json Marshal", err)
		client.SendMsg([]byte("Failed to process data"))

		return
	}

	seq := request.Seq
	cmd := request.Cmd

	var (
		code uint32
		msg  string
		data interface{}
	)

	// request
	fmt.Println("acc_request", cmd, client.Addr)

	// Register with map
	if value, ok := getHandlers(cmd); ok {
		code, msg, data = value(client, seq, requestData)
	} else {
		code = common.RoutingNotExist
		fmt.Println("Processing data routing does not exist", client.Addr, "cmd", cmd)
	}

	msg = common.GetErrorMessage(code, msg)

	responseHead := models.NewResponseHead(seq, cmd, code, msg, data)

	headByte, err := json.Marshal(responseHead)
	if err != nil {
		fmt.Println("Processing data json Marshal", err)

		return
	}

	client.SendMsg(headByte)

	fmt.Println("acc_response send", client.Addr, client.AppId, client.UserId, "cmd", cmd, "code", code)

	return
}
