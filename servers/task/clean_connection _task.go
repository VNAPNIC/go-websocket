/**
* Created by GoLand.
* User: nankai
* Date: 2019-07-31
* Time: 15:17
 */

package task

import (
	"fmt"
	"gowebsocket/servers/websocket"
	"runtime/debug"
	"time"
)

func Init() {
	Timer(3*time.Second, 30*time.Second, cleanConnection, "", nil, nil)

}

// Clean up timed out connections
func cleanConnection(param interface{}) (result bool) {
	result = true

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ClearTimeoutConnections stop", r, string(debug.Stack()))
		}
	}()

	fmt.Println("Scheduled tasks to clean up timeout connections", param)

	websocket.ClearTimeoutConnections()

	return
}
