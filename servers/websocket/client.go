/**
 * Created by GoLand.
 * User: nankai
 * Date: 2019-07-25
 * Time: 16:24
 */

package websocket

import (
	"fmt"
	"runtime/debug"

	"github.com/gorilla/websocket"
)

const (
	// User connection timeout
	heartbeatExpirationTime = 6 * 60
)

// User login
type login struct {
	AppId  uint32
	UserId string
	Client *Client
}

// read client data
func (l *login) GetKey() (key string) {
	key = GetUserKey(l.AppId, l.UserId)

	return
}

// user connection
type Client struct {
	Addr          string          // client address
	Socket        *websocket.Conn // user connection
	Send          chan []byte     // data to be sent
	AppId         uint32          // Platform ID for login app/web/ios
	UserId        string          // User Id, only available after the user logs in
	FirstTime     uint64          // first connect event
	HeartbeatTime uint64          // User's last heartbeat time
	LoginTime     uint64          // Login time is only available after login
}

// initialize
func NewClient(addr string, socket *websocket.Conn, firstTime uint64) (client *Client) {
	client = &Client{
		Addr:          addr,
		Socket:        socket,
		Send:          make(chan []byte, 100),
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
	}

	return
}

// read client data
func (c *Client) GetKey() (key string) {
	key = GetUserKey(c.AppId, c.UserId)

	return
}

// read client data
func (c *Client) read() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("write stop", string(debug.Stack()), r)
		}
	}()

	defer func() {
		fmt.Println("Read client data and close send", c)
		close(c.Send)
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			fmt.Println("Error reading client data", c.Addr, err)

			return
		}

		// handler
		fmt.Println("Read client data processing:", string(message))
		ProcessData(c, message)
	}
}

// write data to client
func (c *Client) write() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("write stop", string(debug.Stack()), r)

		}
	}()

	defer func() {
		clientManager.Unregister <- c
		c.Socket.Close()
		fmt.Println("Client sends data defer", c)
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// send data error close the connection
				fmt.Println("Client sends data to close the connection", c.Addr, "ok", ok)

				return
			}

			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// read client data
func (c *Client) SendMsg(msg []byte) {

	if c == nil {

		return
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("SendMsg stop:", r, string(debug.Stack()))
		}
	}()

	c.Send <- msg
}

// read client data
func (c *Client) close() {
	close(c.Send)
}

// User login
func (c *Client) Login(appId uint32, userId string, loginTime uint64) {
	c.AppId = appId
	c.UserId = userId
	c.LoginTime = loginTime
	// successful login = heartbeat once
	c.Heartbeat(loginTime)
}

// user heartbeat
func (c *Client) Heartbeat(currentTime uint64) {
	c.HeartbeatTime = currentTime

	return
}

// heartbeat timeout
func (c *Client) IsHeartbeatTimeout(currentTime uint64) (timeout bool) {
	if c.HeartbeatTime+heartbeatExpirationTime <= currentTime {
		timeout = true
	}

	return
}

// are you logged in?
func (c *Client) IsLogin() (isLogin bool) {

	// user is logged in
	if c.UserId != "" {
		isLogin = true

		return
	}

	return
}
