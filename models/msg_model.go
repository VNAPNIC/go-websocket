/**
* Created by GoLand.
* User: nankai
* Date: 2019-08-01
* Time: 10:40
 */

package models

import "go-websocket/common"

const (
	MessageTypeText = "text"

	MessageCmdMsg   = "msg"
	MessageCmdEnter = "enter"
	MessageCmdExit  = "exit"
)

// Definition of message
type Message struct {
	Target string `json:"target"` // Target
	Type   string `json:"type"`   // message type text/img/
	Msg    string `json:"msg"`    // message content
	From   string `json:"from"`   // sender
}

func NewTestMsg(from string, Msg string) (message *Message) {

	message = &Message{
		Type: MessageTypeText,
		From: from,
		Msg:  Msg,
	}

	return
}

func getTextMsgData(cmd, uuId, msgId, message string) string {
	textMsg := NewTestMsg(uuId, message)
	head := NewResponseHead(msgId, cmd, common.OK, "Ok", textMsg)

	return head.String()
}

// text message
func GetMsgData(uuId, msgId, cmd, message string) string {

	return getTextMsgData(cmd, uuId, msgId, message)
}

// text message
func GetTextMsgData(uuId, msgId, message string) string {

	return getTextMsgData("msg", uuId, msgId, message)
}

// User entry message
func GetTextMsgDataEnter(uuId, msgId, message string) string {

	return getTextMsgData("enter", uuId, msgId, message)
}

// User entry message
func GetTextMsgDataExit(uuId, msgId, message string) string {

	return getTextMsgData("exit", uuId, msgId, message)
}
