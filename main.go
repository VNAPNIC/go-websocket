/**
* Created by GoLand.
* User: nankai
* Date: 2019-07-25
* Time: 09:59
 */

package main

import (
	"fmt"
	"go-websocket/lib/redislib"
	"go-websocket/routers"
	"go-websocket/servers/grpcserver"
	"go-websocket/servers/task"
	"go-websocket/servers/websocket"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	initConfig()

	initFile()

	initRedis()

	router := gin.Default()
	// Initialize route
	routers.Init(router)
	routers.WebsocketInit()

	// timed task
	task.Init()

	// service registration
	task.ServerInit()

	go websocket.StartWebSocket()
	// grpc
	go grpcserver.Init()

	go open()

	httpPort := viper.GetString("app.httpPort")
	http.ListenAndServe(":"+httpPort, router)

}

// Initialize log
func initFile() {
	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	// Logging to a file.
	logFile := viper.GetString("app.logFile")
	f, _ := os.Create(logFile)
	gin.DefaultWriter = io.MultiWriter(f)
}

func initConfig() {
	viper.SetConfigName("config/app")
	viper.AddConfigPath(".") // Add search path

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config redis:", viper.Get("redis"))

}

func initRedis() {
	redislib.ExampleNewClient()
}

func open() {

	time.Sleep(1000 * time.Millisecond)

	httpUrl := viper.GetString("app.httpUrl")
	httpUrl = "http://" + httpUrl + "/home/index"

	fmt.Println("Visit page experience:", httpUrl)

	cmd := exec.Command("open", httpUrl)
	cmd.Output()
}
