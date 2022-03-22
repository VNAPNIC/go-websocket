/**
 * Created by GoLand.
 * User: nankai
 * Date: 2019-07-27
 * Time: 13:12
 */

package websocket

import (
	"encoding/json"
	"fmt"
	"gowebsocket/common"
	"gowebsocket/lib/cache"
	"gowebsocket/models"
	"time"

	"github.com/go-redis/redis"
)

// ping
func PingController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {

	code = common.OK
	fmt.Println("webSocket_request ping interface", client.Addr, seq, message)

	data = "pong"

	return
}

// User login
func LoginController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {

	code = common.OK
	currentTime := uint64(time.Now().Unix())

	request := &models.Login{}
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		fmt.Println("User login Failed to parse data", seq, err)

		return
	}

	fmt.Println("webSocket_request user login", seq, "ServiceToken", request.ServiceToken)

	// TODO:: Perform user authorization authentication, usually the client passes in TOKEN, then checks whether the TOKEN is legal, and parses the user ID through TOKEN
	// This project is just a demonstration, so go directly to the user ID passed in by the client
	if request.UserId == "" || len(request.UserId) >= 20 {
		code = common.UnauthorizedUserId
		fmt.Println("user login illegal user", seq, request.UserId)

		return
	}

	if !InAppIds(request.AppId) {
		code = common.Unauthorized
		fmt.Println("User login Unsupported platform", seq, request.AppId)

		return
	}

	if client.IsLogin() {
		fmt.Println("User is logged in User is logged in", client.AppId, client.UserId, seq)
		code = common.OperationFailure

		return
	}

	client.Login(request.AppId, request.UserId, currentTime)

	// Storing data
	userOnline := models.UserLogin(serverIp, serverPort, request.AppId, request.UserId, client.Addr, currentTime)
	err := cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = common.ServerError
		fmt.Println("User login SetUserOnlineInfo", seq, err)

		return
	}

	// User login
	login := &login{
		AppId:  request.AppId,
		UserId: request.UserId,
		Client: client,
	}
	clientManager.Login <- login

	fmt.Println("User login successful", seq, client.Addr, request.UserId)

	return
}

// heartbeat interface
func HeartbeatController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {

	code = common.OK
	currentTime := uint64(time.Now().Unix())

	request := &models.HeartBeat{}
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		fmt.Println("The heartbeat interface failed to parse the data", seq, err)

		return
	}

	fmt.Println("webSocket_request heartbeat interface", client.AppId, client.UserId)

	if !client.IsLogin() {
		fmt.Println("Heartbeat interface user is not logged in", client.AppId, client.UserId, seq)
		code = common.NotLoggedIn

		return
	}

	userOnline, err := cache.GetUserOnlineInfo(client.GetKey())
	if err != nil {
		if err == redis.Nil {
			code = common.NotLoggedIn
			fmt.Println("Heartbeat interface user is not logged in", seq, client.AppId, client.UserId)

			return
		} else {
			code = common.ServerError
			fmt.Println("Heartbeat interface GetUserOnlineInfo", seq, client.AppId, client.UserId, err)

			return
		}
	}

	client.Heartbeat(currentTime)
	userOnline.Heartbeat(currentTime)
	err = cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = common.ServerError
		fmt.Println("Heartbeat interface SetUserOnlineInfo", seq, client.AppId, client.UserId, err)

		return
	}

	return
}
