/**
* Created by GoLand.
* User: nankai
* Date: 2019-08-03
* Time: 15:44
 */

package task

import (
	"fmt"
	"gowebsocket/lib/cache"
	"gowebsocket/servers/websocket"
	"runtime/debug"
	"time"
)

func ServerInit() {
	Timer(2*time.Second, 60*time.Second, server, "", serverDefer, "")
}

// service registration
func server(param interface{}) (result bool) {
	result = true

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("service registration stop", r, string(debug.Stack()))
		}
	}()

	server := websocket.GetServer()
	currentTime := uint64(time.Now().Unix())
	fmt.Println("timed task，service registration", param, server, currentTime)

	cache.SetServerInfo(server, currentTime)

	return
}

// service offline
func serverDefer(param interface{}) (result bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("service offline stop", r, string(debug.Stack()))
		}
	}()

	fmt.Println("service offline", param)

	server := websocket.GetServer()
	cache.DelServerInfo(server)

	return
}
