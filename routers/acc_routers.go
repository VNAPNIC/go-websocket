/**
 * Created by GoLand.
 * User: nankai
 * Date: 2019-07-25
 * Time: 16:02
 */

package routers

import (
	"go-websocket/servers/websocket"
)

// Websocket routing
func WebsocketInit() {
	websocket.Register("login", websocket.LoginController)
	websocket.Register("heartbeat", websocket.HeartbeatController)
	websocket.Register("ping", websocket.PingController)
}
