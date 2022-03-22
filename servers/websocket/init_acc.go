/**
 * Created by GoLand.
 * User: nankai
 * Date: 2019-07-25
 * Time: 16:04
 */

package websocket

import (
	"fmt"
	"net/http"
	"time"

	"gowebsocket/helper"
	"gowebsocket/models"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

const (
	defaultAppId = 101 // Default Platform Id
)

var (
	clientManager = NewClientManager()                    // manager
	appIds        = []uint32{defaultAppId, 102, 103, 104} // all platforms

	serverIp   string
	serverPort string
)

func GetAppIds() []uint32 {

	return appIds
}

func GetServer() (server *models.Server) {
	server = models.NewServer(serverIp, serverPort)

	return
}

func IsLocal(server *models.Server) (isLocal bool) {
	if server.Ip == serverIp && server.Port == serverPort {
		isLocal = true
	}

	return
}

func InAppIds(appId uint32) (inAppId bool) {

	for _, value := range appIds {
		if value == appId {
			inAppId = true

			return
		}
	}

	return
}

func GetDefaultAppId() (appId uint32) {
	appId = defaultAppId

	return
}

// starting program
func StartWebSocket() {

	serverIp = helper.GetServerIp()

	webSocketPort := viper.GetString("app.webSocketPort")
	rpcPort := viper.GetString("app.rpcPort")

	serverPort = rpcPort

	http.HandleFunc("/acc", wsPage)

	// add handler
	go clientManager.start()
	fmt.Println("WebSocket startup program succeeded", serverIp, serverPort)

	http.ListenAndServe(":"+webSocketPort, nil)
}

func wsPage(w http.ResponseWriter, req *http.Request) {

	// upgrade protocol
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		fmt.Println("Upgrade Protocol", "ua:", r.Header["User-Agent"], "referer:", r.Header["Referer"])

		return true
	}}).Upgrade(w, req, nil)
	if err != nil {
		http.NotFound(w, req)

		return
	}

	fmt.Println("webSocket established connection:", conn.RemoteAddr().String())

	currentTime := uint64(time.Now().Unix())
	client := NewClient(conn.RemoteAddr().String(), conn, currentTime)

	go client.read()
	go client.write()

	// user connect event
	clientManager.Register <- client
}
