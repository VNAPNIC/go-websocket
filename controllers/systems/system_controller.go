/**
* Created by GoLand.
* User: nankai
* Date: 2019-07-25
* Time: 12:11
 */

package systems

import (
	"fmt"
	"gowebsocket/common"
	"gowebsocket/controllers"
	"gowebsocket/servers/websocket"
	"runtime"

	"github.com/gin-gonic/gin"
)

// Query system status
func Status(c *gin.Context) {

	isDebug := c.Query("isDebug")
	fmt.Println("http_request query system status", isDebug)

	data := make(map[string]interface{})

	numGoroutine := runtime.NumGoroutine()
	numCPU := runtime.NumCPU()

	// number of goroutines
	data["numGoroutine"] = numGoroutine
	data["numCPU"] = numCPU

	// ClientManager information
	data["managerInfo"] = websocket.GetManagerInfo(isDebug)

	controllers.Response(c, common.OK, "", data)
}
