/**
 * Created by GoLand.
 * User: nankai
 * Date: 2019-07-27
 * Time: 14:41
 */

package models

/************************ Request data************************ ***/
// general request data format
type Request struct {
	Seq  string      `json:"seq"`            // Unique Id of the message
	Cmd  string      `json:"cmd"`            // request command word
	Data interface{} `json:"data,omitempty"` // data json
}

// login request data
type Login struct {
	ServiceToken string `json:"serviceToken"` // Verify that the user is logged in
	AppId        uint32 `json:"appId,omitempty"`
	UserId       string `json:"userId,omitempty"`
}

// Heartbeat request data
type HeartBeat struct {
	UserId string `json:"userId,omitempty"`
}
