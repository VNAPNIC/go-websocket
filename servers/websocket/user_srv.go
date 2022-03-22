/**
* Created by GoLand.
* User: nankai
* Date: 2019-07-30
* Time: 12:27
 */

package websocket

import (
	"errors"
	"fmt"
	"gowebsocket/lib/cache"
	"gowebsocket/models"
	"gowebsocket/servers/grpcclient"
	"time"

	"github.com/go-redis/redis"
)

// Query all users
func UserList(appId uint32) (userList []string) {

	userList = make([]string, 0)
	currentTime := uint64(time.Now().Unix())
	servers, err := cache.GetServerAll(currentTime)
	if err != nil {
		fmt.Println("Send a message to all users", err)

		return
	}

	for _, server := range servers {
		var (
			list []string
		)
		if IsLocal(server) {
			list = GetUserList(appId)
		} else {
			list, _ = grpcclient.GetUserList(server, appId)
		}
		userList = append(userList, list...)
	}

	return
}

// Check if the user is online
func CheckUserOnline(appId uint32, userId string) (online bool) {
	// full platform query
	if appId == 0 {
		for _, appId := range GetAppIds() {
			online, _ = checkUserOnline(appId, userId)
			if online == true {
				break
			}
		}
	} else {
		online, _ = checkUserOnline(appId, userId)
	}

	return
}

// Check if the user is online
func checkUserOnline(appId uint32, userId string) (online bool, err error) {
	key := GetUserKey(appId, userId)
	userOnline, err := cache.GetUserOnlineInfo(key)
	if err != nil {
		if err == redis.Nil {
			fmt.Println("GetUserOnlineInfo", appId, userId, err)

			return false, nil
		}

		fmt.Println("GetUserOnlineInfo", appId, userId, err)

		return
	}

	online = userOnline.IsOnline()

	return
}

// send message to user
func SendUserMessage(appId uint32, userId string, msgId, message string) (sendResults bool, err error) {

	data := models.GetTextMsgData(userId, msgId, message)

	client := GetUserClient(appId, userId)

	if client != nil {
		// send locally
		sendResults, err = SendUserMessageLocal(appId, userId, data)
		if err != nil {
			fmt.Println("Send message to user", appId, userId, err)
		}

		return
	}

	key := GetUserKey(appId, userId)
	info, err := cache.GetUserOnlineInfo(key)
	if err != nil {
		fmt.Println("Failed to send message to user", key, err)

		return false, nil
	}
	if !info.IsOnline() {
		fmt.Println("User is not online", key)
		return false, nil
	}
	server := models.NewServer(info.AccIp, info.AccPort)
	msg, err := grpcclient.SendMsg(server, msgId, appId, userId, models.MessageCmdMsg, models.MessageCmdMsg, message)
	if err != nil {
		fmt.Println("Failed to send message to user", key, err)

		return false, err
	}
	fmt.Println("Send message to user successfully -rpc", msg)
	sendResults = true

	return
}

// Send a message to the local user
func SendUserMessageLocal(appId uint32, userId string, data string) (sendResults bool, err error) {

	client := GetUserClient(appId, userId)
	if client == nil {
		err = errors.New("User is not online")

		return
	}

	// send messages
	client.SendMsg([]byte(data))
	sendResults = true

	return
}

// Send a message to all users
func SendUserMessageAll(appId uint32, userId string, msgId, cmd, message string) (sendResults bool, err error) {
	sendResults = true

	currentTime := uint64(time.Now().Unix())
	servers, err := cache.GetServerAll(currentTime)
	if err != nil {
		fmt.Println("Send a message to all users", err)

		return
	}

	for _, server := range servers {
		if IsLocal(server) {
			data := models.GetMsgData(userId, msgId, cmd, message)
			AllSendMessages(appId, userId, data)
		} else {
			grpcclient.SendMsgAll(server, msgId, appId, userId, cmd, message)
		}
	}

	return
}
