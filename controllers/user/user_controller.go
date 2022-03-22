/**
* Created by GoLand.
* User: nankai
* Date: 2019-07-25
* Time: 12:11
 */

package user

import (
	"fmt"
	"gowebsocket/common"
	"gowebsocket/controllers"
	"gowebsocket/lib/cache"
	"gowebsocket/models"
	"gowebsocket/servers/websocket"
	"strconv"

	"github.com/gin-gonic/gin"
)

// View all online users
func List(c *gin.Context) {

	appIdStr := c.Query("appId")
	appIdUint64, _ := strconv.ParseInt(appIdStr, 10, 32)
	appId := uint32(appIdUint64)

	fmt.Println("http_request view all online users", appId)

	data := make(map[string]interface{})

	userList := websocket.UserList(appId)
	data["userList"] = userList
	data["userCount"] = len(userList)

	controllers.Response(c, common.OK, "", data)
}

// Check if the user is online
func Online(c *gin.Context) {

	userId := c.Query("userId")
	appIdStr := c.Query("appId")
	appIdUint64, _ := strconv.ParseInt(appIdStr, 10, 32)
	appId := uint32(appIdUint64)

	fmt.Println("http_request see if the user is online", userId, appIdStr)

	data := make(map[string]interface{})

	online := websocket.CheckUserOnline(appId, userId)
	data["userId"] = userId
	data["online"] = online

	controllers.Response(c, common.OK, "", data)
}

// send message to user
func SendMessage(c *gin.Context) {
	// get parameters
	appIdStr := c.PostForm("appId")
	userId := c.PostForm("userId")
	msgId := c.PostForm("msgId")
	message := c.PostForm("message")
	appIdUint64, _ := strconv.ParseInt(appIdStr, 10, 32)
	appId := uint32(appIdUint64)

	fmt.Println("http_request send message to user", appIdStr, userId, msgId, message)

	// TODO:: Perform user authorization authentication, usually the client passes in TOKEN, then checks whether the TOKEN is legal, and parses the user ID through TOKEN
	// This project is just a demonstration, so go directly to the user ID (userId) passed in by the client

	data := make(map[string]interface{})

	if cache.SeqDuplicates(msgId) {
		fmt.Println("Send message to user repeat submission:", msgId)
		controllers.Response(c, common.OK, "", data)

		return
	}

	sendResults, err := websocket.SendUserMessage(appId, userId, msgId, message)
	if err != nil {
		data["sendResultsErr"] = err.Error()
	}

	data["sendResults"] = sendResults

	controllers.Response(c, common.OK, "", data)
}

// Send a message to everyone
func SendMessageAll(c *gin.Context) {
	// get parameters
	appIdStr := c.PostForm("appId")
	userId := c.PostForm("userId")
	msgId := c.PostForm("msgId")
	message := c.PostForm("message")
	appIdUint64, _ := strconv.ParseInt(appIdStr, 10, 32)
	appId := uint32(appIdUint64)

	fmt.Println("http_request sends a message to all users", appIdStr, userId, msgId, message)

	data := make(map[string]interface{})
	if cache.SeqDuplicates(msgId) {
		fmt.Println("Send message to user repeat submission:", msgId)
		controllers.Response(c, common.OK, "", data)

		return
	}

	sendResults, err := websocket.SendUserMessageAll(appId, userId, msgId, models.MessageCmdMsg, message)
	if err != nil {
		data["sendResultsErr"] = err.Error()

	}

	data["sendResults"] = sendResults

	controllers.Response(c, common.OK, "", data)

}
