/**
* Created by GoLand.
* User: nankai
* Date: 2019-08-01
* Time: 10:46
 */

package models

import "encoding/json"

/************************ Response data************************ ***/
type Head struct {
	Seq      string    `json:"seq"`      // Id of the message
	Cmd      string    `json:"cmd"`      // cmd action of the message
	Response *Response `json:"response"` // message body
}

type Response struct {
	Code    uint32      `json:"code"`
	CodeMsg string      `json:"codeMsg"`
	Data    interface{} `json:"data"` // data json
}

// push data structure
type PushMsg struct {
	Seq  string `json:"seq"`
	Uuid uint64 `json:"uuid"`
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

// set return message
func NewResponseHead(seq string, cmd string, code uint32, codeMsg string, data interface{}) *Head {
	response := NewResponse(code, codeMsg, data)

	return &Head{Seq: seq, Cmd: cmd, Response: response}
}

func (h *Head) String() (headStr string) {
	headBytes, _ := json.Marshal(h)
	headStr = string(headBytes)

	return
}

func NewResponse(code uint32, codeMsg string, data interface{}) *Response {
	return &Response{Code: code, CodeMsg: codeMsg, Data: data}
}
